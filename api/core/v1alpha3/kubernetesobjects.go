package v1alpha3

const KubernetesObjectsKind = "KubernetesObjects"

// KubernetesObjects represents a [Component] composed of Kubernetes API
// objects provided directly from CUE using [APIObjects].
type KubernetesObjects struct {
	Component `json:",inline"`
	Kind      string `json:"kind" cue:"\"KubernetesObjects\""`
}
