package v1alpha2

const KubernetesObjectsKind = "KubernetesObjects"

// KubernetesObjects represents a holos component composed of kubernetes api
// objects provided directly from CUE.
type KubernetesObjects struct {
	HolosComponent `json:",inline" yaml:",inline"`
	Kind           string `json:"kind" yaml:"kind" cue:"\"KubernetesObjects\""`
}
