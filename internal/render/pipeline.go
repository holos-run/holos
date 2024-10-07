package render

import (
	"context"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/artifact"
	"github.com/holos-run/holos/internal/errors"
)

// Platform renders a platform, writing fully rendered manifests to files.
func Platform(ctx context.Context, b holos.Builder) error {
	// Artifacts are currently written by each `holos render component`
	// subprocess, not the parent `holos render platform` process.
	if err := b.Build(ctx, artifact.New()); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// Component renders a component writing fully rendered manifests to files.
func Component(ctx context.Context, b holos.Builder, a holos.Artifact) error {
	if err := b.Build(ctx, a); err != nil {
		return errors.Wrap(err)
	}
	if err := a.Save(ctx); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
