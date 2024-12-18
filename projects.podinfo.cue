package holos

let IMAGE = "ghcr.io/stefanprodan/podinfo"

let PODINFO = #KargoProjectBuilder & {
	Name: "podinfo"
	// Namespaces is used as a template, KargoProjectBuilder will prefix each
	// namespace with the stage name.
	Namespaces: podinfo: _
	// Stages organized by prod and nonprod so we can easily get a handle on all
	// prod stages, for example in the HTTPRoute below.
	Stages: {
		dev: tier:               "nonprod"
		test: tier:              "nonprod"
		uat: tier:               "nonprod"
		"prod-us-east": tier:    "prod"
		"prod-us-central": tier: "prod"
		"prod-us-west": tier:    "prod"
	}
	Components: podinfo: {
		name: "podinfo"
		path: "projects/podinfo/components/podinfo"
		parameters: image: IMAGE
	}

	// Compose the Kargo promotion stages into the holos project components.
	// Project owners are expected to copy the component path into
	// projects/<project name>/components/kargo-stages and customize it as needed
	// to define their promotion process.
	Project: components: "project:\(Name):component:kargo-stages": {
		name: "kargo-stages"
		path: "components/kargo-stages"
		parameters: image:            IMAGE
		parameters: semverConstraint: "^6.0.0"
	}

	KargoProject: {
		name: _

		// We need to use a let variable otherwise name: name is a cyclical
		// reference error.
		let NAME = name
		let WAREHOUSE = {
			kind: "Warehouse"
			name: NAME
		}

		// TODO Figure out a better way to define the promotion process.  This is
		// nice and clear though.  It would be better to pull the information from
		// the stages structure though.
		promotions: {
			"dev-podinfo": requestedFreight: [{
				origin: WAREHOUSE
				sources: direct: true
			}]
			"test-podinfo": requestedFreight: [{
				origin: WAREHOUSE
				sources: stages: ["dev-podinfo"]
			}]
			"uat-podinfo": requestedFreight: [{
				origin: WAREHOUSE
				sources: stages: ["test-podinfo"]
			}]
			// We can at least aggregate all prod stages
			for STAGE in KargoProject.stages if STAGE.tier == "prod" {
				"\(STAGE.name)-podinfo": requestedFreight: [{
					origin: WAREHOUSE
					sources: stages: ["uat-podinfo"]
				}]
			}
		}
	}
}

// Register podinfo as a Holos Project
Projects: podinfo: PODINFO.Project

// Register podinfo as a Kargo Project
KargoProjects: podinfo: PODINFO.KargoProject

// Manage an HTTPRoute for the podinfo service in each namespace.
for NAMESPACE in PODINFO.Project.namespaces {
	HTTPRoutes: (NAMESPACE.metadata.name): _backendRefs: podinfo: {
		namespace: NAMESPACE.metadata.name
		port:      9898
	}
}

HTTPRoutes: podinfo: _backendRefs: {
	// Manage a backend ref for all prod tier stages.
	for COMPONENT in PODINFO.Project.components {
		if COMPONENT._stage != _|_ {
			if COMPONENT._stage.tier == "prod" {
				// The field name just needs to be unique, we don't output it.
				(COMPONENT._namespace): {
					name:      "podinfo"
					namespace: COMPONENT._namespace
					port:      9898
				}
			}
		}
	}
}
