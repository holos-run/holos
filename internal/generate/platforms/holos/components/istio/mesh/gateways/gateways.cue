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
					ExternalSecret: "\(ref.name)": #ExternalCert & {
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
				_listeners: (#WildcardListener & {Name: "admin", Selector: _Selector.Admin}).Output
				_listeners: (#WildcardListener & {Name: "login", Selector: _Selector.Login, Cluster: false}).Output
				listeners: [for x in _listeners {x}]
			}
		}
	}
}

#WildcardListener: {
	Name:    string
	Cluster: false | *true
	Selector: matchLabels: {[string]: string}

	_Hostname: string
	if Cluster == true {
		_Hostname: "\(Name).\(_ClusterName).\(_Platform.Model.org.domain)"
	}
	if Cluster == false {
		_Hostname: "\(Name).\(_Platform.Model.org.domain)"
	}

	Output: [NAME=string]: {name: NAME}
	Output: {
		"\(Name)-apex": {
			hostname: _Hostname
			port:     443
			protocol: "HTTPS"
			tls: {
				certificateRefs: [{
					kind: "Secret"
					name: _Hostname
				}]
			}
			allowedRoutes: namespaces: from:     "Selector"
			allowedRoutes: namespaces: selector: Selector
		}
		"\(Name)-prefix": {
			hostname: "*.\(_Hostname)"
			port:     443
			protocol: "HTTPS"
			tls: {
				certificateRefs: [{
					kind: "Secret"
					name: _Hostname
				}]
			}
			allowedRoutes: namespaces: from:     "Selector"
			allowedRoutes: namespaces: selector: Selector
		}
	}
}

_ProxyProtocol: gatewayTopology: proxyProtocol: {}
