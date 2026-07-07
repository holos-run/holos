package holos

import (
	"encoding/json"
	"github.com/holos-run/holos/api/core/v1beta1:core"
)

// Example of a simple v1beta1 command task.

holos: core.#TaskSet & {
	metadata: {
		name: "command"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command task"
	}
	spec: tasks: {
		echo: {
			kind:   "Command"
			output: "command.gen.yaml"
			command: {
				args: ["/bin/echo", json.Marshal(_ConfigMap)]
				isStdoutOutput: true
			}
		}
		// Artifact sink per schema.md D2 wiring the output to the final
		// artifact path relative to the write-to directory.
		deploy: {
			kind: "Artifact"
			inputs: ["command.gen.yaml"]
			artifact: path: "components/task/command/command.gen.yaml"
		}
	}
}

_ConfigMap: {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: name: holos.metadata.name
}
