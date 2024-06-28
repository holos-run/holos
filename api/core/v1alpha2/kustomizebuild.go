package v1alpha2

// KustomizeBuild renders plain yaml files in the holos component directory
// using kubectl kustomize build.
type KustomizeBuild struct {
	HolosComponent `json:",inline" yaml:",inline"`
	Kind           string `json:"kind" yaml:"kind" cue:"\"KustomizeBuild\""`
}
