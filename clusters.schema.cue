package holos

import "strings"

// #Cluster represents one cluster
#Cluster: {
	// name represents the cluster name.
	name: string & =~"[a-z][a-z0-9]+" & strings.MinRunes(2) & strings.MaxRunes(63)
	// Constrain the regions.  No default, the region must be specified.
	region: "us-east1" | "us-central1" | "us-west1" | "local"
	// Each cluster must be in only one set of clusters.  All but one cluster are
	// workload clusters, so make it the default.
	set: "management" | *"workload"
}

// #Clusters represents a cluster collection structure
#Clusters: {
	// name is the lookup key for the collection.
	[NAME=string]: #Cluster & {
		// name must match the struct field name.
		name: NAME
	}
}

// #ClusterSet represents a set of clusters.
#ClusterSet: {
	// name represents the cluster set name.
	name: string & =~"[a-z][a-z0-9]+" & strings.MinRunes(2) & strings.MaxRunes(63)
	clusters: #Clusters & {
		// Constrain the cluster set to clusters having the same set.  Ensures
		// clusters are never mis-categorized.
		[_]: set: name
	}
}

// #ClusterSets represents a cluster set collection.
#ClusterSets: {
	// name is the lookup key for the collection.
	[NAME=string]: #ClusterSet & {
		// name must match the struct field name.
		name: NAME
	}
}
