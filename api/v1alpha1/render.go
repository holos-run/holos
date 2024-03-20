package v1alpha1

import (
	"context"
	"github.com/holos-run/holos"
)

type Renderer interface {
	GetKind() string
	Render(ctx context.Context, path holos.PathComponent) (*Result, error)
}

// Render produces a Result representing the kubernetes api objects to
// configure. Each of the various holos component types, e.g. Helm, Kustomize,
// et al, should implement the Renderer interface. This process is best
// conceptualized as a data pipeline, for example a component may render a
// result by first calling helm template, then passing the result through
// kustomize, then mixing in overlay api objects.
func Render(ctx context.Context, r Renderer, path holos.PathComponent) (*Result, error) {
	return r.Render(ctx, path)
}
