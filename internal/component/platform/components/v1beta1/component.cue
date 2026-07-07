package holos

// A minimal v1beta1 TaskSet exercising the Resources -> Artifact data flow.
holos: {
	metadata: name: "example"
	spec: tasks: {
		resources: {
			kind:   "Resources"
			output: "example.gen.yaml"
			resources: ConfigMap: example: {
				apiVersion: "v1"
				kind:       "ConfigMap"
				metadata: name: "example"
			}
		}
		deploy: {
			kind: "Artifact"
			inputs: ["example.gen.yaml"]
			artifact: path: "v1beta1/example/example.gen.yaml"
		}
	}
}
