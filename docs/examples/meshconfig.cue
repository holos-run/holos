package holos

// #MeshConfig provides the istio meshconfig in the config key given projects.
#MeshConfig: {
	projects: #Projects
	// clusterName is the value of the --cluster-name flag, the cluster currently being manged / rendered.
	clusterName: string | *#ClusterName

	extensionProviderMap: [Name=_]: {
		name: Name
	}

	config: {
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
		extensionProviders: [for x in extensionProviderMap {x}]
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
		service: string
	}
}
