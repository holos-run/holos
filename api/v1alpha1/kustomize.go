package v1alpha1

import (
	"context"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/pkg/errors"
	"github.com/holos-run/holos/pkg/logger"
	"github.com/holos-run/holos/pkg/util"
)

const KustomizeBuildKind = "KustomizeBuild"

// Kustomize represents resources necessary to execute a kustomize build.
// Intended for at least two use cases:
//
//  1. Process raw yaml file resources in a holos component directory.
//  2. Post process a HelmChart to inject istio, add custom labels, etc...
type Kustomize struct {
	// KustomizeFiles holds file contents for kustomize, e.g. patch files.
	KustomizeFiles FileContentMap `json:"kustomizeFiles,omitempty" yaml:"kustomizeFiles,omitempty"`
	// ResourcesFile is the file name used for api objects in kustomization.yaml
	ResourcesFile string `json:"resourcesFile,omitempty" yaml:"resourcesFile,omitempty"`
}

// KustomizeBuild renders plain yaml files in the holos component directory using kubectl kustomize build.
type KustomizeBuild struct {
	HolosComponent `json:",inline" yaml:",inline"`
}

// Render produces a Result by executing kubectl kustomize on the holos
// component path. Useful for processing raw yaml files.
func (kb *KustomizeBuild) Render(ctx context.Context, path holos.InstancePath) (*Result, error) {
	log := logger.FromContext(ctx)
	result := Result{HolosComponent: kb.HolosComponent}
	// Run kustomize.
	kOut, err := util.RunCmd(ctx, "kubectl", "kustomize", string(path))
	if err != nil {
		log.ErrorContext(ctx, kOut.Stderr.String())
		return nil, errors.Wrap(err)
	}
	// Replace the accumulated output
	result.accumulatedOutput = kOut.Stdout.String()
	// Add CUE based api objects.
	result.addObjectMap(ctx, kb.APIObjectMap)
	return &result, nil
}
