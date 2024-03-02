package holos

import "encoding/json"

#InputKeys: component: "ingress"
#TargetNamespace: "istio-ingress"

#HelmChart & {
	chart: name: "gateway"
	namespace: #TargetNamespace
	values: #GatewayValues & {
		// This component expects the load balancer to send the PROXY protocol header.
		// Refer to: https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.2/guide/service/annotations/#proxy-protocol-v2
		podAnnotations: "proxy.istio.io/config": json.Marshal(_ProxyProtocol)
		// TODO This configuration is specific to the OIS Metal NLB, refactor it out to the metal collection.
		service: {
			type: "NodePort"
			annotations: "service.beta.kubernetes.io/aws-load-balancer-proxy-protocol": "*"
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
	apiObjects: _APIObjects
}

_ProxyProtocol: gatewayTopology: proxyProtocol: {}

// Additional holos specific API Objects
let LoopbackName = #GatewayValues.name + "-loopback"
let LoopbackDescription = "Allows in-cluster traffic to stay in cluster via traffic routing"
let LoopbackLabels = {
	app:   LoopbackName
	istio: "ingressgateway"
}

_APIObjects: {
	Deployment: {
		loopback: #Deployment & {
			_description: LoopbackDescription
			metadata: {
				name:      LoopbackName
				namespace: #TargetNamespace
			}
			spec: {
				selector: matchLabels: LoopbackLabels
				template: {
					metadata: {
						annotations: #CommonAnnotations & {
							_Description:                LoopbackDescription
							"inject.istio.io/templates": "gateway"
						}
						labels: LoopbackLabels & {"sidecar.istio.io/inject": "true"}
					}
					spec: {
						serviceAccountName: "istio-ingressgateway"
						// Allow binding to all ports (such as 80 and 443)
						securityContext: sysctls: [{name: "net.ipv4.ip_unprivileged_port_start", value: "0"}]
						containers: [{
							name:  "istio-proxy"
							image: "auto" // Managed by istiod
							securityContext: capabilities: drop: ["ALL"]
							securityContext: runAsUser:  1337
							securityContext: runAsGroup: 1337
						}]
					}
				}
			}
		}
	}
}
