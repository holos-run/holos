package v1alpha1

// Kustomization holds the rendered flux kustomization api object content for git ops.
type Kustomization struct {
	// KsContent is the yaml representation of the flux kustomization for gitops.
	KsContent string `json:"ksContent,omitempty" yaml:"ksContent,omitempty"`
}
