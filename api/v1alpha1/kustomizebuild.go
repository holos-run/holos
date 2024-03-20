package v1alpha1

import (
	"context"
	"github.com/holos-run/holos"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/util"
	"github.com/holos-run/holos/pkg/wrapper"
)

const KustomizeBuildKind = "KustomizeBuild"

// KustomizeBuild
type KustomizeBuild struct {
	HolosComponent `json:",inline" yaml:",inline"`
}

// Render produces kubernetes api objects from the APIObjectMap
func (kb *KustomizeBuild) Render(ctx context.Context, path holos.PathComponent) (*Result, error) {
	log := logger.FromContext(ctx)
	result := Result{
		TypeMeta:      kb.TypeMeta,
		Metadata:      kb.Metadata,
		Kustomization: kb.Kustomization,
	}
	// Run kustomize.
	kOut, err := util.RunCmd(ctx, "kubectl", "kustomize", string(path))
	if err != nil {
		log.ErrorContext(ctx, kOut.Stderr.String())
		return nil, wrapper.Wrap(err)
	}
	// Replace the accumulated output
	result.accumulatedOutput = kOut.Stdout.String()
	// Add CUE based api objects.
	result.addObjectMap(ctx, kb.APIObjectMap)
	return &result, nil
}
