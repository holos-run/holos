package component

import (
	"github.com/holos-run/holos/internal/component/v1alpha5"
	"github.com/holos-run/holos/internal/component/v1alpha6"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
)

type BuildPlan struct {
	holos.BuildPlan
}

// LoadBuildPlan loads a BuildPlan from a CUE Instance.
func LoadBuildPlan(i *Instance, opts holos.BuildOpts) (bp BuildPlan, err error) {
	err = i.Discriminate(func(tm holos.TypeMeta) error {
		if tm.Kind != "BuildPlan" {
			return errors.Format("unsupported kind: %s, want BuildPlan", tm.Kind)
		}

		switch version := tm.APIVersion; version {
		case "v1alpha5":
			bp = BuildPlan{&v1alpha5.BuildPlan{Opts: opts}}
		case "v1alpha6":
			bp = BuildPlan{&v1alpha6.BuildPlan{Opts: opts}}
		default:
			return errors.Format("unsupported version: %s", version)
		}

		return nil
	})
	if err != nil {
		return bp, errors.Wrap(err)
	}

	// Get the holos: field value from cue.
	v, err := i.HolosValue()
	if err != nil {
		return bp, errors.Wrap(err)
	}

	// Load the platform from the cue value.
	if err := bp.Load(v); err != nil {
		return bp, errors.Wrap(err)
	}

	return bp, err
}
