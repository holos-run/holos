package holos

import "strings"

// #AppInfo represents the data structure for an application deployed onto the
// platform.  This definition constraints the sechema defined at the root.
_AppInfo: #AppInfo & {
	metadata: name:      string
	metadata: namespace: "\(spec.env)-\(metadata.name)"

	spec: env: string
	spec: region: hostname: spec.dns.segments._region
	spec: global: hostname: spec.dns.segments._global

	spec: dns: segments: {
		env: [] | *[spec.env]
		name: [] | [string | *metadata.name]
		cluster: [] | *[_ClusterName]
		domain: [] | *[_Platform.Model.org.domain]
		_region: strings.Join(env+name+cluster+domain, ".")
		_global: strings.Join(env+name+domain, ".")
	}
}

// #AppRoute represents the HTTPRoute resources in the namespace of the Gateway.
#AppRoute: {
	AppInfo: #AppInfo

	Resources: {
		HTTPRoute: (HTTPRouteApp & {Hostname: AppInfo.spec.region.hostname}).HTTPRoute
		HTTPRoute: (HTTPRouteApp & {Hostname: AppInfo.spec.global.hostname}).HTTPRoute
	}

	let HTTPRouteApp = {
		Hostname: string

		HTTPRoute: (Hostname): {
			metadata: namespace:   #IstioGatewaysNamespace
			metadata: labels:      AppInfo.metadata.labels
			metadata: annotations: AppInfo.metadata.annotations
			spec: hostnames: [Hostname]
			spec: parentRefs: [{
				name:      "default"
				namespace: #IstioGatewaysNamespace
			}]
			spec: rules: [
				{
					matches: [{path: {type: "PathPrefix", value: "/"}}]
					backendRefs: [{
						name:      AppInfo.metadata.name
						namespace: AppInfo.metadata.namespace
						port:      AppInfo.spec.port
					}]
				},
				{
					// match the authproxy path prefix
					matches: [{path: {type: "PathPrefix", value: _AuthProxy.pathPrefix}}]
					backendRefs: [{
						name:      _AuthProxy.metadata.name
						namespace: _AuthProxy.metadata.namespace
						port:      _AuthProxy.servicePort
					}]
				},
			]
		}
	}
}
