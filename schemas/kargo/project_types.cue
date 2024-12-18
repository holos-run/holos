package kargo

import (
	holos "example.com/platform/schemas/holos/v1alpha5"
	stage "kargo.akuity.io/stage/v1alpha1"
)

// #Projects defines the structure of a kargo project, useful for kargo
// related components to look up data given a ProjectName.
#Projects: [NAME=string]: #Project & {name: NAME}

// #Project defines the structure of a project configured for progressive
// rollouts with Kargo.  The structure composes one Holos #Project with one
// Kargo #Project.
//
// We define an abstraction over both kinds of projects, holos and kargo,
// because the two are closely related but not the same.  We want to keep the
// concerns separate with loose coupling.
//
// The Holos Project is used to associate multiple related components together
// and isn't concerned with Kargo.
//
// The Kargo Project is used to define a progressive rollout promotion pipeline
// across multiple stages.  For example, Kargo handles the business use case:
// automatically promote from dev to test to uat, then submit a pull request to
// promote to production.  On merge, roll out to production across regions from
// east to west, verifying each deployment is healthy before proceeding to the
// next.
//
// See the #ProjectBuilder for how concrete values of this schema are built from
// a collection of components and stages.
#Project: {
	name:  #Name
	#Name: string

	stages: holos.#Stages

	// promotions maps the promotable component names to pipeline stages.
	promotions: [holos.#Component.#Name]: {
		requestedFreight: stage.#StageSpec.requestedFreight
	}

	// Automatically promote non-prod stages.
	promotionPolicies: [for STAGE in stages if STAGE.tier == "nonprod" {
		stage:                "\(STAGE.name)-\(name)"
		autoPromotionEnabled: true
	}]

	// Compose the holos Project into the kargo project.
	HolosProject: holos.#Project
}

// ProjectBuilder expands components out across the provided stages and
// configures a Kargo Project to manage the promotion process across stages for
// the components in the project.
#ProjectBuilder: {
	Name: string | *"default"
	// Stages to manage resources within.
	Stages: holos.#Stages
	// Namespaces to manage in each Stage.
	Namespaces: [NAME=string]: {name: NAME}
	// Components to manage in each Stage.
	Components: holos.#Components
	// BackendRefs  organized by component names.  The builder will populate the
	// namespace field.
	BackendRefs: [COMPONENT_NAME=string]: [SERVICE_NAME=string]: {
		name: string | *SERVICE_NAME
		port: int | *80
	}

	// Project represents the built kargo project.
	Project: #Project & {
		name:   Name
		stages: Stages

		// TODO: Build promotions
		promotions: {}

		for STAGE in stages {
			for COMPONENT in Components {
				let PARAMS = {
					Component: COMPONENT.name
					Prior:     STAGE.prior
				}
				promotions: (STAGE.name): requestedFreight: (#StageSpecBuilder & PARAMS).spec.requestedFreight
			}
		}

		HolosProject: {
			name: Name

			// Write all artifacts to the project specific directory.
			for KEY, COMPONENT in components {
				components: (KEY): parameters: outputBaseDir: "projects/\(Name)"
			}

			// Manage a namespace for the Kargo Project resource itself.  This namespace
			// is a container for the promotion stages.
			namespaces: (Name): {
				metadata: labels: "kargo.akuity.io/project": "true"
			}

			// Manage the component that manages the Kargo Project resource.
			components: "project:\(Name):component:kargo-project": {
				// A static name is OK because OutputBaseDir is scoped to the project.  If
				// it weren't scoped to the project, multiple projects would clobber each
				// other in the deploy directory.
				name: "kargo-project"
				path: "components/kargo-project"
				parameters: ProjectName: Name
			}

			for STAGE in Stages {
				for NAMESPACE in Namespaces {
					namespaces: "\(STAGE.name)-\(NAMESPACE.name)": _
				}

				for COMPONENT in Components {
					// Unique key to roll the component into the platform spec.
					let COMPONENT_KEY = "project:\(Name):stage:\(STAGE.name):component:\(COMPONENT.name)"

					// Generate a new component with the stage specific name and output.
					let STAGE_COMPONENT = {
						name: "\(STAGE.name)-\(COMPONENT.name)"
						for k, v in COMPONENT if k != "name" {
							(k): v
						}

						// Pass parameters to the component as tags so the component
						// definition can look up project and stage specific values.
						parameters: ProjectName:   Name
						parameters: StageName:     STAGE.name
						parameters: NamespaceName: name
						// Mix in the stage parameters
						parameters: STAGE.parameters

						// Store the stage as a hidden field so it is not output but allows us
						// to select components by stage attributes.  Useful to select all
						// prod tier components for use with an HTTPRoute.
						_stage: STAGE
						// Store the namespace name as well for reference in HTTPRoute backends.
						_namespace: name
					}
					namespaces: (STAGE_COMPONENT._namespace): _
					components: (COMPONENT_KEY):              STAGE_COMPONENT
					// Register HTTPRoutes for each stage.  The backend refs are expected to
					// be composed in from outside the builder so the builder doesn't have
					// to concern itself with service details.  We do this so we know how to
					// build a route for each stage.
					if BackendRefs[COMPONENT.name] != _|_ {
						for BACKEND_REF in BackendRefs[COMPONENT.name] {
							httpRoutes: (STAGE_COMPONENT._namespace): _backendRefs: (BACKEND_REF.name): BACKEND_REF & {
								namespace: STAGE_COMPONENT._namespace
							}
						}
					}
				}
			}
		}
	}
}
