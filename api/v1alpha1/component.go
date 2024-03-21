package v1alpha1

// HolosComponent defines the fields common to all holos component kinds including the Render Result.
type HolosComponent struct {
	TypeMeta `json:",inline" yaml:",inline"`
	// Metadata represents the holos component name
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// APIObjectMap holds the marshalled representation of api objects. Think of
	// these as resources overlaid at the back of the render pipeline.
	APIObjectMap APIObjectMap `json:"apiObjectMap,omitempty" yaml:"apiObjectMap,omitempty"`
	// Kustomization holds the marshalled representation of the flux kustomization
	// which reconciles resources in git with the api server.
	Kustomization `json:",inline" yaml:",inline"`
	// Kustomize represents a kubectl kustomize build post-processing step.
	Kustomize `json:",inline" yaml:",inline"`
	// Skip causes holos to take no action regarding the component.
	Skip bool
}

func (hc *HolosComponent) NewResult() *Result {
	return &Result{HolosComponent: *hc}
}
