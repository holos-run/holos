package v1alpha1

import (
	"fmt"
	"strings"
)

// BuildPlan is the primary interface between CUE and the Holos cli.
type BuildPlan struct {
	TypeMeta `json:",inline" yaml:",inline"`
	// Metadata represents the holos component name
	Metadata ObjectMeta     `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec     BuildPlanSpec  `json:"spec,omitempty" yaml:"spec,omitempty"`
	Platform map[string]any `json:"platform,omitempty" yaml:"platform,omitempty"`
}

type BuildPlanSpec struct {
	Disabled   bool                `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Components BuildPlanComponents `json:"components,omitempty" yaml:"components,omitempty"`
}

type BuildPlanComponents struct {
	HelmChartList         []HelmChart                  `json:"helmChartList,omitempty" yaml:"helmChartList,omitempty"`
	KubernetesObjectsList []KubernetesObjects          `json:"kubernetesObjectsList,omitempty" yaml:"kubernetesObjectsList,omitempty"`
	KustomizeBuildList    []KustomizeBuild             `json:"kustomizeBuildList,omitempty" yaml:"kustomizeBuildList,omitempty"`
	Resources             map[string]KubernetesObjects `json:"resources,omitempty" yaml:"resources,omitempty"`
}

func (bp *BuildPlan) Validate() error {
	errs := make([]string, 0, 2)
	if bp.Kind != BuildPlanKind {
		errs = append(errs, fmt.Sprintf("kind invalid: want: %s have: %s", BuildPlanKind, bp.Kind))
	}
	if bp.APIVersion != APIVersion {
		errs = append(errs, fmt.Sprintf("apiVersion invalid: want: %s have: %s", APIVersion, bp.APIVersion))
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid BuildPlan: " + strings.Join(errs, ", "))
	}
	return nil
}
