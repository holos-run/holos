package holos

#Project: authProxyOrgDomain: "openinfrastructure.co"

let ZitadelProjectID = 257713952794870157

_Projects: #Projects & {
	// The platform project is required and where platform services reside.  ArgoCD, Grafana, Prometheus, etc...
	platform: {
		resourceId: ZitadelProjectID
		// platform level services typically run in the core cluster pair.
		clusters: core1: _
		clusters: core2: _
		// for development, probably wouldn't run these services in the workload clusters.
		clusters: k2: _
		// Services hosted in the platform project
		hosts: argocd:     _
		hosts: grafana:    _
		hosts: prometheus: _
	}

	holos: {
		resourceId: ZitadelProjectID
		clusters: k1: _
		clusters: k2: _
		environments: {
			prod: stage: "prod"
			dev: stage:  "dev"
			jeff: stage: dev.stage
			gary: stage: dev.stage
			nate: stage: dev.stage
		}
	}

	iam: {
		resourceId: ZitadelProjectID
		hosts: login: _
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
