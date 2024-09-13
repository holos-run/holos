package holos

// Manage the component on workload clusters
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: "\(Cluster.name)/{{ .Name }}": {
		path:    "components/httpbin/routes"
		cluster: Cluster.name
	}
}
