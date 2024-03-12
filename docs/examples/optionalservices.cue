package holos

// Controls optional feature flags for services distributed across multiple holos components.
// For example, enable issuing certificates in the provisioner cluster when an optional service is
// enabled for a workload cluster.

#OptionalService: {
	name:    string
	enabled: true | *false
	clusters: [Name=_]: #Platform.clusters[Name]
	clusterNames: [for k, v in clusters {k}]
	namespaces: [Name=_]: #ManagedNamespace & {
		name: Name
	}
	// servers represents istio Gateway.spec.servers.hosts entries
	// Refer to istio/gateway/gateway.cue
	servers: [Name=_]: {
		hosts: [...string]
		port: name:          Name
		port: number:        443
		port: protocol:      "HTTPS"
		tls: credentialName: string
		tls: mode:           "SIMPLE"
	}
	// public tls certs should align to hosts.
	certs: [Name=_]: #Certificate & {
		metadata: name: Name
	}
}

#OptionalServices: {
	[Name=_]: #OptionalService & {
		name: Name
	}
}

for k, v in #OptionalServices {
	#ManagedNamespaces: v.namespaces
}
