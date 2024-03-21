package holos

import "encoding/yaml"

#InputKeys: component: "istiod"
#TargetNamespace: "istio-system"

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-namespaces":       _
		_dependsOn: "\(#InstancePrefix)-istio-base": _

		metadata: name: "prod-mesh-istiod"
		chart: name:    "istiod"
		namespace: #TargetNamespace
		_values: #IstioValues & {
			pilot: {
				// The istio meshconfig ConfigMap is handled in the holos component instead of
				// the upstream chart so extension providers can be collected from holos data.
				configMap: false
				// Set to `type: RuntimeDefault` to use the default profile if available.
				seccompProfile: type: "RuntimeDefault"
			}
		}
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {apiObjects: ConfigMap: istio: #IstioConfigMap}

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
