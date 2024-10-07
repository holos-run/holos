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

import "google.golang.org/protobuf/types/known/structpb"

// APIObject represents the most basic generic form of a single kubernetes api
// object.  Represented as a JSON object internally for compatibility between
// tools, for example loading from CUE.
type APIObject structpb.Struct

// APIObjectMap represents the marshalled yaml representation of kubernetes api
// objects.  Do not produce an APIObjectMap directly, instead use [APIObjects]
// to produce the marshalled yaml representation from CUE data, then provide the
// result to [Component].
type APIObjectMap map[Kind]map[InternalLabel]string

// APIObjects represents Kubernetes API objects defined directly from CUE code.
// Useful to mix in resources to any kind of [Component], for example
// adding an ExternalSecret resource to a [HelmChart].
//
// [Kind] must be the resource kind, e.g. Deployment or Service.
//
// [InternalLabel] is an arbitrary internal identifier to uniquely identify the resource
// within the context of a `holos` command.  Holos will never write the
// intermediate label to rendered output.
//
// Refer to [Component] which accepts an [APIObjectMap] field provided by
// [APIObjects].
type APIObjects struct {
	APIObjects   map[Kind]map[InternalLabel]APIObject `json:"apiObjects"`
	APIObjectMap APIObjectMap                         `json:"apiObjectMap"`
}

// BuildPlan represents a build plan for the holos cli to execute.  The purpose
// of a BuildPlan is to define one or more [Component] kinds.  For example a
// [HelmChart], [KustomizeBuild], or [KubernetesObjects].
//
// A BuildPlan usually has an additional empty [KubernetesObjects] for the
// purpose of using the [Component] DeployFiles field to deploy an ArgoCD
// or Flux gitops resource for the holos component.
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

// BuildPlanSpec represents the specification of the build plan.
type BuildPlanSpec struct {
	// Disabled causes the holos cli to take no action over the [BuildPlan].
	Disabled bool `json:"disabled,omitempty"`
	// TODO: Support generators
	// TODO: Support kustomize pipeline
}

// FileContent represents file contents.
type FileContent string

// FileContentMap represents a mapping of file paths to file contents.
type FileContentMap map[FilePath]FileContent

// FilePath represents a file path.
type FilePath string

// InternalLabel is an arbitrary unique identifier internal to holos itself.
// The holos cli is expected to never write a InternalLabel value to rendered
// output files, therefore use a [InternalLabel] when the identifier must be
// unique and internal.  Defined as a type for clarity and type checking.
//
// A InternalLabel is useful to convert a CUE struct to a list, for example
// producing a list of [APIObject] resources from an [APIObjectMap].  A CUE
// struct using InternalLabel keys is guaranteed to not lose data when rendering
// output because a InternalLabel is expected to never be written to the final
// output.
type InternalLabel string

// Kind is a kubernetes api object kind. Defined as a type for clarity and type
// checking.
type Kind string

// NameLabel is a unique identifier useful to convert a CUE struct to a list
// when the values have a Name field with a default value.  This type is
// intended to indicate the common use case of converting a struct to a list
// where the Name field of the value aligns with the struct field name.
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

// PlatformComponent represents a holos component providing a BuildPlan.
type PlatformComponent struct {
	// Path is the path of the component relative to the platform root.
	Path string `json:"path"`
	// Cluster is the cluster name to provide when rendering the component.
	Cluster string `json:"cluster"`
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
	// Model represents the platform model holos gets from from the
	// PlatformService.GetPlatform rpc method and provides to CUE using a tag.
	Model structpb.Struct `json:"model" cue:"{...}"`
	// Components represents a list of holos components to manage.
	Components []PlatformComponent `json:"components"`
}
