package holos

holos: {
	"apiVersion": "v1alpha6"
	"kind":       "Platform"
	"metadata": {
		"name": "default"
	}
	"spec": {
		"components": [
			{
				"annotations": {
					"app.holos.run/description": "slice command transformer"
				}
				"labels": {
					"holos.run/component.name": "slice"
				}
				"name": "slice"
				"parameters": {
					"outputBaseDir": "outputBaseDir"
				}
				"path": "fixtures/v1alpha6/components/slice"
			},
		]
	}
}
