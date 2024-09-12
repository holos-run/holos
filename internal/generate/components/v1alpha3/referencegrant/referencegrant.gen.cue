package holos

import rg "gateway.networking.k8s.io/referencegrant/v1beta1"

#ReferenceGrant: rg.#ReferenceGrant & {
	metadata: name:      #Istio.Gateway.Namespace
	metadata: namespace: string
	spec: from: [{
		group:     "gateway.networking.k8s.io"
		kind:      "HTTPRoute"
		namespace: #Istio.Gateway.Namespace
	}]
	spec: to: [{
		group: ""
		kind:  "Service"
	}]
}
