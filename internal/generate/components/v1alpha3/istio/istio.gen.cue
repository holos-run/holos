package holos

// #Istio represents platform wide configuration
#Istio: {
	Version: "1.23.1"
	System: Namespace:  "istio-system"
	Gateway: Namespace: "istio-ingress"

	// Constrain Helm values for safer, easier upgrades and consistency across
	// platform components.
	Values: global: istioNamespace: System.Namespace

	// Configure ambient mode
	Values: profile: "ambient"
}

// Register the Namespaces
#Namespaces: (#Istio.System.Namespace):  _
#Namespaces: (#Istio.Gateway.Namespace): _

// Manage istio on workload clusters
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: {
		"\(Cluster.name)/istio-base": {
			path:    "components/istio/base"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istiod": {
			path:    "components/istio/istiod"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-cni": {
			path:    "components/istio/cni"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-ztunnel": {
			path:    "components/istio/ztunnel"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-gateway": {
			path:    "components/istio/gateway"
			cluster: Cluster.name
		}
	}
}
