package holos

#ProjectTemplate: {
	project:        _
	GatewayServers: _

	// Sort GatewayServers by the tls credentialName to issue wildcards
	let GatewayCerts = {
		for FQDN, Server in GatewayServers {
			"\(Server.tls.credentialName)": #Certificate & {
				let CertInfo = Server._CertInfo

				// Store the dnsNames in a struct so they can be collected into a list
				_dnsNames: "\(CertInfo.wildcard)": CertInfo.wildcard

				metadata: name:      CertInfo.canonical & Server.tls.credentialName
				metadata: namespace: "istio-ingress"
				spec: {
					commonName: CertInfo.canonical
					secretName: CertInfo.canonical & Server.tls.credentialName
					dnsNames: [for x in _dnsNames {x}]
					issuerRef: {
						kind: "ClusterIssuer"
						name: "letsencrypt-staging"
					}
				}
			}
		}
	}

	// Resources to be managed on the provisioner cluster.
	provisioner: resources: {
		for stage in project.stages {
			"\(stage.slug)-certs": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					apiObjects: Certificate: GatewayCerts
				}).apiObjectMap
			}
		}
	}
}
