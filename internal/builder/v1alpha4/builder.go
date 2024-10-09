package v1alpha4

import (
	"context"
	"fmt"
	"io"
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

		// Run the Generators
		// TODO: run these concurrently.
		for _, g := range a.Generators {
			switch g.Kind {
			case "Resources":
				if err := b.generateResources(g, am); err != nil {
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

		// Run the Transformers
	}

	return errors.NotImplemented()
}

func (b *BuildPlan) generateResources(g v1alpha4.Generator, am holos.ArtifactMap) error {
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

	data, err := yaml.Marshal(list)
	if err != nil {
		return errors.Format("could not generate %s for %s: %w", g.Output, b.BuildPlan.Metadata.Name, err)
	}

	if err := am.Set(holos.FilePath(g.Output), data); err != nil {
		return errors.Format("could not generate %s for %s: %w", g.Output, b.BuildPlan.Metadata.Name, err)
	}

	return nil
}
