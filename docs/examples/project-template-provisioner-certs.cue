package holos

#ProjectTemplate: {
	project:        _
	GatewayServers: _

	// Sort GatewayServers by the tls credentialName to issue wildcards
	let GatewayCerts = {
		for FQDN, Server in GatewayServers {
			let CertInfo = Server._CertInfo

			// Sort into stage for the holos components, e.g. prod-iam-certs, dev-iam-certs
			"\(CertInfo.stage.slug)": {
				"\(Server.tls.credentialName)": #Certificate & {

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
	}

	// Resources to be managed on the provisioner cluster.
	provisioner: resources: {
		for stage in project.stages {
			"\(stage.slug)-certs": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					apiObjects: Certificate: GatewayCerts[stage.slug]
				}).apiObjectMap
			}
		}
	}
}
