package holos

// Manage on workload clusters
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: "\(Cluster.name)/httproutes": {
		path:    "projects/platform/components/httproutes"
		cluster: Cluster.name
	}
}
