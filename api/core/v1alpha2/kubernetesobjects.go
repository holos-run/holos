package v1alpha2

const KubernetesObjectsKind = "KubernetesObjects"

// KubernetesObjects represents a [HolosComponent] composed of Kubernetes API
// objects provided directly from CUE using [APIObjects].
type KubernetesObjects struct {
	HolosComponent `json:",inline"`
	Kind           string `json:"kind" cue:"\"KubernetesObjects\""`
}
