package holos

import holos "example.com/platform/schemas/holos/v1alpha5"

// Injected from Platform.spec.components.parameters.StageName
StageName: string | *"dev" @tag(StageName)

Stages: holos.#Stages & {
	dev:               _
	test:              _
	uat:               _
	"prod-us-east":    _
	"prod-us-central": _
	"prod-us-est":     _
}
