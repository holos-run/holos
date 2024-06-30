package v1alpha2

import (
	"fmt"
	"strings"
)

// Label is an arbitrary unique identifier.  Defined as a type for clarity and type checking.
type Label string

// Kind is a kubernetes api object kind. Defined as a type for clarity and type checking.
type Kind string

// APIObjectMap is the shape of marshalled api objects returned from cue to the
// holos cli. A map is used to improve the clarity of error messages from cue
// relative to a list.
type APIObjectMap map[Kind]map[Label]string

// FileContentMap represents a mapping of file names to file content.
type FileContentMap map[string]string

// BuildPlan represents a build plan for the holos cli to execute.  A build plan
// is a set of zero or more holos components.
type BuildPlan struct {
	// Kind is a string value representing the resource this object represents.
	Kind string `json:"kind" yaml:"kind" cue:"\"BuildPlan\""`
	// APIVersion represents the versioned schema of this representation of an object.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1alpha2\""`
	// Spec represents the specification.
	Spec BuildPlanSpec `json:"spec" yaml:"spec"`
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

func (bp *BuildPlan) ResultCapacity() (count int) {
	if bp == nil {
		return 0
	}
	count = len(bp.Spec.Components.HelmChartList)
	// +len(bp.Spec.Components.KubernetesObjectsList)
	// +len(bp.Spec.Components.KustomizeBuildList)
	// +len(bp.Spec.Components.Resources)
	return count
}

type BuildPlanSpec struct {
	Disabled   bool                `json:"disabled,omitempty" yaml:"disabled,omitempty"`
	Components BuildPlanComponents `json:"components,omitempty" yaml:"components,omitempty"`
}

type BuildPlanComponents struct {
	Resources             map[string]KubernetesObjects `json:"resources,omitempty" yaml:"resources,omitempty"`
	KubernetesObjectsList []KubernetesObjects          `json:"kubernetesObjectsList,omitempty" yaml:"kubernetesObjectsList,omitempty"`
	HelmChartList         []HelmChart                  `json:"helmChartList,omitempty" yaml:"helmChartList,omitempty"`
	KustomizeBuildList    []KustomizeBuild             `json:"kustomizeBuildList,omitempty" yaml:"kustomizeBuildList,omitempty"`
}

// HolosComponent defines the fields common to all holos component kinds.  Every
// holos component kind should embed HolosComponent.
type HolosComponent struct {
	// Kind is a string value representing the resource this object represents.
	Kind string `json:"kind" yaml:"kind"`
	// APIVersion represents the versioned schema of this representation of an object.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1alpha2\""`
	// Metadata represents data about the holos component such as the Name.
	Metadata Metadata `json:"metadata" yaml:"metadata"`

	// APIObjectMap holds the marshalled representation of api objects. Think of
	// these objects as being mixed into the upstream resources, for example
	// adding an ExternalSecret to a rendered Helm chart.
	APIObjectMap APIObjectMap `json:"apiObjectMap,omitempty" yaml:"apiObjectMap,omitempty"`

	// DeployFiles represents file paths relative to the cluster deploy directory
	// with the value representing the file content.  Intended for defining the
	// ArgoCD Application resource or Flux Kustomization resource from within CUE,
	// but may be used to render any file related to the build plan from CUE.
	DeployFiles FileContentMap `json:"deployFiles,omitempty" yaml:"deployFiles,omitempty"`

	// Kustomize represents a kubectl kustomize build post-processing step.
	Kustomize `json:"kustomize,omitempty" yaml:"kustomize,omitempty"`

	// Skip causes holos to take no action regarding this component.
	Skip bool `json:"skip" yaml:"skip" cue:"bool | *false"`
}

// Metadata represents data about the holos component such as the Name.
type Metadata struct {
	// Name represents the name of the holos component.
	Name string `json:"name" yaml:"name"`
	// Namespace is the primary namespace of the holos component.  A holos
	// component may manage resources in multiple namespaces, in this case
	// consider setting the component namespace to default.
	//
	// This field is optional because not all resources require a namespace,
	// particularly CRD's and DeployFiles functionality.
	// +optional
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

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
