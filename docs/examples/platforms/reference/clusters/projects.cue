package holos

#Project: authProxyOrgDomain: "openinfrastructure.co"

_Projects: #Projects & {
	// The platform project is required and where platform services reside.  ArgoCD, Grafana, Prometheus, etc...
	platform: {
		resourceId: 257713952794870157
		clusters: k1: _
		clusters: k2: _
		stages: dev: authProxyClientID:  "260887327029658738@holos_platform"
		stages: prod: authProxyClientID: "260887404288738416@holos_platform"
		// Services hosted in the platform project
		hosts: argocd:     _
		hosts: grafana:    _
		hosts: prometheus: _
	}

	holos: {
		resourceId: 260446255245690199
		clusters: k1: _
		clusters: k2: _
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
