// Package core contains schemas for a [Platform] and [BuildPlan].  Holos takes
// a [Platform] as input, then iterates over each [Component] to produce a
// [BuildPlan].  Holos processes the [BuildPlan] to produce fully rendered
// manifests, each an [Artifact].
package core

//go:generate ../../../hack/gendoc

// BuildPlan represents an implementation of the [rendered manifest pattern].
// Holos processes a BuildPlan to produce one or more [Artifact] output files.
// BuildPlan artifact files usually contain Kubernetes manifests, but they may
// have any content.
//
// A BuildPlan usually produces two artifacts.  One artifact contains a manifest
// of resources.  A second artifact contains a GitOps resource to manage the
// first, usually an ArgoCD Application resource.
//
// Holos uses CUE to construct a BuildPlan.  A future enhancement will support
// user defined executables providing a BuildPlan to Holos in the style of an
// [external credential provider].
//
// [rendered manifest pattern]: https://akuity.io/blog/the-rendered-manifests-pattern
// [external credential provider]: https://github.com/kubernetes/enhancements/blob/313ad8b59c80819659e1fbf0f165230f633f2b22/keps/sig-auth/541-external-credential-providers/README.md
type BuildPlan struct {
	// Kind represents the type of the resource.
	Kind string `json:"kind" cue:"\"BuildPlan\""`
	// APIVersion represents the versioned schema of the resource.
	APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha5\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata"`
	// Spec specifies the desired state of the resource.
	Spec BuildPlanSpec `json:"spec"`
}

// BuildPlanSpec represents the specification of the [BuildPlan].
type BuildPlanSpec struct {
	// Artifacts represents the artifacts for holos to build.
	Artifacts []Artifact `json:"artifacts"`
	// Disabled causes the holos cli to disregard the build plan.
	Disabled bool `json:"disabled,omitempty"`
}

// BuildPlanSource reflects the origin of a [BuildPlan].  Useful to save a build
// plan to a file, then re-generate it without needing to process a [Platform]
// component collection.
type BuildPlanSource struct {
	// Component reflects the component that produced the build plan.
	Component Component `json:"component,omitempty"`
}

// Artifact represents one fully rendered manifest produced by a [Transformer]
// sequence, which transforms a [Generator] collection.  A [BuildPlan] produces
// an [Artifact] collection.
//
// Each Artifact produces one manifest file artifact.  Generator Output values
// are used as Transformer Inputs.  The Output field of the final [Transformer]
// should have the same value as the Artifact field.
//
// When there is more than one [Generator] there must be at least one
// [Transformer] to combine outputs into one Artifact.  If there is a single
// Generator, it may directly produce the Artifact output.
//
// An Artifact is processed concurrently with other artifacts in the same
// [BuildPlan].  An Artifact should not use an output from another Artifact as
// an input.  Each [Generator] may also run concurrently.  Each [Transformer] is
// executed sequentially starting after all generators have completed.
//
// Output fields are write-once.  It is an error for multiple Generators or
// Transformers to produce the same Output value within the context of a
// [BuildPlan].
type Artifact struct {
	Artifact     FilePath      `json:"artifact,omitempty"`
	Generators   []Generator   `json:"generators,omitempty"`
	Transformers []Transformer `json:"transformers,omitempty"`
	Skip         bool          `json:"skip,omitempty"`
}

// Generator generates Kubernetes resources.  [Helm] and [Resources] are the
// most commonly used, often paired together to mix-in resources to an
// unmodified Helm chart.  A simple [File] generator is also available for use
// with the [Kustomize] transformer.
//
// Each Generator in an [Artifact] must have a distinct Output value for a
// [Transformer] to reference.
//
//  1. [Resources] - Generates resources from CUE code.
//  2. [Helm] - Generates rendered yaml from a [Chart].
//  3. [File] - Generates data by reading a file from the component directory.
type Generator struct {
	// Kind represents the kind of generator.  Must be Resources, Helm, or File.
	Kind string `json:"kind" cue:"\"Resources\" | \"Helm\" | \"File\""`
	// Output represents a file for a Transformer or Artifact to consume.
	Output FilePath `json:"output"`
	// Resources generator. Ignored unless kind is Resources.  Resources are
	// stored as a two level struct.  The top level key is the Kind of resource,
	// e.g. Namespace or Deployment.  The second level key is an arbitrary
	// InternalLabel.  The third level is a map[string]any representing the
	// Resource.
	Resources Resources `json:"resources,omitempty"`
	// Helm generator. Ignored unless kind is Helm.
	Helm Helm `json:"helm,omitempty"`
	// File generator. Ignored unless kind is File.
	File File `json:"file,omitempty"`
}

// Resource represents one kubernetes api object.
type Resource map[string]any

// Resources represents Kubernetes resources.  Most commonly used to mix
// resources into the [BuildPlan] generated from CUE, but may be generated from
// elsewhere.
type Resources map[Kind]map[InternalLabel]Resource

// File represents a simple single file copy [Generator].  Useful with a
// [Kustomize] [Transformer] to process plain manifest files stored in the
// component directory.  Multiple File generators may be used to transform
// multiple resources.
type File struct {
	// Source represents a file sub-path relative to the component path.
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
	// Namespace represents the helm namespace flag
	Namespace string `json:"namespace,omitempty"`
	// APIVersions represents the helm template --api-versions flag
	APIVersions []string `json:"apiVersions,omitempty"`
	// KubeVersion represents the helm template --kube-version flag
	KubeVersion string `json:"kubeVersion,omitempty"`
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

// Transformer combines multiple inputs from prior [Generator] or [Transformer]
// outputs into one output.  [Kustomize] is the most commonly used transformer.
// A simple [Join] is also supported for use with plain manifest files.
//
//  1. [Kustomize] - Patch and transform the output from prior generators or
//     transformers.  See [Introduction to Kustomize].
//  2. [Join] - Concatenate multiple prior outputs into one output.
//
// [Introduction to Kustomize]: https://kubectl.docs.kubernetes.io/guides/config_management/introduction/
type Transformer struct {
	// Kind represents the kind of transformer. Must be Kustomize, or Join.
	Kind string `json:"kind" cue:"\"Kustomize\" | \"Join\""`
	// Inputs represents the files to transform. The Output of prior Generators
	// and Transformers.
	Inputs []FilePath `json:"inputs"`
	// Output represents a file for a subsequent Transformer or Artifact to
	// consume.
	Output FilePath `json:"output"`
	// Kustomize transformer. Ignored unless kind is Kustomize.
	Kustomize Kustomize `json:"kustomize,omitempty"`
	// Join transformer. Ignored unless kind is Join.
	Join Join `json:"join,omitempty"`
}

// Join represents a [Transformer] using [bytes.Join] to concatenate multiple
// inputs into one output with a separator.  Useful for combining output from
// [Helm] and [Resources] together into one [Artifact] when [Kustomize] is
// otherwise unnecessary.
//
// [bytes.Join]: https://pkg.go.dev/bytes#Join
type Join struct {
	Separator string `json:"separator" cue:"string | *\"---\\n\""`
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

// Metadata represents data about the resource such as the Name.
type Metadata struct {
	// Name represents the resource name.
	Name string `json:"name"`
	// Labels represents a resource selector.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// Annotations represents arbitrary non-identifying metadata.  For example
	// holos uses the `cli.holos.run/description` annotation to log resources in a
	// user customized way.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Platform represents a platform to manage.  A Platform specifies a [Component]
// collection and integrates the components together into a holistic platform.
// Holos iterates over the [Component] collection producing a [BuildPlan] for
// each, which holos then executes to render manifests.
//
// Inspect a Platform resource holos would process by executing:
//
//	cue export --out yaml ./platform
type Platform struct {
	// Kind is a string value representing the resource.
	Kind string `json:"kind" cue:"\"Platform\""`
	// APIVersion represents the versioned schema of this resource.
	APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha5\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata"`

	// Spec represents the platform specification.
	Spec PlatformSpec `json:"spec"`
}

// PlatformSpec represents the platform specification.
type PlatformSpec struct {
	// Components represents a collection of holos components to manage.
	Components []Component `json:"components"`
}

// Component represents the complete context necessary to produce a [BuildPlan]
// from a path containing parameterized CUE configuration.
type Component struct {
	// Name represents the name of the component. Injected as the tag variable
	// "holos_component_name".
	Name string `json:"name"`
	// Path represents the path of the component relative to the platform root.
	// Injected as the tag variable "holos_component_path".
	Path string `json:"path"`
	// WriteTo represents the holos render component --write-to flag.  If empty,
	// the default value for the --write-to flag is used.
	WriteTo string `json:"writeTo,omitempty"`
	// Parameters represent user defined input variables to produce various
	// [BuildPlan] resources from one component path.  Injected as CUE @tag
	// variables.  Parameters with a "holos_" prefix are reserved for use by the
	// Holos Authors.  Multiple environments are a prime example of an input
	// parameter that should always be user defined, never defined by Holos.
	Parameters map[string]string `json:"parameters,omitempty"`
	// Labels represent selector labels for the component.  Copied to the
	// resulting BuildPlan.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations represents arbitrary non-identifying metadata.  Use the
	// `cli.holos.run/description` to customize the log message of each BuildPlan.
	Annotations map[string]string `json:"annotations,omitempty"`
}
