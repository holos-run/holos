package holos

import "encoding/json"

let ComponentName = "\(#InstancePrefix)-ingress"

#TargetNamespace: "istio-ingress"

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-stores":           _
		_dependsOn: "\(#InstancePrefix)-istio-base": _
		_dependsOn: "\(#InstancePrefix)-istiod":     _

		metadata: name: ComponentName

		chart: name: "gateway"
		namespace: #TargetNamespace
		_values: #GatewayValues & {
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
		apiObjectMap: OBJECTS.apiObjectMap
		// Auth Proxy
		apiObjectMap: _IngressAuthProxy.Deployment.apiObjectMap
		// Auth Policy
		apiObjectMap: _IngressAuthProxy.Policy.apiObjectMap
		// Auth Policy Exclusions
		apiObjectMap: _AuthPolicyRules.objects.apiObjectMap
	},
]

_ProxyProtocol: gatewayTopology: proxyProtocol: {}

// Additional holos specific API Objects
let Name = #GatewayValues.name
let GatewayLabels = {
	app:   Name
	istio: "ingressgateway"
}
let RedirectMetaName = {
	name:      Name + "-https-redirect"
	namespace: #TargetNamespace
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: "istio-ingress-mtls-cert": #ExternalSecret & {
			metadata: name:      "istio-ingress-mtls-cert"
			metadata: namespace: #TargetNamespace
		}
		Gateway: {
			"\(RedirectMetaName.name)": #Gateway & {
				metadata: RedirectMetaName
				spec: selector: GatewayLabels
				spec: servers: [{
					port: {
						number:   80
						name:     "http2"
						protocol: "HTTP2"
					}
					hosts: ["*"]
					// handled by the VirtualService
					tls: httpsRedirect: false
				}]
			}
		}
		VirtualService: {
			"\(RedirectMetaName.name)": #VirtualService & {
				metadata: RedirectMetaName
				spec: hosts: ["*"]
				spec: gateways: [RedirectMetaName.name]
				spec: http: [{
					match: [{withoutHeaders: ":path": prefix: "/.well-known/acme-challenge/"}]
					redirect: {
						scheme:       "https"
						redirectCode: 302
					}
				}]
			}
		}
		Deployment: {
			loopback: #Deployment & {
				_description: LoopbackDescription
				metadata:     LoopbackMetaName
				spec: {
					selector: matchLabels: LoopbackLabels
					template: {
						metadata: {
							annotations: "inject.istio.io/templates": "gateway"
							annotations: #Description & {
								_Description: LoopbackDescription
							}
							labels: LoopbackLabels & {"sidecar.istio.io/inject": "true"}
						}
						spec: {
							serviceAccountName: "istio-ingressgateway"
							// Allow binding to all ports (such as 80 and 443)
							securityContext: {
								runAsNonRoot: true
								seccompProfile: type: "RuntimeDefault"
								sysctls: [{name: "net.ipv4.ip_unprivileged_port_start", value: "0"}]
							}
							containers: [{
								name:  "istio-proxy"
								image: "auto" // Managed by istiod
								securityContext: {
									allowPrivilegeEscalation: false
									capabilities: drop: ["ALL"]
									runAsUser:  1337
									runAsGroup: 1337
								}
							}]
						}
					}
				}
			}
		}
		Service: {
			loopback: #Service & {
				_description: LoopbackDescription
				metadata:     LoopbackMetaName
				spec: selector:      LoopbackLabels
				spec: ports: [{port: 80, name: "http"}, {port: 443, name: "https"}]
			}
		}
	}
}

let LoopbackName = Name + "-loopback"
let LoopbackDescription = "Allows in-cluster traffic to stay in cluster via traffic routing"
let LoopbackLabels = {
	app:   LoopbackName
	istio: "ingressgateway"
}
let LoopbackMetaName = {
	name:      LoopbackName
	namespace: #TargetNamespace
}
