package v1alpha3

// Component defines the fields common to all holos component kinds.  Every
// holos component kind should embed Component.
type Component struct {
	// Kind is a string value representing the resource this object represents.
	Kind string `json:"kind"`
	// APIVersion represents the versioned schema of this representation of an object.
	APIVersion string `json:"apiVersion" cue:"\"v1alpha3\""`
	// Metadata represents data about the holos component such as the Name.
	Metadata Metadata `json:"metadata"`

	// APIObjectMap holds the marshalled representation of api objects.  Useful to
	// mix in resources to each Component type, for example adding an
	// ExternalSecret to a [HelmChart] Component.  Refer to [APIObjects].
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

// Metadata represents data about the object such as the Name.
type Metadata struct {
	// Name represents the name of the holos component.
	Name string `json:"name"`
	// Namespace is the primary namespace of the holos component.  A holos
	// component may manage resources in multiple namespaces, in this case
	// consider setting the component namespace to default.
	//
	// This field is optional because not all resources require a namespace,
	// particularly CRDs and DeployFiles functionality.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}
