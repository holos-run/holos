package holos

import (
	"encoding/json"
	"github.com/holos-run/holos/api/core/v1alpha6:core"
)

// Example of a simple v1alpha6 command transformer.

holos: core.#BuildPlan & {
	buildContext: _
	let TMPDIR = buildContext.tempDir
	metadata: {
		name: "join"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) transformer"
	}
	spec: artifacts: [{
		artifact: "components/\(metadata.name)/\(metadata.name).gen.yaml"
		generators: [
			{
				kind:   "Command"
				output: "1.yaml"
				command: args: ["touch", "\(TMPDIR)/\(output)"]
			},
			{
				kind:   "Command"
				output: "2.yaml"
				command: args: ["touch", "\(TMPDIR)/\(output)"]
			},
			{
				kind:   "Command"
				output: "3.yaml"
				command: args: ["echo", json.Marshal(_ConfigMap)]
				command: stdout: true
			},
		]
		transformers: [
			{
				kind: "Join"
				inputs: [for g in generators {g.output}]
				output: "4.yaml"
				join: separator: "---\n"
			},
			{
				kind: "Kustomize"
				inputs: ["4.yaml"]
				output: artifact
				kustomize: kustomization: resources: inputs
			},
		]
	}]
}

_ConfigMap: {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: name: holos.metadata.name
}
