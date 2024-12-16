@if(!NoIstio)
package holos

Projects: network: {
	namespaces: (Istio.System.Namespace):  _
	namespaces: (Istio.Gateway.Namespace): _

	components: {
		"istio-base": {
			name: "istio-base"
			path: "projects/network/components/istio-base"
		}
		"istiod": {
			name: "istiod"
			path: "projects/network/components/istiod"
		}
		"istio-cni": {
			name: "istio-cni"
			path: "projects/network/components/istio-cni"
		}
		"istio-ztunnel": {
			name: "istio-ztunnel"
			path: "projects/network/components/istio-ztunnel"
		}
		"istio-gateway": {
			name: "istio-gateway"
			path: "projects/network/components/istio-gateway"
		}
		"httproutes": {
			name: "httproutes"
			path: "projects/network/components/httproutes"
		}
	}
}

#Istio: {
	Version: string
	System: Namespace:  string
	Gateway: Namespace: string
	Values: {...}
}

// #Istio represents platform wide configuration
Istio: #Istio & {
	Version: "1.23.1"
	System: Namespace:  "istio-system"
	Gateway: Namespace: "istio-ingress"

	// Constrain Helm values for safer, easier upgrades and consistency across
	// platform components.
	Values: global: istioNamespace: System.Namespace
	// Configure ambient mode
	Values: profile: "ambient"
}
