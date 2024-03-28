package holos

import "strings"

// Platform level definition of a project.
#Project: {
	name: string
	// All projects have at least a prod environment and stage.
	environments: prod: stage: "prod"
	environments: prod: dnsSegments: []
	stages: prod: _
	stages: dev:  _
	// Short hostnames to construct fqdns.
	hosts: (name): _
}

#ProjectTemplate: {
	project: #Project

	// ExtAuthzHosts maps host names to the backend environment namespace for ExtAuthz.
	let ExtAuthzHosts = {
		for stage in project.stages {
			(stage.name): {}
		}

		for env in project.environments {
			(env.stage): {
				for host in project.hosts {
					let NAME = "https-\(project.name)-\(env.name)-\(host.name)"
					let SEGMENTS = [host.name] + env.dnsSegments + [#Platform.org.domain]
					let HOST = strings.Join(SEGMENTS, ".")
					(NAME): #GatewayServer & {
						hosts: ["\(env.namespace)/\(HOST)"]
						// name must be unique across all servers in all gateways
						port: name:     NAME
						port: number:   443
						port: protocol: "HTTPS"
						// TODO: Manage a certificate with each host in the dns alt names.
						tls: credentialName: HOST
						tls: mode:           "SIMPLE"
					}

					for cluster in project.clusters {
						let NAME = "https-\(cluster.name)-\(project.name)-\(env.name)-\(host.name)"
						let SEGMENTS = [host.name] + env.dnsSegments + [cluster.name, #Platform.org.domain]
						let HOST = strings.Join(SEGMENTS, ".")
						(NAME): #GatewayServer & {
							hosts: ["\(env.namespace)/\(HOST)"]
							// name must be unique across all servers in all gateways
							port: name:     NAME
							port: number:   443
							port: protocol: "HTTPS"
							// TODO: Manage a certificate with each host in the dns alt names.
							tls: credentialName: HOST
							tls: mode:           "SIMPLE"
						}
					}
				}
			}
		}
	}

	workload: resources: {
		for stage in project.stages {
			// System namespace for each project stage
			let SystemName = "\(stage.slug)-system"
			(SystemName): #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					apiObjects: Namespace: (SystemName):   _
					apiObjects: SecretStore: (SystemName): _
				}).apiObjectMap
			}

			// Project namespace for each project environment
			"\(stage.slug)-namespaces": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					for env in project.environments if env.stage == stage.name {
						apiObjects: Namespace: (env.slug):   _
						apiObjects: SecretStore: (env.slug): _
					}
				}).apiObjectMap
			}

			// Istio Gateway
			"\(stage.slug)-gateway": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					apiObjects: Gateway: (stage.slug): #Gateway & {
						spec: servers: [for host in ExtAuthzHosts[stage.name] {host}]
					}

					for host in ExtAuthzHosts[stage.name] {
						apiObjects: ExternalSecret: (host.tls.credentialName): metadata: namespace: "istio-ingress"
					}
				}).apiObjectMap
			}
		}
	}

	provisioner: resources: {
		for stage in project.stages {
			"\(stage.slug)-namespaces": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					apiObjects: Namespace: "\(stage.slug)-system": _
					for env in project.environments if env.stage == stage.name {
						apiObjects: Namespace: (env.slug): _
					}
				}).apiObjectMap
			}

			"\(stage.slug)-certs": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					for host in ExtAuthzHosts[stage.name] {
						let CN = host.tls.credentialName
						apiObjects: Certificate: (CN): #Certificate & {
							metadata: name:      CN
							metadata: namespace: "istio-ingress"
							spec: {
								commonName: CN
								dnsNames: [CN]
								secretName: CN
								issuerRef: {
									kind: "ClusterIssuer"
									name: "letsencrypt"
								}
							}
						}
					}

				}).apiObjectMap
			}
		}
	}
}

// #GatewayServer defines the value of the istio Gateway.spec.servers field.
#GatewayServer: {
	// The ip or the Unix domain socket to which the listener should
	// be bound to.
	bind?:            string
	defaultEndpoint?: string

	// One or more hosts exposed by this gateway.
	hosts: [...string]

	// An optional name of the server, when set must be unique across
	// all servers.
	name?: string

	// The Port on which the proxy should listen for incoming
	// connections.
	port: {
		// Label assigned to the port.
		name: string

		// A valid non-negative integer port number.
		number: int

		// The protocol exposed on the port.
		protocol:    string
		targetPort?: int
	}

	// Set of TLS related options that govern the server's behavior.
	tls?: {
		// REQUIRED if mode is `MUTUAL` or `OPTIONAL_MUTUAL`.
		caCertificates?: string

		// Optional: If specified, only support the specified cipher list.
		cipherSuites?: [...string]

		// For gateways running on Kubernetes, the name of the secret that
		// holds the TLS certs including the CA certificates.
		credentialName?: string

		// If set to true, the load balancer will send a 301 redirect for
		// all http connections, asking the clients to use HTTPS.
		httpsRedirect?: bool

		// Optional: Maximum TLS protocol version.
		maxProtocolVersion?: "TLS_AUTO" | "TLSV1_0" | "TLSV1_1" | "TLSV1_2" | "TLSV1_3"

		// Optional: Minimum TLS protocol version.
		minProtocolVersion?: "TLS_AUTO" | "TLSV1_0" | "TLSV1_1" | "TLSV1_2" | "TLSV1_3"

		// Optional: Indicates whether connections to this port should be
		// secured using TLS.
		mode?: "PASSTHROUGH" | "SIMPLE" | "MUTUAL" | "AUTO_PASSTHROUGH" | "ISTIO_MUTUAL" | "OPTIONAL_MUTUAL"

		// REQUIRED if mode is `SIMPLE` or `MUTUAL`.
		privateKey?: string

		// REQUIRED if mode is `SIMPLE` or `MUTUAL`.
		serverCertificate?: string

		// A list of alternate names to verify the subject identity in the
		// certificate presented by the client.
		subjectAltNames?: [...string]

		// An optional list of hex-encoded SHA-256 hashes of the
		// authorized client certificates.
		verifyCertificateHash?: [...string]

		// An optional list of base64-encoded SHA-256 hashes of the SPKIs
		// of authorized client certificates.
		verifyCertificateSpki?: [...string]
	}
}
