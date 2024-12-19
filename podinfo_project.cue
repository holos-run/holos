package holos

import "example.com/platform/schemas/kargo"

// let IMAGE = "ghcr.io/stefanprodan/podinfo"
let IMAGE = "quay.io/holos/stefanprodan/podinfo"
let HTTPROUTE_LABEL = "holos.run/httproute.project"

let PROJECT = kargo.#ProjectBuilder & {
	Name: "podinfo"
	// Namespaces are used as a template, KargoProjectBuilder will prefix each
	// namespace with the stage name.
	Namespaces: (Name): metadata: labels: (HTTPROUTE_LABEL): Name

	// Stages organized by prod and nonprod so we can easily get a handle on all
	// prod stages, for example in the HTTPRoute below.
	Stages: {
		let PARAMS = {
			name: string
			parameters: message: "Hello! Stage: \(name)"
		}
		let NONPROD = PARAMS & {
			tier: "nonprod"
			parameters: version: "6.7.0"
		}
		dev: NONPROD & {prior: "direct"}
		test: NONPROD & {prior: "dev"}
		uat: NONPROD & {prior: "test"}
		let PROD = PARAMS & {
			tier:  "prod"
			prior: "uat"
			// We have to stringify all injected tags.  This is a reason to switch to
			// passing the component over standard input as structured data.
			parameters: replicaCount: "2"
		}
		"prod-us-east": PROD & {
			parameters: version: "6.6.1"
		}
		"prod-us-central": PROD & {
			parameters: version: "6.6.2"
		}
		"prod-us-west": PROD & {
			parameters: version: "6.7.0"
		}
	}

	Components: (Name): {
		name: Name
		path: "projects/\(Name)/components/\(Name)"
		parameters: image: IMAGE
	}

	// Compose the Kargo promotion stages into the holos project components.
	// Project owners are expected to copy the component path into
	// projects/<project name>/components/kargo-stages and customize it as needed
	// to define their promotion process.
	Project: HolosProject: components: "project:\(Name):component:kargo-stages": {
		name: "kargo-stages"
		path: "components/kargo-stages"
		parameters: image:            IMAGE
		parameters: semverConstraint: "^6.0.0"
	}
}

// Register the project as a Holos Project.
Projects: (PROJECT.Name): PROJECT.Project.HolosProject

// Register the project as a Kargo Project.
KargoProjects: (PROJECT.Name): PROJECT.Project

// Compose stage specific httproutes with the platform selecting namespaces.
for NS in PROJECT.Project.HolosProject.namespaces {
	for K, V in NS.metadata.labels {
		if K == HTTPROUTE_LABEL && V == PROJECT.Name {
			HTTPRoutes: (NS.metadata.name): _backendRefs: (NS.metadata.name): {
				name:      PROJECT.Name
				namespace: NS.metadata.name
				port:      9898
			}
		}

		// Manage a backend ref for all prod tier stages.
		if K == "holos.run/stage.tier" && V == "prod" {
			HTTPRoutes: (PROJECT.Name): _backendRefs: (NS.metadata.name): {
				name:      PROJECT.Name
				namespace: NS.metadata.name
				port:      9898
			}
		}
	}
}
