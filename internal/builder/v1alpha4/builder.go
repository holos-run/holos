package v1alpha4

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/api/core/v1alpha4"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

// Platform represents a platform builder.
type Platform struct {
	Platform    v1alpha4.Platform
	Concurrency int
	Stderr      io.Writer
}

// Build builds a Platform by concurrently building a BuildPlan for each
// platform component.  No artifact files are written directly, only indirectly
// by rendering each component.
func (p *Platform) Build(ctx context.Context, _ holos.ArtifactMap) error {
	parentStart := time.Now()
	components := p.Platform.Spec.Components
	total := len(components)
	g, ctx := errgroup.WithContext(ctx)
	// Limit the number of concurrent goroutines due to CUE memory usage concerns
	// while rendering components.  One more for the producer.
	g.SetLimit(p.Concurrency + 1)
	// Spawn a producer because g.Go() blocks when the group limit is reached.
	g.Go(func() error {
		for idx := range components {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Capture idx to avoid issues with closure.  Fixed in Go 1.22.
				idx := idx
				buildContext := &components[idx]
				// Worker go routine.  Blocks if limit has been reached.
				g.Go(func() error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						start := time.Now()
						log := logger.FromContext(ctx).With(
							"path", buildContext.Path,
							"cluster", buildContext.Cluster,
							"environment", buildContext.Environment,
							"num", idx+1,
							"total", total,
						)
						log.DebugContext(ctx, "render component")

						tags := make([]string, 0, 2+len(buildContext.Tags))
						tags = append(tags, "component="+buildContext.Path)
						tags = append(tags, "environment="+buildContext.Environment)
						tags = append(tags, buildContext.Tags...)

						// Execute a sub-process to limit CUE memory usage.
						args := []string{
							"render",
							"component",
							"--cluster-name", buildContext.Cluster,
							"--tags", strings.Join(tags, ","),
							buildContext.Path,
						}
						result, err := util.RunCmd(ctx, "holos", args...)
						if err != nil {
							_, _ = io.Copy(p.Stderr, result.Stderr)
							return errors.Wrap(fmt.Errorf("could not render component: %w", err))
						}

						duration := time.Since(start)
						msg := fmt.Sprintf(
							"rendered %s for cluster %s in %s",
							filepath.Base(buildContext.Path),
							buildContext.Cluster,
							duration,
						)
						log.InfoContext(ctx, msg, "duration", duration)
						return nil
					}
				})
			}
		}
		return nil
	})

	// Wait for completion and return the first error (if any)
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Since(parentStart)
	msg := fmt.Sprintf("rendered platform in %s", duration)
	logger.FromContext(ctx).InfoContext(ctx, msg, "duration", duration, "version", p.Platform.APIVersion)
	return nil
}

// BuildPlan represents a component builder.
type BuildPlan struct {
	BuildPlan   v1alpha4.BuildPlan
	Concurrency int
	Stderr      io.Writer
}

// Build builds a BuildPlan into Artifact files.
func (b *BuildPlan) Build(ctx context.Context, am holos.ArtifactMap) error {
	name := b.BuildPlan.Metadata.Name
	component := b.BuildPlan.Spec.Component
	log := logger.FromContext(ctx).With("name", name, "component", component)
	log.DebugContext(ctx, "building "+name)

	if b.BuildPlan.Spec.Disabled {
		log.WarnContext(ctx, fmt.Sprintf("could not build %s: spec.disabled field is true", name))
		return nil
	}

	for _, a := range b.BuildPlan.Spec.Artifacts {
		if a.Skip {
			log.WarnContext(ctx, fmt.Sprintf("skipped artifact %s: skip field is true", a.Artifact))
			continue
		}

		// TODO(jeff): concurrent generators.
		for _, g := range a.Generators {
			switch g.Kind {
			case "Resources":
				if err := b.generateResources(log, g, am); err != nil {
					return errors.Wrap(err)
				}
			case "Helm":
				return errors.NotImplemented()
			case "File":
				return errors.NotImplemented()
			default:
				return errors.Format("could not build %s: unsupported kind %s", name, g.Kind)
			}
		}

		for _, t := range a.Transformers {
			switch t.Kind {
			case "Kustomize":
				if err := b.kustomize(ctx, log, t, am); err != nil {
					return errors.Wrap(err)
				}
			case "Join":
				return errors.NotImplemented()
			default:
				return errors.Format("could not build %s: unsupported kind %s", name, t.Kind)
			}
		}
	}

	return errors.NotImplemented()
}

func (b *BuildPlan) generateResources(
	log *slog.Logger,
	g v1alpha4.Generator,
	am holos.ArtifactMap,
) error {
	var size int
	for _, m := range g.Resources {
		size += len(m)
	}
	list := make([]v1alpha4.Resource, 0, size)

	for _, m := range g.Resources {
		for _, r := range m {
			list = append(list, r)
		}
	}

	msg := fmt.Sprintf("could not generate %s for %s", g.Output, b.BuildPlan.Metadata.Name)

	buf, err := marshal(list)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	if err := am.Set(holos.FilePath(g.Output), buf.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}

	log.Debug("set artifact: " + string(g.Output))
	return nil
}

func (b *BuildPlan) kustomize(
	ctx context.Context,
	log *slog.Logger,
	t v1alpha4.Transformer,
	am holos.ArtifactMap,
) error {
	tempDir, err := os.MkdirTemp("", "holos.kustomize")
	if err != nil {
		return errors.Wrap(err)
	}
	defer util.Remove(ctx, tempDir)
	msg := fmt.Sprintf("could not transform %s for %s", t.Output, b.BuildPlan.Metadata.Name)

	// Write the kustomization
	data, err := yaml.Marshal(t.Kustomize.Kustomization)
	if err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	path := filepath.Join(tempDir, "kustomization.yaml")
	if err := os.WriteFile(path, data, 0666); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	log.DebugContext(ctx, "wrote: "+path)

	// Write the inputs
	for _, input := range t.Inputs {
		data, ok := am.Get(holos.FilePath(input))
		if !ok {
			return errors.Format("%s: could not get artifact %s", msg, input)
		}
		path := filepath.Join(tempDir, string(input))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		if err := os.WriteFile(path, data, 0666); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		log.DebugContext(ctx, "wrote: "+path)
	}

	// Execute kustomize
	result, err := util.RunCmd(ctx, "kubectl", "kustomize", tempDir)
	if err != nil {
		log.ErrorContext(ctx, result.Stderr.String())
		return errors.Format("%s: could not run kustomize: %w", msg, err)
	}

	// Store the artifact
	if err := am.Set(holos.FilePath(t.Output), result.Stdout.Bytes()); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	log.Debug("set artifact: " + string(t.Output))

	return nil
}

func marshal(list []v1alpha4.Resource) (buf bytes.Buffer, err error) {
	encoder := yaml.NewEncoder(&buf)
	defer encoder.Close()
	for _, item := range list {
		if err = encoder.Encode(item); err != nil {
			err = errors.Wrap(err)
			return
		}
	}
	return
}
