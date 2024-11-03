package holos

import author "github.com/holos-run/holos/api/author/v1alpha5"

// #Cluster is our definition of a cluster, building on the author api and
// composing in additional fields.
#Cluster: author.#Cluster & {
	aspects: os: "linux" | "windows"
}

// #Fleets is our definition of the fleet collection, building on the author api
// and composing in additional fields.
#Fleets: author.#StandardFleets & {
	[_]: clusters: [_]: #Cluster
}

// Manage a workload cluster named local for use with the guides.
_Fleets: #Fleets & {
	workload: clusters: local: aspects: os: "linux"
}
