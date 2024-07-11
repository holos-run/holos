package holos

import "encoding/yaml"

// Produce a helm chart build plan.
(#Helm & Chart).Output

_MeshConfig: (#MeshConfig & {}).Output

let Chart = {
	Name:      "istiod"
	Version:   #IstioVersion
	Namespace: "istio-system"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #IstioValues & {
		pilot: {
			// The istio meshconfig ConfigMap is handled in the holos component instead of
			// the upstream chart so extension providers can be collected from holos data.
			configMap: false
			// Set to `type: RuntimeDefault` to use the default profile if available.
			seccompProfile: type: "RuntimeDefault"
		}
	}

	Resources: ConfigMap: istio: {
		metadata: name:      "istio"
		metadata: namespace: Namespace
		data: mesh:          yaml.Marshal(_MeshConfig)
		data: meshNetworks:  "networks: {}"
	}
}
