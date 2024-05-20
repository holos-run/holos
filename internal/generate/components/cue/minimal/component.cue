package holos

import v1 "github.com/holos-run/holos/api/v1alpha1"
import "encoding/yaml"

let ComponentName = "{{ .ComponentName }}"

// The BuildPlan represents the kubernetes api objects to manage.  CUE returns
// the build plan to the holos CLI for rendering to plain yaml files.
v1.#BuildPlan & {
	spec: components: resources: "\(ComponentName)": {
		metadata: name: ComponentName
		apiObjectMap: OBJECTS.apiObjectMap
	}
}

// OBJECTS represents the kubernetes api objects to manage.
let OBJECTS = v1.#APIObjects & {
  // Add Kubernetes API Objects to manage here.
	apiObjects: ConfigMap: "\(ComponentName)": {
		metadata: {
			name:      ComponentName
			namespace: "default"
		}
		data: platform: yaml.Marshal(PLATFORM)
	}
}

// This is an example of how to refer to the Platform model.
let PLATFORM = {
	spec: model: _Platform.spec.model
}
