package holos

// #Fleets represents the clusters in the platform.
#Fleets: {
	management: clusters: management: _
	workload: clusters: aws1:         _
}

// Platform components to manage.
#Platform: Components: {
	for Fleet in #Fleets {
		for Cluster in Fleet.clusters {
			"\(Cluster.name)/namespaces": {
				path:    "components/namespaces"
				cluster: Cluster.name
			}
		}
	}
}
