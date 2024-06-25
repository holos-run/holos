package holos

// #MeshConfig provides the istio mesh config in the Output field.
#MeshConfig: {
	Output: {
		accessLogEncoding: string | *"JSON"
		accessLogFile:     string | *"/dev/stdout"
		defaultConfig: {
			discoveryAddress: string | *"istiod.istio-system.svc:15012"
			tracing: zipkin: address: string | *"zipkin.istio-system:9411"
		}
		defaultProviders: metrics: [...string] | *["prometheus"]
		enablePrometheusMerge: false | *true
		rootNamespace:         string | *"istio-system"
		trustDomain:           string | *"cluster.local"
		// NOTE: Use a Pod annotation to configure proxyProtocol, the proxyProtocol
		// field here has no effect when tested on istio 1.22.0. Refer to the
		// default Gateway in the gateways component for where the proxyProtocol is
		// configured.
		gatewayTopology: numTrustedProxies: 2
		extensionProviders: [
			// The auth proxy attached to Gateway default
			#ExtAuthzProxy & {
				name: _AuthProxy.provider
				envoyExtAuthzHttp: service: "\(_AuthProxy.metadata.name).\(_AuthProxy.metadata.namespace).svc.cluster.local"
			},
		]
	}
}

// #ExtAuthzProxy defines the provider configuration for an istio external authorization auth proxy.
#ExtAuthzProxy: {
	name: string
	envoyExtAuthzHttp: {
		headersToDownstreamOnDeny: [
			"content-type",
			"set-cookie",
		]
		headersToUpstreamOnAllow: [
			"authorization",
			"path",
			_AuthProxy.idTokenHeader,
			"x-forwaded-access-token",
			// For Backstage oauth2-proxy auth provider
			"x-forwarded-email",
			"x-forwarded-user",
			"x-forwarded-preferred-username",
			"x-auth-request-email",
			"x-auth-request-groups",
		]
		includeAdditionalHeadersInCheck: "X-Auth-Request-Redirect": "%REQ(x-forwarded-proto)%://%REQ(:authority)%%REQ(:path)%%REQ(:query)%"
		includeRequestHeadersInCheck: [
			"authorization",
			"cookie",
			"x-forwarded-for",
		]
		port:    _AuthProxy.servicePort
		service: string
	}
}
