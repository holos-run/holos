package v1alpha3

// KustomizeBuild represents a [Component] that renders plain yaml files in
// the holos component directory using `kubectl kustomize build`.
type KustomizeBuild struct {
	Component `json:",inline"`
	Kind      string `json:"kind" cue:"\"KustomizeBuild\""`
}
