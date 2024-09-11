package holos

// Manage podinfo on workload clusters only
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: "\(Cluster.name)/podinfo": {
		path:    "components/podinfo"
		cluster: Cluster.name
	}
}
