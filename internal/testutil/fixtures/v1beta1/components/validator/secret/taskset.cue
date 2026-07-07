package holos

import "github.com/holos-run/holos/api/core/v1beta1:core"

holos: core.#TaskSet & {
	buildContext: _
	metadata: {
		name: "secret"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command validator"
	}
	spec: tasks: {
		resources: {
			kind:   "Resources"
			output: "secret.gen.yaml"
			"resources": Secret: (holos.metadata.name): {
				apiVersion: "v1"
				kind:       "Secret"
				metadata: name: holos.metadata.name
			}
		}
		validate: {
			kind: "Command"
			inputs: ["secret.gen.yaml"]
			command: args: [
				buildContext.holosExecutable,
				"cue",
				"vet",
				"./fixtures/v1beta1/policy",
				"--path=\"kind\"",
				"--path=strings.ToLower(kind)",
				"--path=metadata.name",
				"\(buildContext.tempDir)/secret.gen.yaml",
			]
		}
		deploy: {
			kind: "Artifact"
			inputs: ["secret.gen.yaml"]
			dependsOn: validate: {}
			artifact: path: "components/validator/secret/secret.gen.yaml"
		}
	}
}
