package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "gateway"
	Namespace: #IstioGatewaysNamespace

	Resources: {
		// Manage a service account to prevent ArgoCD from pruning it.
		ServiceAccount: "default-istio": {
			metadata: namespace: Namespace
			metadata: labels: {
				"gateway.istio.io/managed":               "istio.io-gateway-controller"
				"gateway.networking.k8s.io/gateway-name": "default"
				"istio.io/gateway-name":                  "default"
			}
		}

		// The default gateway with all listeners attached to tls certs.
		Gateway: default: {
			metadata: namespace: Namespace
			spec: {
				// Work with a struct of listeners instead of a list.
				_listeners: (#WildcardListener & {Name: "admin", Selector: _Selector.GrantSubdomainAdmin, Cluster: true}).Output
				_listeners: (#WildcardListener & {Name: "login", Selector: _Selector.GrantSubdomainLogin, Cluster: false}).Output
				_listeners: (#WildcardListener & {Name: "app", Selector: _Selector.GrantSubdomainApp, Cluster: false}).Output
				_listeners: (#WildcardListener & {Name: "app", Selector: _Selector.GrantSubdomainApp, Cluster: true}).Output
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
	_Prefix:   string
	if Cluster == true {
		_Hostname: "\(Name).\(_ClusterName).\(_Platform.Model.org.domain)"
		_Prefix:   "region-\(Name)"
	}
	if Cluster == false {
		_Hostname: "\(Name).\(_Platform.Model.org.domain)"
		_Prefix:   "global-\(Name)"
	}

	Output: [NAME=string]: {name: NAME}
	Output: {
		"\(_Prefix)-apex": {
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
		"\(_Prefix)-prefix": {
			hostname: "*.\(_Hostname)"
			port:     443
			protocol: "HTTPS"
			tls: {
				certificateRefs: [{
					kind: "Secret"
					name: "any.\(_Hostname)"
				}]
			}
			allowedRoutes: namespaces: from:     "Selector"
			allowedRoutes: namespaces: selector: Selector
		}
	}
}

_ProxyProtocol: gatewayTopology: proxyProtocol: {}
