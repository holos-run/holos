package component

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/component/v1alpha5"
	"github.com/holos-run/holos/internal/component/v1alpha6"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"gopkg.in/yaml.v3"
)

type BuildPlan struct {
	holos.BuildPlan
}

// New returns a new Component used to obtain a BuildPlan.
func New(root string, path string, cfg Config) *Component {
	return &Component{
		Config: cfg,
		Root:   root,
		Path:   path,
	}
}

// Component implements the holos render component command.
type Component struct {
	Config
	// Root represents the cue module root directory.
	Root string
	// Path represents the component path relative to Root.
	Path string
}

// TypeMeta returns the [holos.TypeMeta] of the resource the component produces.
// Useful to discriminate behavior.  If the type meta file does not exist
// TypeMeta returns a v1alpha5 APIVersion BuildPlan Kind.
func (c *Component) TypeMeta() (tm holos.TypeMeta, err error) {
	// if typemeta.yaml does not exist, assume v1alpha5 BuildPlan
	tmPath := filepath.Join(c.Root, c.Path, holos.TypeMetaFile)
	if _, err = os.Stat(tmPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return tm, errors.Wrap(err)
		}
		slog.Debug(fmt.Sprintf("could not load %s assuming v1alpha5", tmPath), "path", tmPath, "err", err)
		tm.APIVersion = "v1alpha5"
		tm.Kind = "BuildPlan"
		return tm, nil
	}
	data, err := os.ReadFile(tmPath)
	if err != nil {
		return tm, errors.Wrap(err)
	}
	if err = yaml.Unmarshal(data, &tm); err != nil {
		return tm, errors.Wrap(err)
	}
	return tm, nil
}

// Render renders the component BuildPlan.
func (c *Component) Render(ctx context.Context) error {
	tm, err := c.TypeMeta()
	if err != nil {
		return errors.Format("could not discriminate component type: %w", err)
	}

	switch tm.APIVersion {
	case "v1alpha6":
		if err := c.render(ctx, tm); err != nil {
			return errors.Format("could not render component: %w", err)
		}
	case "v1alpha5":
		if err := c.renderAlpha5(ctx); err != nil {
			return errors.Format("could not render v1alpha5 component: %w", err)
		}
	default:
		return errors.Format("unsupported version: %v", tm.APIVersion)
	}
	return nil
}

// BuildPlan returns the BuildPlan for the component.
func (c *Component) BuildPlan(tm holos.TypeMeta, opts holos.BuildOpts) (BuildPlan, error) {
	// Generic build plan wrapper for all api versions.
	var bp BuildPlan
	// All versions allow tags explicitly injected using the --inject flag.
	tags := c.TagMap.Tags()
	// discriminate the version.
	switch tm.APIVersion {
	case "v1alpha6":
		// Prepare runtime build context for injection as a cue tag.
		bc := v1alpha6.NewBuildContext(opts.BuildContext)
		buildContextTags, err := bc.Tags()
		if err != nil {
			return bp, errors.Format("could not get build context tag: %w", err)
		}
		// Append the standard tags for the component name, labels, annotations.
		tags = append(tags, opts.Tags...)
		// Append build context tags such as the holos managed temp directory.
		tags = append(tags, buildContextTags...)
		// the version specific build plan itself embedded into the wrapper.
		bp = BuildPlan{BuildPlan: &v1alpha6.BuildPlan{Opts: opts}}
	case "v1alpha5":
		// Append the standard tags for the component name, labels, annotations.
		tags = append(tags, opts.Tags...)
		bp = BuildPlan{BuildPlan: &v1alpha5.BuildPlan{Opts: opts}}
	default:
		return bp, errors.Format("unsupported version: %s", tm.APIVersion)
	}

	inst, err := BuildInstance(c.Root, c.Path, tags)
	if err != nil {
		return bp, errors.Format("could not load cue instance: %w", err)
	}

	// Get the holos field value from cue.
	v, err := inst.HolosValue()
	if err != nil {
		return bp, errors.Wrap(err)
	}

	// Load the BuildPlan from the cue value.
	if err := bp.Load(v); err != nil {
		return bp, errors.Wrap(err)
	}

	return bp, nil
}

// render implements the behavior of holos render component for v1alpha6 and
// later component versions.  The typemeta.yaml file located in the component
// directory must be present and is used to discriminate the apiVersion prior to
// building the CUE instance.  Useful to determine which build tags need to be
// injected depending on the apiVersion of the component.
func (c *Component) render(ctx context.Context, tm holos.TypeMeta) error {
	if tm.Kind != "BuildPlan" {
		return errors.Format("unsupported kind: %s, want BuildPlan", tm.Kind)
	}
	// temp directory is an important part of the build context.
	tempDir, err := os.MkdirTemp("", "holos.render")
	if err != nil {
		return errors.Format("could not make temp dir: %w", err)
	}
	defer util.Remove(ctx, tempDir)

	// Runtime configuration of the build.
	opts := holos.NewBuildOpts(c.Path)
	opts.Stderr = c.Stderr
	opts.Concurrency = c.Concurrency
	opts.WriteTo = filepath.Join(c.Root, c.WriteTo)
	opts.BuildContext.TempDir = tempDir

	log := logger.FromContext(ctx)
	log.DebugContext(ctx, fmt.Sprintf("rendering %s kind %s version %s", c.Path, tm.Kind, tm.APIVersion), "kind", tm.Kind, "apiVersion", tm.APIVersion, "path", c.Path)

	// Get the BuildPlan from cue.
	bp, err := c.BuildPlan(tm, opts)
	if err != nil {
		return errors.Wrap(err)
	}
	// Execute the build.
	if err := bp.Build(ctx); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

// renderComponentAlpha5 implements the behavior of holos render component for
// v1alpha5 and earlier.  This method loads the CUE Instance to discriminate the
// apiVersion, which is too late to pass tags properly.
//
// Deprecated: use render() instead
func (c *Component) renderAlpha5(ctx context.Context) error {
	// Manage a temp directory for the build artifacts.  The concrete value is
	// needed prior to exporting the BuildPlan from the CUE instance.
	tempDir, err := os.MkdirTemp("", "holos.render")
	if err != nil {
		return errors.Format("could not make temp dir: %w", err)
	}
	defer util.Remove(ctx, tempDir)

	// Runtime configuration of the build.
	opts := holos.NewBuildOpts(c.Path)
	opts.Stderr = c.Stderr
	opts.Concurrency = c.Concurrency
	opts.WriteTo = filepath.Join(c.Root, c.WriteTo)

	tm := holos.TypeMeta{
		Kind:       "BuildPlan",
		APIVersion: "v1alpha5",
	}

	// Get the BuildPlan from cue.
	bp, err := c.BuildPlan(tm, opts)
	if err != nil {
		return errors.Wrap(err)
	}

	// Execute the build.
	if err := bp.Build(ctx); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
