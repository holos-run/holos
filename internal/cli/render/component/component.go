package component

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/builder"
	"github.com/holos-run/holos/internal/builder/v1alpha6"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/util"
	"gopkg.in/yaml.v2"
)

// New returns a new Component renderer.
func New(path string, cfg Config) *Component {
	return &Component{
		Config: cfg,
		Path:   path,
	}
}

// Component implements the holos render component command.
type Component struct {
	Config
	// Path represents the component path relative to the cue module root.
	Path string
}

func (c *Component) Render(ctx context.Context) error {
	log := logger.FromContext(ctx)

	// if typemeta.yaml does not exist, render using <= v1alpha5 behavior.
	typeMetaPath := filepath.Join(c.Path, holos.TypeMetaFile)
	if _, err := os.Stat(typeMetaPath); err != nil {
		log.DebugContext(ctx, fmt.Sprintf("could not load %s falling back to deprecated builder", typeMetaPath), "path", typeMetaPath, "err", err)
		return c.renderAlpha5(ctx)
	}

	// Render using the new typemeta.yaml discriminator in v1alpha6 and later.
	tm, err := loadTypeMeta(typeMetaPath)
	if err != nil {
		return errors.Format("could not discriminate component type: %w", err)
	}
	if err := c.render(ctx, tm); err != nil {
		return errors.Format("could not render component: %w", err)
	}
	return nil
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
	// so we can append version specific tags.
	tags := c.TagMap.Tags()
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
	opts.WriteTo = c.WriteTo
	opts.BuildContext.TempDir = tempDir

	log := logger.FromContext(ctx)
	log.DebugContext(ctx, fmt.Sprintf("rendering %s kind %s version %s", c.Path, tm.Kind, tm.APIVersion), "kind", tm.Kind, "apiVersion", tm.APIVersion, "path", c.Path)

	// generic build plan wrapper for all versions.
	var bp builder.BuildPlan
	switch version := tm.APIVersion; version {
	case "v1alpha6":
		// Prepare runtime build context for injection as a cue tag.
		bc := v1alpha6.NewBuildContext(opts.BuildContext)
		buildContextTags, err := bc.Tags()
		if err != nil {
			return errors.Format("could not get build context tag: %w", err)
		}
		tags = append(tags, buildContextTags...)
		// the version specific build plan itself embedded into the wrapper.
		bp = builder.BuildPlan{BuildPlan: &v1alpha6.BuildPlan{Opts: opts}}
	default:
		return errors.Format("unsupported version: %s, must be at least v1alpha6 when typemeta.yaml is present", version)
	}

	// Load the CUE instance to export the BuildPlan with version-specific tags.
	inst, err := builder.LoadInstance(c.Path, nil, tags)
	if err != nil {
		return errors.Format("could not load cue instance: %w", err)
	}

	// Get the holos: field value from cue.
	v, err := inst.HolosValue()
	if err != nil {
		return errors.Wrap(err)
	}

	// Load the BuildPlan from the cue value.
	if err := bp.Load(v); err != nil {
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

	// Load the CUE instance to export the BuildPlan.
	inst, err := builder.LoadInstance(c.Path, nil, c.TagMap.Tags())
	if err != nil {
		return errors.Format("could not load cue instance: %w", err)
	}

	// Runtime configuration of the build.
	opts := holos.NewBuildOpts(c.Path)
	opts.Stderr = c.Stderr
	opts.Concurrency = c.Concurrency
	opts.WriteTo = c.WriteTo

	// Export the BuildPlan from the CUE instance.
	bp, err := builder.LoadBuildPlan(inst, opts)
	if err != nil {
		return errors.Wrap(err)
	}

	// Execute the build.
	if err := bp.Build(ctx); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func loadTypeMeta(path string) (tm holos.TypeMeta, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return tm, errors.Wrap(err)
	}
	if err = yaml.Unmarshal(data, &tm); err != nil {
		return tm, errors.Wrap(err)
	}
	return tm, nil
}
