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
		let NONPROD = {
			tier: "nonprod"
			parameters: version: "6.7.0"
		}
		dev: NONPROD & {prior: "direct"}
		test: NONPROD & {prior: "dev"}
		uat: NONPROD & {prior: "test"}
		let PROD = {
			tier:  "prod"
			prior: "uat"
		}
		"prod-us-east": PROD & {parameters: version: "6.6.1"}
		"prod-us-central": PROD & {parameters: version: "6.6.2"}
		"prod-us-west": PROD & {parameters: version: "6.7.0"}
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
