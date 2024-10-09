// Package v1alpha4 contains the core API contract between the holos cli and CUE
// configuration code.  Platform designers, operators, and software developers
// use this API to write configuration in CUE which `holos` loads.  The overall
// shape of the API defines imperative actions `holos` should carry out to
// render the complete yaml that represents a Platform.
//
// [Platform] defines the complete configuration of a platform.  With the holos
// reference platform this takes the shape of one management cluster and at
// least two workload clusters.  Each cluster has multiple [Component] resources
// applied to it.
//
// Each holos component path, e.g. `components/namespaces` produces exactly one
// [BuildPlan] which produces an [Artifact].  An [Artifact] is a collection of
// fully rendered manifest files written to the filesystem.
package v1alpha4

//go:generate ../../../hack/gendoc

// BuildPlan represents a build plan for holos to execute.  Each [Platform]
// component produces exactly one BuildPlan.
type BuildPlan struct {
	// Kind represents the type of the resource.
	Kind string `json:"kind" cue:"\"BuildPlan\""`
	// APIVersion represents the versioned schema of the resource.
	APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha4\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata"`
	// Spec specifies the desired state of the resource.
	Spec BuildPlanSpec `json:"spec"`
}

// BuildPlanSpec represents the specification of the [BuildPlan].
type BuildPlanSpec struct {
	// Component represents the component that produced the build plan.
	// Represented as a path relative to the platform root.
	Component string `json:"component"`
	// Disabled causes the holos cli to disregard the build plan.
	Disabled bool `json:"disabled,omitempty"`
	// Steps represent build steps for holos to execute
	Steps []BuildStep `json:"steps"`
}

// BuildStep represents the holos rendering pipeline for a [BuildPlan].
//
// Each [Generator] may be executed concurrently with other generators in the
// same collection. Each [Transformer] is executed sequentially, the first after
// all generators have completed.
//
// Each BuildStep produces one manifest file artifact.  [Generator] manifests are
// implicitly joined into one artifact file if there is no [Transformer] that
// would otherwise combine them.
type BuildStep struct {
	Artifact     FilePath      `json:"artifact,omitempty"`
	Generators   []Generator   `json:"generators,omitempty"`
	Transformers []Transformer `json:"transformers,omitempty"`
	Skip         bool          `json:"skip,omitempty"`
}

// Generator generates an intermediate manifest for a [BuildStep].
//
// Each Generator in a [BuildStep] must have a distinct manifest value for a
// [Transformer] to reference.
type Generator struct {
	// Kind represents the kind of generator.  Must be Resources, Helm, or File.
	Kind string `json:"kind" cue:"\"Resources\" | \"Helm\" | \"File\""`
	// Manifest represents the output file for subsequent transformers.
	Manifest string `json:"manifest"`
	// Resources generator. Ignored unless kind is Resources.
	Resources Resources `json:"resources,omitempty"`
	// Helm generator. Ignored unless kind is Helm.
	Helm Helm `json:"helm,omitempty"`
	// File generator. Ignored unless kind is File.
	File File `json:"file,omitempty"`
}

// Resource represents one kubernetes api object.
type Resource map[string]any

// Resources represents a kubernetes resources [Generator] from CUE.
type Resources map[Kind]map[InternalLabel]Resource

// File represents a simple single file copy [Generator].  Useful with a
// [Kustomize] [Transformer] to process plain manifest files stored in the
// component directory.  Multiple File generators may be used to transform
// multiple resources.
type File struct {
	// Source represents a file to read relative to the component path, the
	// [BuildPlanSpec] Component field.
	Source FilePath `json:"source"`
}

// Helm represents a [Chart] manifest [Generator].
type Helm struct {
	// Chart represents a helm chart to manage.
	Chart Chart `json:"chart"`
	// Values represents values for holos to marshal into values.yaml when
	// rendering the chart.
	Values Values `json:"values"`
	// EnableHooks enables helm hooks when executing the `helm template` command.
	EnableHooks bool `json:"enableHooks,omitempty"`
}

// Values represents [Helm] Chart values generated from CUE.
type Values map[string]any

// Chart represents a [Helm] Chart.
type Chart struct {
	// Name represents the chart name.
	Name string `json:"name"`
	// Version represents the chart version.
	Version string `json:"version"`
	// Release represents the chart release when executing helm template.
	Release string `json:"release"`
	// Repository represents the repository to fetch the chart from.
	Repository Repository `json:"repository,omitempty"`
}

// Repository represents a [Helm] [Chart] repository.
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Transformer transforms [Generator] manifests within a [BuildStep].
type Transformer struct {
	// Kind represents the kind of transformer.  Must be Kustomize.
	Kind string `json:"kind" cue:"\"Kustomize\""`
	// Manifest represents the output file for subsequent transformers.
	Manifest string `json:"manifest,omitempty"`
	// Kustomize transformer. Ignored unless kind is Kustomize.
	Kustomize Kustomize `json:"kustomize,omitempty"`
}

// Kustomize represents a kustomization [Transformer].
type Kustomize struct {
	// Kustomization represents the decoded kustomization.yaml file
	Kustomization Kustomization `json:"kustomization"`
	// Files holds file contents for kustomize, e.g. patch files.
	Files FileContentMap `json:"files,omitempty"`
}

// Kustomization represents a kustomization.yaml file for use with the
// [Kustomize] [Transformer].  Untyped to avoid tightly coupling holos to
// kubectl versions which was a problem for the Flux maintainers.  Type checking
// is expected to happen in CUE against the kubectl version the user prefers.
type Kustomization map[string]any

// FileContent represents file contents.
type FileContent string

// FileContentMap represents a mapping of file paths to file contents.
type FileContentMap map[FilePath]FileContent

// FilePath represents a file path.
type FilePath string

// InternalLabel is an arbitrary unique identifier internal to holos itself.
// The holos cli is expected to never write a InternalLabel value to rendered
// output files, therefore use a InternalLabel when the identifier must be
// unique and internal.  Defined as a type for clarity and type checking.
type InternalLabel string

// Kind is a discriminator. Defined as a type for clarity and type checking.
type Kind string

// NameLabel is a unique identifier useful to convert a CUE struct to a list
// when the values have a Name field with a default value.  NameLabel indicates
// the common use case of converting a struct to a list where the Name field of
// the value aligns with the outer struct field name.
//
// For example:
//
//	Outer: [NAME=_]: Name: NAME
type NameLabel string

// Platform represents a platform to manage.  A Platform resource informs holos
// which components to build.  The platform resource also acts as a container
// for the platform model form values provided by the PlatformService.  The
// primary use case is to collect the cluster names, cluster types, platform
// model, and holos components to build into one resource.
type Platform struct {
	// Kind is a string value representing the resource.
	Kind string `json:"kind" cue:"\"Platform\""`
	// APIVersion represents the versioned schema of this resource.
	APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha4\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata"`

	// Spec represents the specification.
	Spec PlatformSpec `json:"spec"`
}

// Metadata represents data about the resource such as the Name.
type Metadata struct {
	// Name represents the resource name.
	Name string `json:"name"`
}

// PlatformSpec represents the specification of a Platform.  Think of a platform
// specification as a list of platform components to apply to a list of
// kubernetes clusters combined with the user-specified Platform Model.
type PlatformSpec struct {
	// Components represents a list of holos components to manage.
	Components []BuildContext `json:"components"`
}

// BuildContext represents the context necessary to render a component into a
// BuildPlan.  Useful to capture parameters passed down from a Platform spec for
// the purpose of idempotent rebuilds.
type BuildContext struct {
	// Path is the path of the component relative to the platform root.
	Path string `json:"path"`
	// Cluster is the cluster name to provide when rendering the component.
	Cluster string `json:"cluster"`
	// Environment for example, dev, test, stage, prod
	Environment string `json:"environment,omitempty"`
	// Model represents the platform model holos gets from from the
	// PlatformService.GetPlatform rpc method and provides to CUE using a tag.
	Model map[string]any `json:"model"`
	// Tags represents cue tags to provide when rendering the component.
	Tags []string `json:"tags,omitempty"`
}
