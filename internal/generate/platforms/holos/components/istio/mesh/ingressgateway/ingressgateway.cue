package holos

import "encoding/json"

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "gateway"
	Version:   #IstioVersion
	Namespace: "istio-ingress"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #GatewayValues & {
		// This component expects the load balancer to send the PROXY protocol header.
		// Refer to: https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.2/guide/service/annotations/#proxy-protocol-v2
		podAnnotations: "proxy.istio.io/config": json.Marshal(_ProxyProtocol)
		service: {
			type: "LoadBalancer"
			annotations: "service.beta.kubernetes.io/aws-load-balancer-proxy-protocol": "*"
			annotations: "service.beta.kubernetes.io/aws-load-balancer-alpn-policy":    "HTTP2Preferred"
			annotations: "service.beta.kubernetes.io/aws-load-balancer-scheme":         "internet-facing"
			externalTrafficPolicy: "Local"
			// Add 30000 to the port to get the Nodeport
			ports: [
				{
					name:       "status-port"
					port:       15021
					protocol:   "TCP"
					targetPort: 15021
					nodePort:   30021
				},
				{
					name:       "http2"
					port:       80
					protocol:   "TCP"
					targetPort: 80
					nodePort:   30080
				},
				{
					name:       "https"
					port:       443
					protocol:   "TCP"
					targetPort: 443
					nodePort:   30443
				},
			]
		}
	}
}

_ProxyProtocol: gatewayTopology: proxyProtocol: {}
