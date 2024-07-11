package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "httpbin-routes"
	Namespace: #IstioGatewaysNamespace

	Resources: {
		Resources: [_]: [_]: metadata: namespace: Namespace
		HTTPRoute: (#HTTPRouteClone & {Name: "httpbin"}).Output
	}
}

#HTTPRouteClone: {
	Name: string
	let Host = Name + ".\(_Platform.Model.org.domain)"
	Output: "\(Name)": {
		metadata: namespace: #IstioGatewaysNamespace
		metadata: labels: app: "httpbin"
		spec: hostnames: [Host]
		spec: parentRefs: [{
			name:      "default"
			namespace: #IstioGatewaysNamespace
		}]
		spec: rules: [
			{
				matches: [{path: {type: "PathPrefix", value: "/"}}]
				backendRefs: [{
					name:      _HTTPBin.metadata.name
					namespace: _HTTPBin.metadata.namespace
					port:      _HTTPBin.port
				}]
			},
		]
	}
}
