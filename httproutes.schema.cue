package holos

import v1 "gateway.networking.k8s.io/httproute/v1"

// HTTPRoutes is where routes are registered.  The httproutes component manages
// routes by composing this struct into a BuildPlan.
HTTPRoutes: #HTTPRoutes

// #HTTPRoutes defines the schema of managed HTTPRoute resources for the
// platform.
#HTTPRoutes: {
	// For the guides, we simplify this down to a flat namespace.
	[Name=string]: v1.#HTTPRoute & {
		let HOST = Name + "." + Organization.Domain

		_backendRefs: [NAME=string]: {
			name:      string | *NAME
			namespace: string
			port:      number | *80
		}

		metadata: name:      Name
		metadata: namespace: Istio.Gateway.Namespace
		metadata: labels: app: Name
		spec: hostnames: [HOST]
		spec: parentRefs: [{
			name:      "default"
			namespace: metadata.namespace
		}]
		spec: rules: [
			{
				matches: [{path: {type: "PathPrefix", value: "/"}}]
				backendRefs: [for x in _backendRefs {x}]
			},
		]
	}
}
