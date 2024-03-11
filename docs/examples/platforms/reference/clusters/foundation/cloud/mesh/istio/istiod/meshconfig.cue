package holos

// Istio meshconfig
// TODO: Generate per-project extauthz providers.
_MeshConfig: {
	accessLogEncoding: "JSON"
	accessLogFile:     "/dev/stdout"
	defaultConfig: {
		discoveryAddress: "istiod.istio-system.svc:15012"
		tracing: zipkin: address: "zipkin.istio-system:9411"
	}
	defaultProviders: metrics: ["prometheus"]
	enablePrometheusMerge: true
	// For PROXY PROTOCOL at the ingress gateway.
	gatewayTopology: {
		numTrustedProxies: 2
	}
	rootNamespace: "istio-system"
	trustDomain:   "cluster.local"
	extensionProviders: [{
		name: "cluster-trace"
		zipkin: {
			maxTagLength: 56
			port:         9411
			service:      "zipkin.istio-system.svc"
		}
	}, {
		name: "cluster-gatekeeper"
		envoyExtAuthzHttp: {
			headersToDownstreamOnDeny: [
				"content-type",
				"set-cookie",
			]
			headersToUpstreamOnAllow: [
				"authorization",
				"path",
				"x-auth-request-user",
				"x-auth-request-email",
				"x-auth-request-access-token",
			]
			includeAdditionalHeadersInCheck: "X-Auth-Request-Redirect": "%REQ(x-forwarded-proto)%://%REQ(:authority)%%REQ(:path)%%REQ(:query)%"
			includeRequestHeadersInCheck: [
				"authorization",
				"cookie",
				"x-forwarded-for",
			]
			port:    4180
			service: "oauth2-proxy.istio-ingress.svc.cluster.local"
		}
	}, {
		name: "core-authorizer"
		envoyExtAuthzHttp: {
			headersToDownstreamOnDeny: [
				"content-type",
				"set-cookie",
			]
			headersToUpstreamOnAllow: [
				"authorization",
				"path",
				"x-auth-request-user",
				"x-auth-request-email",
				"x-auth-request-access-token",
			]
			includeAdditionalHeadersInCheck: "X-Auth-Request-Redirect": "%REQ(x-forwarded-proto)%://%REQ(:authority)%%REQ(:path)%%REQ(:query)%"
			includeRequestHeadersInCheck: [
				"authorization",
				"cookie",
				"x-forwarded-for",
			]
			port:    4180
			service: "oauth2-proxy.prod-core-system.svc.cluster.local"
		}
	}]
}
