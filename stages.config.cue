package holos

import holos "example.com/platform/schemas/holos/v1alpha5"

// Injected from Platform.spec.components.parameters.StageName
StageName: string | *"dev" @tag(StageName)

Stages: holos.#Stages & {
	let NONPROD = {
		tier: "nonprod"
	}
	dev: NONPROD
	test: NONPROD & {prior: dev.name}
	uat: NONPROD & {prior: test.name}

	let PROD = {
		tier:  "prod"
		prior: uat.name
	}
	"prod-us-east":    PROD
	"prod-us-central": PROD
	"prod-us-west":    PROD
}
