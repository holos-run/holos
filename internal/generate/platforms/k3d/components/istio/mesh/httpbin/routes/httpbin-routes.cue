package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "httpbin-routes"
	Namespace: #IstioGatewaysNamespace

	Resources: {
		Resources: [_]: [_]: metadata: namespace: Namespace
		// Multiple HTTPRoutes to test Chrome http2 connection reuse with *.admin
		// wildcard cert.
		HTTPRoute: (#HTTPRouteClone & {Name: "httpbin"}).Output
		HTTPRoute: (#HTTPRouteClone & {Name: "httpbin1"}).Output
		HTTPRoute: (#HTTPRouteClone & {Name: "httpbin2"}).Output
		HTTPRoute: (#HTTPRouteClone & {Name: "httpbin3"}).Output
		HTTPRoute: (#HTTPRouteClone & {Name: "httpbin4"}).Output
	}
}

#HTTPRouteClone: {
	Name: string
	let Host = Name + ".admin.\(_ClusterName).\(_Platform.Model.org.domain)"
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
			{
				// match the authproxy path prefix
				matches: [{path: {type: "PathPrefix", value: _AuthProxy.pathPrefix}}]
				backendRefs: [{
					name:      _AuthProxy.metadata.name
					namespace: _AuthProxy.metadata.namespace
					port:      _AuthProxy.servicePort
				}]
			},
		]
	}
}
