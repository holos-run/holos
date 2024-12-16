package holos

// Injected from Platform.spec.components.parameters.StageName
StageName: string | *"dev" @tag(StageName)

Stages: #Stages & {
	dev:               _
	test:              _
	uat:               _
	"prod-us-east":    _
	"prod-us-central": _
	"prod-us-est":     _
}
