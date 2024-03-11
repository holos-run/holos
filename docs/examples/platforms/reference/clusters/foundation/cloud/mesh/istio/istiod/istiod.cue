package holos

import "encoding/yaml"

#InputKeys: component: "istiod"
#TargetNamespace: "istio-system"

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name: "istiod"
	}
	values: #IstioValues & {
		pilot: {
			// The istio meshconfig ConfigMap is handled in the holos component instead of
			// the upstream chart so extension providers can be collected from holos data.
			configMap: false
			// Set to `type: RuntimeDefault` to use the default profile if available.
			seccompProfile: type: "RuntimeDefault"
		}
	}
	apiObjects: ConfigMap: istio: #IstioConfigMap
}

#IstioConfigMap: #ConfigMap & {
	metadata: {
		name:      "istio"
		namespace: #TargetNamespace
	}
	data: {
		mesh:         yaml.Marshal(_MeshConfig)
		meshNetworks: "networks: {}"
	}
}
