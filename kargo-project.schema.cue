@if(!NoKargo)
package holos

import stage "kargo.akuity.io/stage/v1alpha1"

// #KargoProjects defines the structure of a kargo project, useful for kargo
// related components to look up data given a ProjectName.
#KargoProjects: [NAME=string]: #KargoProject & {name: NAME}

#KargoProject: {
	name:   string
	stages: #Stages

	// promotions maps the promotable component names to pipeline stages.
	promotions: [COMPONENT_NAME=string]: requestedFreight: stage.#StageSpec.requestedFreight

	// Automatically promote non-prod stages.
	promotionPolicies: [for STAGE in stages if STAGE.tier == "nonprod" {stage: "\(STAGE.name)-\(name)", autoPromotionEnabled: true}]
}

// KargoProjectBuilder expands components out across the provided stages and
// configures a Kargo Project to manage the promotion process across stages for
// the components in the project.
#KargoProjectBuilder: {
	Name: string | *"default"
	// Stages to manage resources within.
	Stages: #Stages
	// Namespaces to manage in each Stage.
	Namespaces: [NAME=string]: {name: NAME}
	// Components to manage in each Stage.
	Components: #Components

	Project: #Project & {
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
			}
		}
	}

	KargoProject: #KargoProject & {
		name:   Name
		stages: Stages

		for STAGE in Stages {
			for COMPONENT in Components {
				let NAME = "\(STAGE.name)-\(COMPONENT.name)"
				let PARAMS = {
					Component: COMPONENT.name
					Prior:     STAGE.prior
				}
				promotions: (NAME): requestedFreight: (#StageSpecBuilder & PARAMS).spec.requestedFreight
			}
		}
	}
}
