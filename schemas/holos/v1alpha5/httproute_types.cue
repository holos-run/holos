package v1alpha5

import (
	v1 "gateway.networking.k8s.io/httproute/v1"
	rg "gateway.networking.k8s.io/referencegrant/v1beta1"
)

#HTTPRoutes: [string]: v1.#HTTPRoute

#HTTPRouteBuilder: {
	Name:             string
	Domain:           string
	GatewayNamespace: string | *"istio-ingress"

	BackendRefs: [NAME=string]: {
		name:      string | *NAME
		namespace: string
		port:      number | *80
	}

	HTTPRoute: v1.#HTTPRoute & {
		let HOST = Name + "." + Domain

		metadata: name:      Name
		metadata: namespace: GatewayNamespace
		metadata: labels: app: Name
		spec: hostnames: [HOST]
		spec: parentRefs: [{
			name:      "default"
			namespace: metadata.namespace
		}]
		spec: rules: [{
			matches: [{path: {type: "PathPrefix", value: "/"}}]
			backendRefs: [for x in BackendRefs {x}]
		}]
	}
}

// #ReferenceGrantBuilder builds a ReferenceGrant.  Useful from within a
// component definition to grant the HTTPRoute access to the namespace the
// component is managed in.
//
// Usage:
//  Component: Resources: #ReferenceGrantBuilder & {Namespace: NAMESPACE}
#ReferenceGrantBuilder: {
	Namespace:        string
	GatewayNamespace: string | *"istio-ingress"

	ReferenceGrant: (GatewayNamespace): rg.#ReferenceGrant & {
		metadata: name:      GatewayNamespace
		metadata: namespace: Namespace
		spec: from: [{
			group:     "gateway.networking.k8s.io"
			kind:      "HTTPRoute"
			namespace: "istio-ingress"
		}]
		spec: to: [{
			group: ""
			kind:  "Service"
		}]
	}
}
