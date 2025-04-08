package holos

import (
	"encoding/json"
	"github.com/holos-run/holos/api/core/v1alpha6:core"
)

// Example of a simple v1alpha6 command generator.

holos: core.#BuildPlan & {
	metadata: {
		name: "simple"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command generator"
	}
	spec: artifacts: [{
		artifact: "components/generator/\(metadata.name)/\(metadata.name).gen.yaml"
		generators: [{
			kind:   "Command"
			output: artifact
			command: {
				args: ["/bin/echo", json.Marshal(_ConfigMap)]
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
