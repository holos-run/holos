package v1alpha1

// Kustomization holds the rendered flux kustomization api object content for git ops.
type Kustomization struct {
	// KsContent is the yaml representation of the flux kustomization for gitops.
	KsContent string `json:"ksContent,omitempty" yaml:"ksContent,omitempty"`
	// KustomizeFiles holds file contents for kustomize, e.g. patch files.
	KustomizeFiles FileContentMap `json:"kustomizeFiles,omitempty" yaml:"kustomizeFiles,omitempty"`
	// ResourcesFile is the file name used for api objects in kustomization.yaml
	ResourcesFile string `json:"resourcesFile,omitempty" yaml:"resourcesFile,omitempty"`
}
