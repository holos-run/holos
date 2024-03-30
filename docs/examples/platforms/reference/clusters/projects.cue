package holos

#Project: authProxyOrgDomain: "openinfrastructure.co"

_Projects: #Projects & {
	holos: {
		resourceId: 260446255245690199
		clusters: {
			k1: _
			k2: _
		}
		stages: dev: authProxyClientID:  "260505543108527218@holos"
		stages: prod: authProxyClientID: "260506079325128023@holos"
		environments: {
			prod: stage: "prod"
			dev: stage:  "dev"
			jeff: stage: dev.stage
			gary: stage: dev.stage
			nate: stage: dev.stage
		}
	}

	iam: {
		resourceId: 260582480954787159
		clusters: {
			core1: _
			core2: _
		}
		stages: dev: authProxyClientID:  "260582521186616432@iam"
		stages: prod: authProxyClientID: "260582633862399090@iam"
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
