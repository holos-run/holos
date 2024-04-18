package holos

// Platform level definition of a project.
#Project: {
	name: string

	// All projects have at least a prod environment and stage.
	stages: prod: stageSegments: []
	environments: prod: stage: "prod"
	environments: prod: envSegments: []
	stages: dev: _
	environments: dev: stage: "dev"
	environments: dev: envSegments: []

	// environments share the stage segments of their stage.
	environments: [_]: {
		stage:         string
		stageSegments: stages[stage].stageSegments
	}
}
