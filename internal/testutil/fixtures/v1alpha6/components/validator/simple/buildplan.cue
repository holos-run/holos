package holos

import "github.com/holos-run/holos/api/core/v1alpha6:core"

holos: core.#BuildPlan & {
	buildContext: _
	metadata: {
		name: "simple"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command validator"
	}
	spec: artifacts: [{
		artifact: "components/validator/\(metadata.name)/\(metadata.name).gen.yaml"
		generators: [{
			kind:   "Resources"
			output: artifact
			resources: ConfigMap: simple: {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: name: holos.metadata.name
			}
		}]
		validators: [{
			kind: "Command"
			inputs: [artifact]
			command: args: [
				buildContext.holosExecutable,
				"cue",
				"vet",
				"./fixtures/v1alpha6/policy",
				"--path=\"kind\"",
				"--path=strings.ToLower(kind)",
				"--path=metadata.name",
				"\(buildContext.tempDir)/\(artifact)",
			]
		}]
	}]
}
