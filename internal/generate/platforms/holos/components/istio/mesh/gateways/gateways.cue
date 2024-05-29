package holos

import "encoding/json"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "gateway"
	Namespace: #IstioGatewaysNamespace

	Resources: {
		// The *.admin.clustername.example.com cert
		for listener in Gateway.default.spec.listeners {
			for ref in listener.tls.certificateRefs {
				if ref.kind == "Secret" {
					ExternalSecret: admin: #ExternalCert & {
						metadata: name: ref.name
					}
				}
			}
		}

		// The default gateway with all listeners attached to tls certs.
		Gateway: default: {
			metadata: {
				namespace: Namespace
				annotations: "service.beta.kubernetes.io/aws-load-balancer-proxy-protocol": "*"
				annotations: "service.beta.kubernetes.io/aws-load-balancer-alpn-policy":    "HTTP2Preferred"
				annotations: "service.beta.kubernetes.io/aws-load-balancer-scheme":         "internet-facing"
				annotations: "proxy.istio.io/config":                                       json.Marshal(_ProxyProtocol)
			}
			spec: {
				// Work with a struct of listeners instead of a list.
				_listeners: (#WildcardListener & {Name: "admin"}).Output
				listeners: [for x in _listeners {x}]
			}
		}
	}
}

#WildcardListener: {
	Name: string

	Output: [NAME=string]: {name: NAME}
	Output: {
		"\(Name)-apex": {
			hostname: "\(Name).\(_ClusterName).\(_Platform.Model.org.domain)"
			port:     443
			protocol: "HTTPS"
			tls: {
				certificateRefs: [{
					kind: "Secret"
					name: "\(Name).\(_ClusterName).\(_Platform.Model.org.domain)"
				}]
			}
			allowedRoutes: namespaces: from:     "Selector"
			allowedRoutes: namespaces: selector: _AdminSelector
		}
		"\(Name)-prefix": {
			hostname: "*.\(Name).\(_ClusterName).\(_Platform.Model.org.domain)"
			port:     443
			protocol: "HTTPS"
			tls: {
				certificateRefs: [{
					kind: "Secret"
					name: "\(Name).\(_ClusterName).\(_Platform.Model.org.domain)"
				}]
			}
			allowedRoutes: namespaces: from:     "Selector"
			allowedRoutes: namespaces: selector: _AdminSelector
		}
	}
}

_ProxyProtocol: gatewayTopology: proxyProtocol: {}
