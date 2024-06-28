package v1alpha2

// TypeMeta describes an individual object in an API response or request with
// strings representing the type of the object and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
type TypeMeta struct {
	// Kind is a string value representing the resource this object represents.
	Kind string `json:"kind" yaml:"kind"`
	// APIVersion defines the versioned schema of this representation of an object.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1alpha2\""`
}

func (tm *TypeMeta) GetKind() string {
	return tm.Kind
}

func (tm *TypeMeta) GetAPIVersion() string {
	return tm.APIVersion
}

// Discriminator discriminates the kind of an api object.
type Discriminator interface {
	// GetKind returns Kind.
	GetKind() string
	// GetAPIVersion returns APIVersion.
	GetAPIVersion() string
}

// ObjectMeta represents metadata of a holos component object. The fields are a
// copy of upstream kubernetes api machinery but are holos objects distinct from
// kubernetes api objects.
type ObjectMeta struct {
	// Name uniquely identifies the holos component instance and must be suitable as a file name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Namespace confines a holos component to a single namespace via kustomize if set.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}
