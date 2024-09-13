package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Objects = {
	Name:      "{{ .Name }}"
	Namespace: #Istio.Gateway.Namespace

	Resources: [_]: [_]: metadata: namespace: Namespace
	Resources: HTTPRoute: (#HTTPRouteClone & {Name: "httpbin"}).BuildPlan
}

#HTTPRouteClone: {
	Name: string
	let Host = "\(Name).\(#Platform.Domain)"
	Output: "\(Name)": {
		metadata: namespace: _
		metadata: name:      Name
		metadata: labels: app: Name
		spec: hostnames: [Host]
		spec: parentRefs: [{
			name:      "default"
			namespace: metadata.namespace
		}]
		spec: rules: [
			{
				matches: [{path: {type: "PathPrefix", value: "/"}}]
				backendRefs: [{
					name:      Name
					namespace: #HTTPBin.Namespace
					port:      #HTTPBin.Port
				}]
			},
		]
	}
}
