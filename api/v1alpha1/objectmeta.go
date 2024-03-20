package v1alpha1

// ObjectMeta represents metadata of a holos component object. The fields are a
// copy of upstream kubernetes api machinery but are by holos objects distinct
// from kubernetes api objects.
type ObjectMeta struct {
	// Name uniquely identifies the holos component instance and must be suitable as a file name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Namespace confines a holos component to a single namespace via kustomize if set.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// Labels are not used but are copied from api machinery ObjectMeta for completeness.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// Annotations are not used but are copied from api machinery ObjectMeta for completeness.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
