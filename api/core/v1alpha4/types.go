// # Core API
//
// Package v1alpha4 contains the core API contract between the holos cli and CUE
// configuration code.  Platform designers, operators, and software developers
// use this API to write configuration in CUE which holos loads.  The Core API
// is declarative.  Each resource represents a desired state necessary for holos
// to fully render Kubernetes manifests into plain files.
//
// The following resources provide important context for the Core API.  The
// [Author API] is intended for component authors as a convenient adapter for
// the Core API resources Holos expects.
//
//  1. [Technical Overview]
//  2. [Quickstart]
//  3. [Author API]
//
// # Platform
//
// [Platform] defines the complete configuration of a platform.  A platform
// represents a [Component] collection.
//
// Inspect a Platform resource holos would process by executing:
//
//	cue export --out yaml ./platform
//
// # Component
//
// A [Component] is the combination of CUE code along one path relative to the
// platform root directory plus data injected from the [PlatformSpec] via CUE tags.
// The platform configuration root is the directory containing cue.mod.
//
// A [Component] always produces exactly one [BuildPlan].
//
// # BuildPlan
//
// A [BuildPlan] contains an [Artifact] collection.  A BuildPlan often produces
// two artifacts, one containing the fully rendered Kubernetes API resources,
// the other containing an additional resource to manage the former with GitOps.
// For example, a BuildPlan for a podinfo component produces a manifest
// containing a Deployment and a Service, along with a second manifest
// containing an ArgoCD Application.
//
// Inspect a BuildPlan resource holos render component would process by executing:
//
//	cue export --out yaml ./projects/platform/components/namespaces
//
// # Artifact
//
// An [Artifact] is one fully rendered manifest file produced from the final
// [Transformer] in a sequence of transformers.  An Artifact may also be
// produced directly from a [Generator], but this use case is uncommon.
//
// # Transformer
//
// A [Transformer] takes multiple inputs from prior [Generator] or [Transformer]
// outputs, then transforms the data into one output.  [Kustomize] is the most
// commonly used transformer, though a simple [Join] is also supported.
//
//  1. [Kustomize] - Patch and transform the output from prior generators or
//     transformers.  See [Introduction to Kustomize].
//  2. [Join] - Concatenate multiple prior outputs into one output.
//
// # Generators
//
// A [Generator] generates Kubernetes resources.  [Helm] and [Resources] are the
// most commonly used, often paired together to mix-in resources to an
// unmodified Helm chart.  A simple [File] generator is also available for use
// with the [Kustomize] transformer.
//
//  1. [Resources] - Generates resources from CUE code.
//  2. [Helm] - Generates rendered yaml from a [Chart].
//  3. [File] - Generates data by reading a file from the component directory.
//
// [Introduction to Kustomize]: https://kubectl.docs.kubernetes.io/guides/config_management/introduction/
// [Author API]: https://holos.run/docs/api/author/
// [Quickstart]: https://holos.run/docs/quickstart/
// [Technical Overview]: https://holos.run/docs/technical-overview/
package v1alpha4

//go:generate ../../../hack/gendoc

// BuildPlan represents a build plan for holos to execute.  Each [Platform]
// component produces exactly one BuildPlan.
//
// One or more [Artifact] files are produced by a BuildPlan, representing the
// fully rendered manifests for the Kubernetes API Server.
//
// # Example BuildPlan
//
// Command:
//
//	cue export --out yaml ./projects/platform/components/namespaces
//
// Output:
//
//	kind: BuildPlan
//	apiVersion: v1alpha4
//	metadata:
//	  name: dev-namespaces
//	spec:
//	  component: projects/platform/components/namespaces
//	  artifacts:
//	    - artifact: clusters/no-cluster/components/dev-namespaces/dev-namespaces.gen.yaml
//	      generators:
//	        - kind: Resources
//	          output: resources.gen.yaml
//	          resources:
//	            Namespace:
//	              dev-jeff:
//	                metadata:
//	                  name: dev-jeff
//	                  labels:
//	                    kubernetes.io/metadata.name: dev-jeff
//	                kind: Namespace
//	                apiVersion: v1
//	      transformers:
//	        - kind: Kustomize
//	          inputs:
//	            - resources.gen.yaml
//	          output: clusters/no-cluster/components/dev-namespaces/dev-namespaces.gen.yaml
//	          kustomize:
//	            kustomization:
//	              commonLabels:
//	                holos.run/component.name: dev-namespaces
//	              resources:
//	                - resources.gen.yaml
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
	// Artifacts represents the artifacts for holos to build.
	Artifacts []Artifact `json:"artifacts"`
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

// Generator generates an intermediate manifest for a [Artifact].
//
// Each Generator in a [Artifact] must have a distinct Output value for a
// [Transformer] to reference.
//
// Refer to [Resources], [Helm], and [File].
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

// Resources represents a kubernetes resources [Generator] from CUE.
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

// Transformer transforms [Generator] manifests within a [Artifact].
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

// PlatformSpec represents the specification of a [Platform].  Think of a
// platform spec as a [Component] collection for multiple kubernetes clusters
// combined with the user-specified Platform Model.
type PlatformSpec struct {
	// Components represents a list of holos components to manage.
	Components []Component `json:"components"`
}

// Component represents the complete context necessary to produce a [BuildPlan].
// Component carries information injected from holos render platform to holos
// render component to produce each [BuildPlan].
//
// All of these fields are passed to the holos render component command using
// flags, which in turn are injected to CUE using tags.  For clarity, CUE field
// and tag names should match the struct json tag names below.
type Component struct {
	// Name represents the name of the component, injected as a tag to set the
	// BuildPlan metadata.name field.  Necessary for clear user feedback during
	// platform rendering.
	Name string `json:"name"`
	// Component represents the path of the component relative to the platform root.
	Component string `json:"component"`
	// Cluster is the cluster name to provide when rendering the component.
	Cluster string `json:"cluster"`
	// Environment for example, dev, test, stage, prod
	Environment string `json:"environment,omitempty"`
	// Model represents the platform model holos gets from from the
	// PlatformService.GetPlatform rpc method and provides to CUE using a tag.
	Model map[string]any `json:"model"`
	// Tags represents cue tags to inject when rendering the component.  The json
	// struct tag names of other fields in this struct are reserved tag names not
	// to be used in the tags collection.
	Tags []string `json:"tags,omitempty"`
	// WriteTo represents the holos render component --write-to flag.  If empty,
	// the default value for the --write-to flag is used.
	WriteTo string `json:"writeTo,omitempty"`
}

// Tags represents standardized fields injected into the component [BuildPlan]
// from the [Platform].
//
// Note, tags should have a reasonable default value to easily use cue eval and
// cue export without needing to make a bunch of decisions about tag values.
//
// Example:
//
//	import core "github.com/holos-run/holos/api/core/v1alpha4"
//	_Tags: core.#Tags & {
//	  cluster:     _ @tag(cluster, type=string)
//	  environment: _ @tag(environment, type=string)
//	  component:   _ @tag(component, type=string)
//	  name:        _ @tag(name, type=string)
//	}
type Tags struct {
	// Name represents the BuildPlan metadata.name field injected from the Platform.
	Name string `json:"name" cue:"string | *\"no-name\""`
	// Cluster represents the cluster name injected from
	Cluster string `json:"cluster" cue:"string | *\"no-cluster\""`
	// Environment represents the build plan environment.
	Environment string `json:"environment" cue:"string | *\"no-environment\""`
	// Component represents the path of the component relative to the platform root.
	Component string `json:"component" cue:"string | *\"no-component\""`
}
