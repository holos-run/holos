// Code generated by cue get go. DO NOT EDIT.

//cue:generate cue get go github.com/holos-run/holos/api/core/v1alpha4

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

// APIObject represents the most basic generic form of a single kubernetes api
// object.  Represented as a JSON object internally for compatibility between
// tools, for example loading from CUE.
#APIObject: {...}

// APIObjects represents kubernetes resources generated from CUE.
#APIObjects: {[string]: [string]: #APIObject}

// HelmValues represents helm chart values generated from CUE.
#HelmValues: {...}

// Kustomization represents a kustomization.yaml file.  Untyped to avoid tightly
// coupling holos to kubectl versions which was a problem for the Flux
// maintainers.  Type checking is expected to happen in CUE against the kubectl
// version the user prefers.
#Kustomization: {...}

// BuildPlan represents a build plan for holos to execute.
#BuildPlan: {
	// Kind represents the type of the resource.
	kind: string & "BuildPlan" @go(Kind)

	// APIVersion represents the versioned schema of the resource.
	apiVersion: string & (string | *"v1alpha4") @go(APIVersion)

	// Metadata represents data about the resource such as the Name.
	metadata: #Metadata @go(Metadata)

	// Spec specifies the desired state of the resource.
	spec: #BuildPlanSpec @go(Spec)
}

// BuildPlanSpec represents the specification of the build plan.
#BuildPlanSpec: {
	// Component represents the component that produced the build plan.
	// Represented as a path relative to the platform root.
	component: string @go(Component)

	// Disabled causes the holos cli to disregard the build plan.
	disabled?: bool @go(Disabled)

	// Steps represent build steps for holos to execute
	steps: [...#BuildStep] @go(Steps,[]BuildStep)
}

#BuildStep: {
	// Name represents the build step name, often the same as the build plan name.
	// Used to construct the output manifest and gitops filenames.
	name: string @go(Name)

	// Skip causes holos to skip over this build step.
	skip?:      bool       @go(Skip)
	generator?: #Generator @go(Generator)
	transformers?: [...#Transformer] @go(Transformers,[]Transformer)
	paths: #ArtifactPaths @go(Paths)
}

// Generator generates an artifact.
#Generator: {
	helmEnabled?: bool  @go(HelmEnabled)
	helm?:        #Helm @go(Helm)

	// HelmFile represents the intermediate file for the transformer.
	helmFile:          string & (string | *"helm.gen.yaml") @go(HelmFile)
	kustomizeEnabled?: bool                                 @go(KustomizeEnabled)
	kustomize?:        #Kustomize                           @go(Kustomize)

	// KustomizeFile represents the intermediate file for the transformer.
	kustomizeFile:      string & (string | *"kustomize.gen.yaml") @go(KustomizeFile)
	apiObjectsEnabled?: bool                                      @go(APIObjectsEnabled)
	apiObjects?:        #APIObjects                               @go(APIObjects)

	// APIObjectsFile represents the intermediate file for the transformer.
	apiObjectsFile: string & (string | *"api-objects.gen.yaml") @go(APIObjectsFile)
}

#Transformer: {
	kind:       string & "Kustomize" @go(Kind)
	kustomize?: #Kustomize           @go(Kustomize)
}

// Kustomize represents resources necessary to execute a kustomize build.
#Kustomize: {
	// Kustomization represents the decoded kustomization.yaml file
	kustomization: #Kustomization @go(Kustomization)

	// Files holds file contents for kustomize, e.g. patch files.
	files?: #FileContentMap @go(Files)
}

// ArtifactPaths represents filesystem paths relative to the write to directory
// (default is deploy/) to store artifacts.  Mainly used to specify the
// directory where resource manifests are written and a separate directory for a
// gitops resource manifest.
//
// Intended for holos to determine where to write the output of the transformer
// stage, which combines multiple generators into one stream.
#ArtifactPaths: {
	// Manifest represents the path to store fully rendered resource manifest
	// artifacts.
	manifest?: string @go(Manifest)

	// GitOps represents the path to store fully rendered gitops artifacts.  For
	// example, an ArgoCD Application or a Flux Kustomization resource.
	gitops?: string @go(Gitops)
}

#Helm: {
	// Chart represents a helm chart to manage.
	chart: #Chart @go(Chart)

	// Values represents values for holos to marshal into values.yaml when
	// rendering the chart.
	values: #HelmValues @go(Values)

	// EnableHooks enables helm hooks when executing the `helm template` command.
	enableHooks?: bool @go(EnableHooks)
}

// Chart represents a helm chart.
#Chart: {
	// Name represents the chart name.
	name: string @go(Name)

	// Version represents the chart version.
	version: string @go(Version)

	// Release represents the chart release when executing helm template.
	release: string @go(Release)

	// Repository represents the repository to fetch the chart from.
	repository?: #Repository @go(Repository)
}

// Repository represents a helm chart repository.
#Repository: {
	name: string @go(Name)
	url:  string @go(URL)
}

// FileContent represents file contents.
#FileContent: string

// FileContentMap represents a mapping of file paths to file contents.
#FileContentMap: {[string]: #FileContent}

// FilePath represents a file path.
#FilePath: string

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
#InternalLabel: string

// Kind is a kubernetes api object kind. Defined as a type for clarity and type
// checking.
#Kind: string

// NameLabel is a unique identifier useful to convert a CUE struct to a list
// when the values have a Name field with a default value.  This type is
// intended to indicate the common use case of converting a struct to a list
// where the Name field of the value aligns with the struct field name.
#NameLabel: string

// Platform represents a platform to manage.  A Platform resource informs holos
// which components to build.  The platform resource also acts as a container
// for the platform model form values provided by the PlatformService.  The
// primary use case is to collect the cluster names, cluster types, platform
// model, and holos components to build into one resource.
#Platform: {
	// Kind is a string value representing the resource.
	kind: string & "Platform" @go(Kind)

	// APIVersion represents the versioned schema of this resource.
	apiVersion: string & (string | *"v1alpha4") @go(APIVersion)

	// Metadata represents data about the resource such as the Name.
	metadata: #Metadata @go(Metadata)

	// Spec represents the specification.
	spec: #PlatformSpec @go(Spec)
}

// Metadata represents data about the resource such as the Name.
#Metadata: {
	// Name represents the resource name.
	name: string @go(Name)
}

// PlatformSpec represents the specification of a Platform.  Think of a platform
// specification as a list of platform components to apply to a list of
// kubernetes clusters combined with the user-specified Platform Model.
#PlatformSpec: {
	// Components represents a list of holos components to manage.
	components: [...#BuildContext] @go(Components,[]BuildContext)
}

// BuildContext represents the context necessary to render a component into a
// BuildPlan.  Useful to capture parameters passed down from a Platform spec for
// the purpose of idempotent rebuilds.
#BuildContext: {
	// Path is the path of the component relative to the platform root.
	path: string @go(Path)

	// Cluster is the cluster name to provide when rendering the component.
	cluster: string @go(Cluster)

	// Environment for example, dev, test, stage, prod
	environment?: string @go(Environment)

	// Model represents the platform model holos gets from from the
	// PlatformService.GetPlatform rpc method and provides to CUE using a tag.
	model: {...} @go(Model,map[string]any)

	// Tags represents cue tags to provide when rendering the component.
	tags?: [...string] @go(Tags,[]string)
}
