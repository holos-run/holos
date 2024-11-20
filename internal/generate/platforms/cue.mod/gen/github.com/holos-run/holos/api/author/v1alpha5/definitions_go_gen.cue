// Code generated by cue get go. DO NOT EDIT.

//cue:generate cue get go github.com/holos-run/holos/api/author/v1alpha5

// Package author contains a standard set of schemas for component authors to
// generate common [core] BuildPlans.
//
// Holos values stability, flexibility, and composition.  This package
// intentionally defines only the minimal necessary set of structures.
// Component authors are encouraged to define their own structures building on
// our example [topics].
//
// The Holos Maintainers may add definitions to this package if the community
// identifies nearly all users must define the exact same structure.  Otherwise,
// definitions should be added as a customizable example in [topics].
//
// For example, structures representing a cluster and environment almost always
// need to be defined.  Their definition varies from one organization to the
// next.  Therefore, customizable definitions for a cluster and environment are
// best maintained in [topics], not standardized in this package.
//
// [core]: https://holos.run/docs/api/core/
// [topics]: https://holos.run/docs/topics/
package author

import "github.com/holos-run/holos/api/core/v1alpha5:core"

// Platform assembles a core Platform in the Resource field for the holos render
// platform command.  Use the Components field to register components with the
// platform.
#Platform: {
	Name: string
	Components: {[string]: core.#Component} @go(,map[NameLabel]core.Component)
	Resource: core.#Platform
}

// ComponentConfig represents the configuration common to all kinds of
// components for use with the holos render component command.  All component
// kinds may be transformed with [kustomize] configured with the
// [KustomizeConfig] field.
//
//   - [Helm] charts.
//   - [Kubernetes] resources generated from CUE.
//   - [Kustomize] bases.
//
// [kustomize]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/
#ComponentConfig: {
	// Name represents the BuildPlan metadata.name field.  Used to construct the
	// fully rendered manifest file path.
	Name: string

	// Labels represent the BuildPlan metadata.labels field.
	Labels: {[string]: string} @go(,map[string]string)

	// Annotations represent the BuildPlan metadata.annotations field.
	Annotations: {[string]: string} @go(,map[string]string)

	// Path represents the path to the component producing the BuildPlan.
	Path: string

	// Parameters are useful to reuse a component with various parameters.
	// Injected as CUE @tag variables.  Parameters with a "holos_" prefix are
	// reserved for use by the Holos Authors.
	Parameters: {[string]: string} @go(,map[string]string)

	// OutputBaseDir represents the output base directory used when assembling
	// artifacts.  Useful to organize components by clusters or other parameters.
	// For example, holos writes resource manifests to
	// {WriteTo}/{OutputBaseDir}/components/{Name}/{Name}.gen.yaml
	OutputBaseDir: string & (string | *"")

	// Resources represents kubernetes resources mixed into the rendered manifest.
	Resources: core.#Resources

	// KustomizeConfig represents the configuration kustomize.
	KustomizeConfig: #KustomizeConfig

	// Artifacts represents additional artifacts to mix in.  Useful for adding
	// GitOps resources.  Each Artifact is unified without modification into the
	// BuildPlan.
	Artifacts: {[string]: core.#Artifact} @go(,map[NameLabel]core.Artifact)
}

// Helm assembles a BuildPlan rendering a helm chart.  Useful to mix in
// additional resources from CUE and transform the helm output with kustomize.
#Helm: {
	#ComponentConfig

	// Chart represents a Helm chart.
	Chart: core.#Chart

	// Values represents data to marshal into a values.yaml for helm.
	Values: core.#Values

	// EnableHooks enables helm hooks when executing the `helm template` command.
	EnableHooks: bool & (true | *false)

	// Namespace sets the helm chart namespace flag if provided.
	Namespace?: string

	// APIVersions represents the helm template --api-versions flag
	APIVersions?: [...string] @go(,[]string)

	// KubeVersion represents the helm template --kube-version flag
	KubeVersion?: string

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan: core.#BuildPlan
}

// Kubernetes assembles a BuildPlan containing inline resources exported from
// CUE.
#Kubernetes: {
	#ComponentConfig

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan: core.#BuildPlan
}

// Kustomize assembles a BuildPlan rendering manifests from a [kustomize]
// kustomization.
//
// [kustomize]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/
#Kustomize: {
	#ComponentConfig

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan: core.#BuildPlan
}

// KustomizeConfig represents the configuration for [kustomize] post processing.
// Use the Files field to mix in plain manifest files located in the component
// directory.  Use the Resources field to mix in manifests from network urls.
//
// [kustomize]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/
#KustomizeConfig: {
	// Kustomization represents the kustomization used to transform resources.
	// Note the resources field is internally managed from the Files and Resources fields.
	Kustomization?: {...} @go(,map[string]any)

	// Files represents files to copy from the component directory for kustomization.
	Files: {[string]: Source: string} & {[NAME=_]: Source: NAME} @go(,map[string]struct{Source string})

	// Resources represents additional entries to included in the resources list.
	Resources: {[string]: Source: string} & {[NAME=_]: Source: NAME} @go(,map[string]struct{Source string})

	// CommonLabels represents common labels added without including selectors.
	CommonLabels: {[string]: string} @go(,map[string]string)
}

// NameLabel represents the common use case of converting a struct to a list
// where the name field of each value unifies with the field name of the outer
// struct.
//
// For example:
//
//	S: [NameLabel=string]: name: NameLabel
//	S: jeff: _
//	S: gary: _
//	S: nate: _
//	L: [for x in S {x}]
//	// L is [{name: "jeff"}, {name: "gary"}, {name: "nate"}]
#NameLabel: string
