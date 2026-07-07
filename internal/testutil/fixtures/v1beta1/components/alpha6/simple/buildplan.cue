package holos

import (
	"encoding/json"
	"github.com/holos-run/holos/api/core/v1alpha6:core"
)

// A v1alpha6 component colocated in the v1beta1 fixture tree so the mixed
// platform1 fixture proves one platform renders components with different
// api versions in one invocation.

holos: core.#BuildPlan & {
	metadata: {
		name: "simple"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) command generator"
	}
	spec: artifacts: [{
		artifact: "components/alpha6/\(metadata.name)/\(metadata.name).gen.yaml"
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
