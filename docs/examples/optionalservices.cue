// Controls optional feature flags for services distributed across multiple holos components.
// For example, enable issuing certificates in the provisioner cluster when an optional service is
// enabled for a workload cluster.
package holos

import "list"

#OptionalService: {
	name:    string
	enabled: true | *false
	clusters: [Name=_]: #Platform.clusters[Name]
	clusterNames: [for c in clusters {c.name}]

	managedNamespaces: [Name=_]: #ManagedNamespace & {
		namespace: metadata: name: Name
		clusterNames: ["provisioner", for c in clusters {c.name}]
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

for svc in #OptionalServices {
	for nsName, ns in svc.managedNamespaces {
		if svc.enabled && list.Contains(ns.clusterNames, #ClusterName) {
			#ManagedNamespaces: "\(nsName)": ns
		}
	}
}
