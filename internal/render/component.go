// Deprecated: use cli instead
package render

import (
	"context"
	"fmt"

	"github.com/holos-run/holos"
	core "github.com/holos-run/holos/api/core/v1alpha3"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
	"github.com/holos-run/holos/internal/util"
)

const KubernetesObjectsKind = "KubernetesObjects"

// KubernetesObjects represents CUE output which directly provides Kubernetes api objects to holos.
type KubernetesObjects struct {
	Component core.KubernetesObjects `json:"component" yaml:"component"`
}

// Render produces kubernetes api objects from the APIObjectMap of the holos component.
func (o *KubernetesObjects) Render(ctx context.Context, path holos.InstancePath) (*Result, error) {
	result := NewResult(o.Component.Component)
	result.addObjectMap(ctx, o.Component.APIObjectMap)
	if err := result.kustomize(ctx); err != nil {
		return nil, errors.Wrap(fmt.Errorf("could not kustomize: %w", err))
	}
	return result, nil
}

// KustomizeBuild renders plain yaml files in the holos component directory
// using kubectl kustomize build.
type KustomizeBuild struct {
	Component core.KustomizeBuild `json:"component" yaml:"component"`
}

// Render produces a Result by executing kubectl kustomize on the holos
// component path. Useful for processing raw yaml files.
func (kb *KustomizeBuild) Render(ctx context.Context, path holos.InstancePath) (*Result, error) {
	if kb == nil {
		return nil, nil
	}
	log := logger.FromContext(ctx)
	result := NewResult(kb.Component.Component)
	// Run kustomize.
	kOut, err := util.RunCmd(ctx, "kubectl", "kustomize", string(path))
	if err != nil {
		log.ErrorContext(ctx, kOut.Stderr.String())
		return nil, errors.Wrap(err)
	}
	// Replace the accumulated output
	result.accumulatedOutput = kOut.Stdout.String()
	// Add CUE based api objects.
	result.addObjectMap(ctx, kb.Component.APIObjectMap)
	return result, nil
}
