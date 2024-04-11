package holos

// NOTE: Beyond the base reference platform, services should typically be added to #OptionalServices instead of directly to a managed namespace.

// ManagedNamespace is a namespace to manage across all clusters in the holos platform.
#ManagedNamespace: {
	namespace: {
		metadata: {
			name: string
			labels: [string]: string
		}
	}
	// clusterNames represents the set of clusters the namespace is managed on.  Usually all clusters.
	clusterNames: [...string]
	for cluster in clusterNames {
		clusters: (cluster): name: cluster
	}
}

// #ManagedNamepsaces is the union of all namespaces across all cluster types and optional services.
// Holos adopts the namespace sameness position of SIG Multicluster, refer to https://github.com/kubernetes/community/blob/dd4c8b704ef1c9c3bfd928c6fa9234276d61ad18/sig-multicluster/namespace-sameness-position-statement.md
#ManagedNamespaces: {
	[Name=_]: #ManagedNamespace & {
		namespace: metadata: name: Name
	}
}
