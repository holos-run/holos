package holos

// Platform level definition of a project.
#Project: {
	name: string

	// All projects have at least a prod and dev environment and stage.

	// Omit the prod stage segment from hostnames.  foo.holos.run not foo.prod.holos.run
	stages: prod: stageSegments: []
	environments: prod: stage: "prod"
	// Omit the prod env segment from hostnames.  foo.holos.run not prod.foo.holos.run
	environments: prod: envSegments: []

	stages: dev: _
	environments: dev: stage: "dev"
	// Omit the dev env segment from hostnames.  foo.dev.holos.run not dev.foo.dev.holos.run
	environments: dev: envSegments: []

	// environments share the stage segments of their stage.
	environments: [_]: {
		stage:         string
		stageSegments: stages[stage].stageSegments
	}
}
