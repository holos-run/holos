package holos

// Manage on workload clusters only
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: "\(Cluster.name)/bank-frontend": {
		path:    "projects/bank-of-holos/frontend/bank-frontend"
		cluster: Cluster.name
	}
}
