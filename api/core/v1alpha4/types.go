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

// APIObject represents the most basic generic form of a single kubernetes api
// object.  Represented as a JSON object internally for compatibility between
// tools, for example loading from CUE.
type APIObject map[string]any

// APIObjects represents kubernetes resources generated from CUE.
type APIObjects map[Kind]map[InternalLabel]APIObject

// HelmValues represents helm chart values generated from CUE.
type HelmValues map[string]any

// Kustomization represents a kustomization.yaml file.  Untyped to avoid tightly
// coupling holos to kubectl versions which was a problem for the Flux
// maintainers.  Type checking is expected to happen in CUE against the kubectl
// version the user prefers.
type Kustomization map[string]any

// BuildPlan represents a build plan for holos to execute.
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
	// Disabled causes the holos cli to disregard the build plan.
	Disabled bool        `json:"disabled,omitempty"`
	Steps    []BuildStep `json:"steps"`
}

type BuildStep struct {
	// Skip causes holos to skip over this build step.
	Skip         bool          `json:"skip,omitempty"`
	Generator    Generator     `json:"generator,omitempty"`
	Transformers []Transformer `json:"transformers,omitempty"`
	Paths        BuildPaths    `json:"paths"`
}

// Generator generates an artifact.
type Generator struct {
	HelmEnabled bool `json:"helmEnabled,omitempty"`
	Helm        Helm `json:"helm,omitempty"`

	KustomizeEnabled bool      `json:"kustomizeEnabled,omitempty"`
	Kustomize        Kustomize `json:"kustomize,omitempty"`

	APIObjectsEnabled bool       `json:"apiObjectsEnabled,omitempty"`
	APIObjects        APIObjects `json:"apiObjects,omitempty"`
}

type Transformer struct {
	Kind      string    `json:"kind" cue:"\"Kustomize\""`
	Kustomize Kustomize `json:"kustomize,omitempty"`
}

// Kustomize represents resources necessary to execute a kustomize build.
type Kustomize struct {
	// Kustomization represents the decoded kustomization.yaml file
	Kustomization Kustomization `json:"kustomization"`
	// Files holds file contents for kustomize, e.g. patch files.
	Files FileContentMap `json:"files,omitempty"`
}

// BuildPaths represents filesystem paths relative to the platform root.
type BuildPaths struct {
	// Component represents the component directory producing a build plan.
	Component string `json:"component"`
	// Manifest represents the directory to store fully rendered resource manifest
	// artifacts.
	Manifest string `json:"manifest,omitempty"`
	// Application represents the directory to store ArgoCD Application manifests
	// for GitOps.
	Application string `json:"application,omitempty"`
	// Flux represents the directory to store Flux Kustomization manifests
	// for GitOps.
	Flux string `json:"flux,omitempty"`
}

type Helm struct {
	// Chart represents a helm chart to manage.
	Chart Chart `json:"chart"`
	// Values represents values for holos to marshal into values.yaml when
	// rendering the chart.
	Values HelmValues `json:"values"`
	// EnableHooks enables helm hooks when executing the `helm template` command.
	EnableHooks bool `json:"enableHooks,omitempty"`
}

// Chart represents a helm chart.
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

// Repository represents a helm chart repository.
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
