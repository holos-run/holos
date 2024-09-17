package holos

// #Istio represents platform wide configuration
// Manage istio on workload clusters
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: {
		"\(Cluster.name)/istio-base": {
			path:    "projects/platform/components/istio/base"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istiod": {
			path:    "projects/platform/components/istio/istiod"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-cni": {
			path:    "projects/platform/components/istio/cni"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-ztunnel": {
			path:    "projects/platform/components/istio/ztunnel"
			cluster: Cluster.name
		}
	}
}
