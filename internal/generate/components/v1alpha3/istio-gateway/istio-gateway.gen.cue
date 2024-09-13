package holos

// #Istio represents platform wide configuration
#Istio: Gateway: Namespace: "istio-ingress"

// Register the Namespaces
#Namespaces: (#Istio.Gateway.Namespace): _

// Manage istio on workload clusters
for Cluster in #Fleets.workload.clusters {
	#Platform: Components: {
		"\(Cluster.name)/{{ .Name }}": {
			path:    "components/istio/gateway"
			cluster: Cluster.name
		}
	}
}
