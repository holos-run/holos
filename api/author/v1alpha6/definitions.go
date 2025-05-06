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

import core "github.com/holos-run/holos/api/core/v1alpha6"

//go:generate ../../../hack/gendoc

// Platform assembles a core Platform in the Resource field for the holos render
// platform command.  Use the Components field to register components with the
// platform.
type Platform struct {
	Name       string                       `json:"name" yaml:"name" cue:"string | *\"default\""`
	Components map[NameLabel]core.Component `json:"components" yaml:"components"`
	Resource   core.Platform                `json:"resource" yaml:"resource"`
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
type ComponentConfig struct {
	// Name represents the BuildPlan metadata.name field.  Used to construct the
	// fully rendered manifest file path.
	Name string
	// Labels represent the BuildPlan metadata.labels field.
	Labels map[string]string
	// Annotations represent the BuildPlan metadata.annotations field.
	Annotations map[string]string

	// Path represents the path to the component producing the BuildPlan.
	Path string
	// Parameters are useful to reuse a component with various parameters.
	// Injected as CUE @tag variables.  Parameters with a "holos_" prefix are
	// reserved for use by the Holos Authors.
	Parameters map[string]string
	// OutputBaseDir represents the output base directory used when assembling
	// artifacts.  Useful to organize components by clusters or other parameters.
	// For example, holos writes resource manifests to
	// {WriteTo}/{OutputBaseDir}/components/{Name}/{Name}.gen.yaml
	OutputBaseDir string `cue:"string | *\"\""`

	// Resources represents kubernetes resources mixed into the rendered manifest.
	Resources core.Resources
	// KustomizeConfig represents the kustomize configuration.
	KustomizeConfig KustomizeConfig
	// Validators represent checks that must pass for output to be written.
	Validators map[NameLabel]core.Validator
	// Artifacts represents additional artifacts to mix in.  Useful for adding
	// GitOps resources.  Each Artifact is unified without modification into the
	// BuildPlan.
	Artifacts map[NameLabel]core.Artifact
}

// Helm assembles a BuildPlan rendering a helm chart.  Useful to mix in
// additional resources from CUE and transform the helm output with kustomize.
type Helm struct {
	ComponentConfig `json:",inline"`

	// Chart represents a Helm chart.
	Chart core.Chart
	// Values represents data to marshal into a values.yaml for helm.
	Values core.Values
	// ValueFiles represents value files for migration from helm value
	// hierarchies.  Use Values instead.
	ValueFiles []core.ValueFile `json:",omitempty"`
	// EnableHooks enables helm hooks when executing the `helm template` command.
	EnableHooks bool `cue:"true | *false"`
	// Namespace sets the helm chart namespace flag if provided.
	Namespace string `json:",omitempty"`
	// APIVersions represents the helm template --api-versions flag
	APIVersions []string `json:",omitempty"`
	// KubeVersion represents the helm template --kube-version flag
	KubeVersion string `json:",omitempty"`

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan core.BuildPlan
}

// Kubernetes assembles a BuildPlan containing inline resources exported from
// CUE.
type Kubernetes struct {
	ComponentConfig `json:",inline"`

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan core.BuildPlan
}

// Kustomize assembles a BuildPlan rendering manifests from a [kustomize]
// kustomization.
//
// [kustomize]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/
type Kustomize struct {
	ComponentConfig `json:",inline"`

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan core.BuildPlan
}

// KustomizeConfig represents the configuration for [kustomize] post processing.
// Use the Files field to mix in plain manifest files located in the component
// directory.  Use the Resources field to mix in manifests from network urls.
//
// [kustomize]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/
type KustomizeConfig struct {
	// Kustomization represents the kustomization used to transform resources.
	// Note the resources field is internally managed from the Files and Resources fields.
	Kustomization map[string]any `json:",omitempty"`
	// Files represents files to copy from the component directory for kustomization.
	Files map[string]struct{ Source string } `cue:"{[NAME=_]: Source: NAME}"`
	// Resources represents additional entries to included in the resources list.
	Resources map[string]struct{ Source string } `cue:"{[NAME=_]: Source: NAME}"`
	// CommonLabels represents common labels added without including selectors.
	CommonLabels map[string]string
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
type NameLabel string
