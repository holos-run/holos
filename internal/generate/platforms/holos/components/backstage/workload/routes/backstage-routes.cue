package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "\(_Component.metadata.name)-routes"
	Namespace: #IstioGatewaysNamespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	Resources: HTTPRoute: (_Component.metadata.name): {
		spec: hostnames: [_Component.spec.hostname]
		spec: parentRefs: [{
			name:      "default"
			namespace: #IstioGatewaysNamespace
		}]
		spec: rules: [
			{
				matches: [{path: {type: "PathPrefix", value: "/"}}]
				backendRefs: [{
					name:      _Component.metadata.name
					namespace: _Component.metadata.namespace
					port:      _Component.spec.port
				}]
			},
		]
	}
}
