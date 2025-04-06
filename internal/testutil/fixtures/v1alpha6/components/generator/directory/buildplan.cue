package holos

import (
	"encoding/json"
	"github.com/holos-run/holos/api/core/v1alpha6:core"
)

// Example of a simple v1alpha6 command generator.

holos: core.#BuildPlan & {
	metadata: {
		name: "directory"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command generator"
	}
	spec: artifacts: [{
		// The user should be able to specify a directory as the artifact, holos
		// copies everything in it.
		artifact: "components/\(metadata.name)"
		generators: [{
			kind:   "Command"
			output: "\(artifact)/\(metadata.name).gen.yaml"
			command: {
				args: ["/bin/echo", json.Marshal(_ConfigMap)]
				stdout: true
			}
		}]
	}]
}

_ConfigMap: {
	apiVersion: "v1"
	kind:       "ConfigMap"
	metadata: name: holos.metadata.name
}
