package holos

// A mixed platform proving per-component backwards compatibility: one
// v1alpha6 component and one v1beta1 component render in one invocation.

holos: {
	"apiVersion": "v1beta1"
	"kind":       "Platform"
	"metadata": {
		"name": "default"
	}
	"spec": {
		"components": [
			{
				"annotations": {
					"app.holos.run/description": "simple command generator"
				}
				"labels": {
					"holos.run/component.name": "simple"
				}
				"name": "simple"
				"path": "fixtures/v1beta1/components/alpha6/simple"
			},
			{
				"annotations": {
					"app.holos.run/description": "command command task"
				}
				"labels": {
					"holos.run/component.name": "command"
				}
				"name": "command"
				"path": "fixtures/v1beta1/components/task/command"
			},
		]
	}
}
