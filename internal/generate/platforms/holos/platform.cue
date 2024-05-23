package holos

// _Fleets represents the clusters in the platform.
_Fleets: {
	management: clusters: management: _
	workload: clusters: aws1:         _
}

// Namespaces to manage.
_Namespaces: "holos-system": _

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
