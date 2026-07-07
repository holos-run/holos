package holos

import "github.com/holos-run/holos/api/core/v1beta1:core"

// Example of a resources -> command -> kustomize task chain.

holos: core.#TaskSet & {
	metadata: {
		name: "kustomize"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) task chain"
	}
	spec: tasks: {
		resources: {
			kind:   "Resources"
			output: "resources.gen.yaml"
			"resources": ConfigMap: (holos.metadata.name): {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: name: holos.metadata.name
			}
		}
		copy: {
			kind: "Command"
			inputs: ["resources.gen.yaml"]
			output: "copy.gen.yaml"
			command: {
				args: ["cat"]
				stdin:          "resources.gen.yaml"
				isStdoutOutput: true
			}
		}
		transform: {
			kind: "Kustomize"
			inputs: ["copy.gen.yaml"]
			output: "kustomize.gen.yaml"
			kustomize: kustomization: resources: inputs
		}
		deploy: {
			kind: "Artifact"
			inputs: ["kustomize.gen.yaml"]
			artifact: path: "components/task/kustomize/kustomize.gen.yaml"
		}
	}
}
