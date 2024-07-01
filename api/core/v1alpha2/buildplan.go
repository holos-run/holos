package v1alpha2

// FilePath represents a file path.
type FilePath string

// FileContent represents file contents.
type FileContent string

// FileContentMap represents a mapping of file paths to file contents.  Paths
// are relative to the `holos` output "deploy" directory, and may contain
// sub-directories.
type FileContentMap map[FilePath]FileContent

// BuildPlan represents a build plan for the holos cli to execute.  The purpose
// of a BuildPlan is to define one or more [HolosComponent] kinds.  For example a
// [HelmChart], [KustomizeBuild], or [KubernetesObjects].
//
// A BuildPlan usually has an additional empty [KubernetesObjects] for the
// purpose of using the [HolosComponent] DeployFiles field to deploy an ArgoCD
// or Flux gitops resource for the holos component.
type BuildPlan struct {
	Kind       string        `json:"kind" cue:"\"BuildPlan\""`
	APIVersion string        `json:"apiVersion" cue:"string | *\"v1alpha2\""`
	Spec       BuildPlanSpec `json:"spec"`
}

// BuildPlanSpec represents the specification of the build plan.
type BuildPlanSpec struct {
	// Disabled causes the holos cli to take no action over the [BuildPlan].
	Disabled bool `json:"disabled,omitempty"`
	// Components represents multiple [HolosComponent] kinds to manage.
	Components BuildPlanComponents `json:"components,omitempty"`
}

type BuildPlanComponents struct {
	Resources             map[Label]KubernetesObjects `json:"resources,omitempty"`
	KubernetesObjectsList []KubernetesObjects         `json:"kubernetesObjectsList,omitempty"`
	HelmChartList         []HelmChart                 `json:"helmChartList,omitempty"`
	KustomizeBuildList    []KustomizeBuild            `json:"kustomizeBuildList,omitempty"`
}

// HolosComponent defines the fields common to all holos component kinds.  Every
// holos component kind should embed HolosComponent.
type HolosComponent struct {
	// Kind is a string value representing the resource this object represents.
	Kind string `json:"kind"`
	// APIVersion represents the versioned schema of this representation of an object.
	APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha2\""`
	// Metadata represents data about the holos component such as the Name.
	Metadata Metadata `json:"metadata"`

	// APIObjectMap holds the marshalled representation of api objects.  Useful to
	// mix in resources to each HolosComponent type, for example adding an
	// ExternalSecret to a HelmChart HolosComponent.  Refer to [APIObjects].
	APIObjectMap APIObjectMap `json:"apiObjectMap,omitempty"`

	// DeployFiles represents file paths relative to the cluster deploy directory
	// with the value representing the file content.  Intended for defining the
	// ArgoCD Application resource or Flux Kustomization resource from within CUE,
	// but may be used to render any file related to the build plan from CUE.
	DeployFiles FileContentMap `json:"deployFiles,omitempty"`

	// Kustomize represents a kubectl kustomize build post-processing step.
	Kustomize `json:"kustomize,omitempty"`

	// Skip causes holos to take no action regarding this component.
	Skip bool `json:"skip" cue:"bool | *false"`
}

// Metadata represents data about the holos component such as the Name.
type Metadata struct {
	// Name represents the name of the holos component.
	Name string `json:"name"`
	// Namespace is the primary namespace of the holos component.  A holos
	// component may manage resources in multiple namespaces, in this case
	// consider setting the component namespace to default.
	//
	// This field is optional because not all resources require a namespace,
	// particularly CRD's and DeployFiles functionality.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// Kustomize represents resources necessary to execute a kustomize build.
// Intended for at least two use cases:
//
//  1. Process a [KustomizeBuild] [HolosComponent] which represents raw yaml
//     file resources in a holos component directory.
//  2. Post process a [HelmChart] [HolosComponent] to inject istio, patch jobs,
//     add custom labels, etc...
type Kustomize struct {
	// KustomizeFiles holds file contents for kustomize, e.g. patch files.
	KustomizeFiles FileContentMap `json:"kustomizeFiles,omitempty"`
	// ResourcesFile is the file name used for api objects in kustomization.yaml
	ResourcesFile string `json:"resourcesFile,omitempty"`
}
