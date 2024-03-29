package holos

import (
	"strings"
)

// Platform level definition of a project.
#Project: {
	name: string
	// All projects have at least a prod environment and stage.
	environments: prod: stage: "prod"
	environments: prod: dnsSegments: []
	stages: prod: _
	stages: dev:  _
	// Ensure at least the project name is a short hostname.  Additional may be added.
	hosts: (name): _
}

#ProjectTemplate: {
	project: #Project

	// GatewayServers maps Gateway spec.servers #GatewayServer values indexed by stage then name.
	let GatewayServers = {
		// Initialize all stages, even if they have no environments.
		for stage in project.stages {
			(stage.name): {}
		}

		for env in project.environments {
			(env.stage): {
				for host in project.hosts {
					let NAME = "https-\(env.name)-\(project.name)-\(host.name)"
					let SEGMENTS = [host.name] + env.dnsSegments + [#Platform.org.domain]
					let HOST = strings.Join(SEGMENTS, ".")
					(NAME): #GatewayServer & {
						hosts: ["\(env.namespace)/\(HOST)"]
						port: name:          NAME
						port: number:        443
						port: protocol:      "HTTPS"
						tls: credentialName: HOST
						tls: mode:           "SIMPLE"
					}

					for cluster in project.clusters {
						let NAME = "https-\(cluster.name)-\(project.name)-\(env.name)-\(host.name)"
						let SEGMENTS = [host.name] + env.dnsSegments + [cluster.name, #Platform.org.domain]
						let HOST = strings.Join(SEGMENTS, ".")
						(NAME): #GatewayServer & {
							hosts: ["\(env.namespace)/\(HOST)"]
							port: name:          NAME
							port: number:        443
							port: protocol:      "HTTPS"
							tls: credentialName: HOST
							tls: mode:           "SIMPLE"
						}
					}
				}
			}
		}
	}

	workload: resources: {
		// Provide resources only if the project is managed on --cluster-name
		if project.clusters[#ClusterName] != _|_ {
			for stage in project.stages {
				// Istio Gateway
				"\(stage.slug)-gateway": #KubernetesObjects & {
					apiObjectMap: (#APIObjects & {
						apiObjects: Gateway: (stage.slug): #Gateway & {
							spec: servers: [for server in GatewayServers[stage.name] {server}]
						}

						for host in GatewayServers[stage.name] {
							apiObjects: ExternalSecret: (host.tls.credentialName): metadata: namespace: "istio-ingress"
						}
					}).apiObjectMap
				}

				// Manage httpbin in each environment
				for Env in project.environments if Env.stage == stage.name {
					let Component = "\(Env.slug)-httpbin"
					(Component): #KubernetesObjects & {
						apiObjectMap: (#APIObjects & {
							let Project = project
							apiObjects: (HTTPBIN & {env: Env, project: Project}).apiObjects
						}).apiObjectMap
					}
				}
			}
		}
	}

	provisioner: resources: {
		for stage in project.stages {
			"\(stage.slug)-certs": #KubernetesObjects & {
				apiObjectMap: (#APIObjects & {
					for host in GatewayServers[stage.name] {
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

let HTTPBIN = {
	name:    string | *"httpbin"
	project: #Project
	env:     #Environment
	let Name = name

	let Metadata = {
		name:      Name
		namespace: env.namespace
		labels: app: name
	}

	apiObjects: {
		Deployment: (Name): #Deployment & {
			metadata: Metadata

			spec: selector: matchLabels: Metadata.labels
			spec: template: {
				metadata: labels: Metadata.labels & #IstioSidecar
				spec: securityContext: seccompProfile: type: "RuntimeDefault"
				spec: containers: [{
					name:  Name
					image: "quay.io/holos/mccutchen/go-httpbin"
					ports: [{containerPort: 8080}]
					securityContext: {
						seccompProfile: type: "RuntimeDefault"
						allowPrivilegeEscalation: false
						runAsNonRoot:             true
						runAsUser:                8192
						runAsGroup:               8192
						capabilities: drop: ["ALL"]
					}}]
			}
		}
		Service: (Name): #Service & {
			metadata: Metadata
			spec: selector: Metadata.labels
			spec: ports: [
				{port: 80, targetPort: 8080, protocol: "TCP", name: "http"},
			]
		}
		VirtualService: (Name): #VirtualService & {
			metadata: Metadata
			let Project = project
			let Env = env
			spec: hosts: [for host in (#EnvHosts & {project: Project, env: Env}).hosts {host.name}]
			spec: gateways: ["istio-ingress/\(env.stageSlug)"]
			spec: http: [{route: [{destination: host: Name}]}]
		}
	}
}
