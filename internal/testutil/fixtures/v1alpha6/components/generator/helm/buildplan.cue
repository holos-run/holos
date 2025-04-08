package holos

import "github.com/holos-run/holos/api/core/v1alpha6:core"

holos: core.#BuildPlan & {
	metadata: {
		name: "helm"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) generator"
	}
	spec: artifacts: [{
		artifact: "components/generator/\(metadata.name)/\(metadata.name).gen.yaml"
		generators: [{
			kind:   "Helm"
			output: artifact
			helm: chart: {
				name:    "mychart"
				version: "0.1.0"
				release: name
			}
		}]
	}]
}
