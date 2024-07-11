package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "authroutes"
	Namespace: #IstioGatewaysNamespace

	Resources: {
		Resources: [_]: [_]: metadata: namespace: Namespace
		HTTPRoute: (#HTTPRouteClone & {Name: "authproxy"}).Output
	}
}

#HTTPRouteClone: {
	Name: string
	Output: "\(Name)": {
		metadata: namespace: #IstioGatewaysNamespace
		metadata: labels: "app.kubernetes.io/name": Name
		spec: parentRefs: [{
			name:      "default"
			namespace: #IstioGatewaysNamespace
		}]
		spec: rules: [
			{
				// match the authproxy path prefix
				matches: [{path: {type: "PathPrefix", value: _AuthProxy.pathPrefix}}]
				backendRefs: [{
					name:      _AuthProxy.metadata.name
					namespace: _AuthProxy.metadata.namespace
					port:      _AuthProxy.servicePort
				}]
			},
			{
				matches: [{path: {type: "PathPrefix", value: "/holos/httpbin/"}}]
				filters: [{
					type: "URLRewrite"
					urlRewrite: {
						path: {
							type:               "ReplacePrefixMatch"
							replacePrefixMatch: "/"
						}
					}
				}]
				backendRefs: [{
					name:      _HTTPBin.metadata.name
					namespace: _HTTPBin.metadata.namespace
					port:      _HTTPBin.port
				}]
			},
		]
	}
}
