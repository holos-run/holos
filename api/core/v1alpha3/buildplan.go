package v1alpha3

// FilePath represents a file path.
type FilePath string

// FileContent represents file contents.
type FileContent string

// FileContentMap represents a mapping of file paths to file contents.
type FileContentMap map[FilePath]FileContent

// BuildPlan represents a build plan for the holos cli to execute.  The purpose
// of a BuildPlan is to define one or more [Component] kinds.  For example a
// [HelmChart], [KustomizeBuild], or [KubernetesObjects].
//
// A BuildPlan usually has an additional empty [KubernetesObjects] for the
// purpose of using the [Component] DeployFiles field to deploy an ArgoCD
// or Flux gitops resource for the holos component.
type BuildPlan struct {
	Kind       string        `json:"kind" cue:"\"BuildPlan\""`
	APIVersion string        `json:"apiVersion" cue:"string | *\"v1alpha3\""`
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
	Resources             map[InternalLabel]KubernetesObjects `json:"resources,omitempty"`
	KubernetesObjectsList []KubernetesObjects                 `json:"kubernetesObjectsList,omitempty"`
	HelmChartList         []HelmChart                         `json:"helmChartList,omitempty"`
	KustomizeBuildList    []KustomizeBuild                    `json:"kustomizeBuildList,omitempty"`
}

// Kustomize represents resources necessary to execute a kustomize build.
// Intended for at least two use cases:
//
//  1. Process a [KustomizeBuild] [Component] which represents raw yaml
//     file resources in a holos component directory.
//  2. Post process a [HelmChart] [Component] to inject istio, patch jobs,
//     add custom labels, etc...
type Kustomize struct {
	// KustomizeFiles holds file contents for kustomize, e.g. patch files.
	KustomizeFiles FileContentMap `json:"kustomizeFiles,omitempty"`
	// ResourcesFile is the file name used for api objects in kustomization.yaml
	ResourcesFile string `json:"resourcesFile,omitempty"`
}
