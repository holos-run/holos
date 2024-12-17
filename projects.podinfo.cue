package holos

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
	}
}

Projects: podinfo: PODINFO.Project

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
