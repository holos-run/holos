@if(!NoKargo)
package holos

// KargoProjectBuilder expands components out across the provided stages and
// configured a Kargo Project resource to manage the promotion process across
// stages for the project.
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

		// TODO: Manage the component that manages the Kargo Stage resources.

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
}
