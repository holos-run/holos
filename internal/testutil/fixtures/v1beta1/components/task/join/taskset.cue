package holos

import (
	"encoding/json"
	"github.com/holos-run/holos/api/core/v1beta1:core"
)

// Example of two producer tasks joined into one output.

holos: core.#TaskSet & {
	metadata: {
		name: "join"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) task"
	}
	spec: tasks: {
		alfa: {
			kind:   "Command"
			output: "alfa.gen.yaml"
			command: {
				args: ["/bin/echo", json.Marshal(_ConfigMap & {metadata: name: "join-alfa"})]
				isStdoutOutput: true
			}
		}
		bravo: {
			kind:   "Command"
			output: "bravo.gen.yaml"
			command: {
				args: ["/bin/echo", json.Marshal(_ConfigMap & {metadata: name: "join-bravo"})]
				isStdoutOutput: true
			}
		}
		combine: {
			kind: "Join"
			inputs: ["alfa.gen.yaml", "bravo.gen.yaml"]
			output: "join.gen.yaml"
			join: separator: "---\n"
		}
		deploy: {
			kind: "Artifact"
			inputs: ["join.gen.yaml"]
			artifact: path: "components/task/join/join.gen.yaml"
		}
	}
}

_ConfigMap: {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: name: string
}
