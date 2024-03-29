package holos

_Projects: #Projects & {
	holos: {
		resourceId: 260446255245690199
		clusters: {
			k1: _
			k2: _
		}
		environments: {
			prod: stage: "prod"
			dev: stage:  "dev"
			jeff: stage: dev.stage
			gary: stage: dev.stage
			nate: stage: dev.stage
		}
	}

	iam: {
		clusters: {
			core1: _
			core2: _
		}
	}
}

// Manage namespaces for platform project environments.
for project in _Projects {
	for ns in project.managedNamespaces {
		if ns.clusters[#ClusterName] != _|_ {
			#ManagedNamespaces: (ns.namespace.metadata.name): ns
		}
	}
}
