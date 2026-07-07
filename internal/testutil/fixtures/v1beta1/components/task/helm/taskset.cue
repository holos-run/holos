package holos

import "github.com/holos-run/holos/api/core/v1beta1:core"

holos: core.#TaskSet & {
	metadata: {
		name: "helm"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) task"
	}
	spec: tasks: {
		helm: {
			kind:   "Helm"
			output: "helm.gen.yaml"
			helm: chart: {
				name:    "mychart"
				version: "0.1.0"
				release: holos.metadata.name
			}
		}
		deploy: {
			kind: "Artifact"
			inputs: ["helm.gen.yaml"]
			artifact: path: "components/task/helm/helm.gen.yaml"
		}
	}
}
