package v1alpha1

import (
	"context"
	"github.com/holos-run/holos"
)

const KubernetesObjectsKind = "KubernetesObjects"

// KubernetesObjects represents CUE output which directly provides Kubernetes api objects to holos.
type KubernetesObjects struct {
	HolosComponent `json:",inline" yaml:",inline"`
}

// Render produces kubernetes api objects from the APIObjectMap
func (o *KubernetesObjects) Render(ctx context.Context, path holos.PathComponent) (*Result, error) {
	result := Result{
		TypeMeta:      o.TypeMeta,
		Metadata:      o.Metadata,
		Kustomization: o.Kustomization,
	}
	result.addObjectMap(ctx, o.APIObjectMap)
	return &result, nil
}
