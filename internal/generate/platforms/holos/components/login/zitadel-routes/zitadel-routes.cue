package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Objects = {
	Name:      "zitadel-routes"
	Namespace: #IstioGatewaysNamespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	Resources: HTTPRoute: zitadel: {
		spec: hostnames: [_ExternalDomain, "*.\(_ExternalDomain)"]
		spec: parentRefs: [{
			name:      "default"
			namespace: #IstioGatewaysNamespace
		}]
		spec: rules: [
			{
				matches: [{path: {type: "PathPrefix", value: "/"}}]
				backendRefs: [
					{
						name:      "zitadel"
						port:      8080
						namespace: _ZitadelNamespace
					},
				]
			},
		]
	}
}
