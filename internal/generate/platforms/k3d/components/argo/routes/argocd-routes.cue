package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "argo-routes"
	Namespace: #IstioGatewaysNamespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	Resources: HTTPRoute: argocd: {
		spec: hostnames: [_ArgoCD.hostname]
		spec: parentRefs: [{
			name:      "default"
			namespace: #IstioGatewaysNamespace
		}]
		spec: rules: [
			{
				matches: [{path: {type: "PathPrefix", value: "/"}}]
				backendRefs: [{
					name:      "argocd-server"
					port:      80
					namespace: _ArgoCD.metadata.namespace
				}]
			},
		]
	}
}
