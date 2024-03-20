package v1alpha1

// HolosComponent defines the common fields for all holos components.
type HolosComponent struct {
	TypeMeta `json:",inline" yaml:",inline"`
	// Metadata represents the holos component name
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// APIObjectMap holds the marshalled representation of api objects.
	APIObjectMap APIObjectMap `json:"apiObjectMap,omitempty" yaml:"apiObjectMap,omitempty"`
	// Kustomization holds the marshalled representation of the flux kustomization.
	Kustomization `json:",inline" yaml:",inline"`
}
