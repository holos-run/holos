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
		name: "simple"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command transformer"
	}
	spec: artifacts: [{
		artifact: "components/transformer/\(metadata.name)/\(metadata.name).gen.yaml"
		generators: [
			{
				kind:   "Command"
				output: "1.yaml"
				command: args: ["echo", json.Marshal(_ConfigMap)]
				command: isStdoutOutput: true
			},
			{
				kind:   "Command"
				output: "2.yaml"
				command: args: ["touch", "\(TMPDIR)/\(output)"]
			},
			{
				kind:   "Command"
				output: "3.yaml"
				command: args: ["touch", "\(TMPDIR)/\(output)"]
			},
		]
		transformers: [{
			kind:   "Command"
			output: artifact
			command: {
				args: ["cat", for g in generators {"\(TMPDIR)/\(g.output)"}]
				isStdoutOutput: true
			}
		}]
	}]
}

_ConfigMap: {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: name: holos.metadata.name
}
