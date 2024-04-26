package holos

import ( "encoding/yaml"

)

// The platform configmap is a simple component that manages a configmap named
// platform in the default namespace.  The purpose is to exercise end to end
// validation of platform configuration values provided by the holos web ui to
// each cluster in the platform.
platform: #Platform & {metadata: name: "bare"}
let PLATFORM = platform

// spec represents the output provided to holos
spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "platform-configmap"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

// OBJECTS represents the kubernetes api objects to manage.
let OBJECTS = #APIObjects & {
	apiObjects: ConfigMap: platform: {
		metadata: {
			name:      "platform"
			namespace: "default"
		}
		data: platform: yaml.Marshal(PLATFORM)
	}
}
