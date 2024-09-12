package holos

// Platform wide configuration
#HTTPBin: {
	Namespace: "{{ .Namespace }}"
	Port:      80
}

// Register the namespace
#Namespaces: (#HTTPBin.Namespace): _

// Manage the component on workload clusters
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: "\(Cluster.name)/{{ .Name }}": {
		path:    "components/httpbin/workload"
		cluster: Cluster.name
	}
}
