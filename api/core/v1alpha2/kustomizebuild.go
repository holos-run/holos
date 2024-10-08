package v1alpha2

// KustomizeBuild represents a [HolosComponent] that renders plain yaml files in
// the holos component directory using `kubectl kustomize build`.
type KustomizeBuild struct {
	HolosComponent `json:",inline"`
	Kind           string `json:"kind" cue:"\"KustomizeBuild\""`
}
