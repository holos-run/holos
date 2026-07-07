package holos

import "github.com/holos-run/holos/api/core/v1beta1:core"

// Diamond shaped DAG: two independent producers join into one output, a
// command validator gates the artifact sink through an explicit dependsOn
// edge per schema.md D2.

holos: core.#TaskSet & {
	metadata: {
		name: "dag"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) diamond"
	}
	spec: tasks: {
		alfa: {
			kind:   "Resources"
			output: "alfa.gen.yaml"
			resources: ConfigMap: alfa: {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: name: "dag-alfa"
			}
		}
		bravo: {
			kind:   "Resources"
			output: "bravo.gen.yaml"
			resources: ConfigMap: bravo: {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: name: "dag-bravo"
			}
		}
		combine: {
			kind: "Join"
			inputs: ["alfa.gen.yaml", "bravo.gen.yaml"]
			output: "dag.gen.yaml"
			join: separator: "---\n"
		}
		// A command with only inputs validates, gating the sink below.
		validate: {
			kind: "Command"
			inputs: ["dag.gen.yaml"]
			command: args: ["true"]
		}
		deploy: {
			kind: "Artifact"
			inputs: ["dag.gen.yaml"]
			dependsOn: validate: {}
			artifact: path: "components/task/dag/dag.gen.yaml"
		}
	}
}
