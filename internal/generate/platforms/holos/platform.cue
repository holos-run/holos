package holos

// _Fleets represents the clusters in the platform.
_Fleets: {
	management: clusters: management: _
	workload: clusters: aws1:         _
}

// Namespaces to manage.
_Namespaces: "holos-system": _

// Include all project namespaces in the platform namespaces.
for project in _Projects {
	_Namespaces: project.spec.namespaces
}

// Projects to manage.
_Projects: holos: spec: namespaces: "holos-system": _
_Projects: argocd: spec: namespaces: argocd:        _

// Platform components to manage.
_Platform: Components: {
	for Fleet in _Fleets {
		for Cluster in Fleet.clusters {
			"\(Cluster.name)/namespaces": {
				path:    "components/namespaces"
				cluster: Cluster.name
			}
		}
	}
}
