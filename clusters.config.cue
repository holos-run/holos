package holos

// For this demo we'll configure a single management cluster.
Clusters: #Clusters & {
	management: region: "local"
	management: set:    "management"
}

// ClusterSets is dynamically built from the Clusters structure.
ClusterSets: #ClusterSets & {
	// Map every cluster into the correct set.
	for CLUSTER in Clusters {
		(CLUSTER.set): clusters: (CLUSTER.name): CLUSTER
	}
}
