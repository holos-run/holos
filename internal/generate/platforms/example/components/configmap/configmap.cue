package holos

import "encoding/yaml"
import v1 "github.com/holos-run/holos/api/v1alpha1"

let PLATFORM = {message: "TODO: Load the platform from the API."}

// Provide a BuildPlan to the holos cli to render k8s api objects.
v1.#BuildPlan & {
	spec: components: resources: platformConfigmap: {
		metadata: name: "platform-configmap"
		apiObjectMap: OBJECTS.apiObjectMap
	}
}

// OBJECTS represents the kubernetes api objects to manage.
let OBJECTS = v1.#APIObjects & {
	apiObjects: ConfigMap: platform: {
		metadata: {
			name:      "platform"
			namespace: "default"
		}
		data: platform: yaml.Marshal(PLATFORM)
	}
}
