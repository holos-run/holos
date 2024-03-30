package holos

// Platform level definition of a project.
#Project: {
	name: string

	// All projects have at least a prod environment and stage.
	stages: prod: stageSegments: []
	environments: prod: stage: "prod"
	environments: prod: envSegments: []
	stages: dev: _
	environments: dev: stage: "dev"
	environments: dev: envSegments: []
	// Ensure at least the project name is a short hostname.  Additional may be added.
	hosts: (name): _

	// environments share the stage segments of their stage.
	environments: [_]: {
		stage:         string
		stageSegments: stages[stage].stageSegments
	}
}

#ProjectTemplate: {
	project: #Project
	let Project = project

	// GatewayServers maps Gateway spec.servers #GatewayServer values indexed by stage then name.
	let GatewayServers = {
		// Initialize all stages, even if they have no environments.
		for stage in project.stages {
			(stage.name): {
				let Stage = stage

				// entry for the auth-proxy domain
				for host in (#StageDomains & {project: Project, stage: Stage, prefixes: [["auth"]]}).hosts {
					(host.name): #GatewayServer & {
						hosts: ["\(stage.namespace)/\(host.name)"]
						port: host.port
						tls: credentialName: host.name
						tls: mode:           "SIMPLE"
					}
				}
			}
		}

		// For each stage, construct entries for the Gateway spec.servers.hosts field.
		for env in project.environments {
			(env.stage): {
				let Env = env
				for host in (#EnvHosts & {project: Project, env: Env}).hosts {
					(host.name): #GatewayServer & {
						hosts: ["\(env.namespace)/\(host.name)"]
						port: host.port
						tls: credentialName: host.name
						tls: mode:           "SIMPLE"
					}
				}
			}
		}
	}

	workload: resources: {
		// Provide resources only if the project is managed on --cluster-name
		if project.clusters[#ClusterName] != _|_ {
			for stage in project.stages {
				let Stage = stage

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

				// Manage auth-proxy in each stage
				"\(stage.slug)-authproxy": #KubernetesObjects & {
					apiObjectMap: (#APIObjects & {
						apiObjects: (AUTHPROXY & {stage: Stage, project: Project}).apiObjects
					}).apiObjectMap
				}

				// Manage httpbin in each environment
				for Env in project.environments if Env.stage == stage.name {
					"\(Env.slug)-httpbin": #KubernetesObjects & {
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

// AUTHPROXY configures one oauth2-proxy deployment for each host in each stage of a project.  Multiple deployments per stage are used to narrow down the cookie domain.
let AUTHPROXY = {
	name:    string | *"authproxy"
	project: #Project
	stage:   #Stage
	let Name = name
	let Project = project
	let Stage = stage

	let Metadata = {
		name:      Name
		namespace: stage.namespace
		labels: {
			"app.kubernetes.io/name":     name
			"app.kubernetes.io/instance": stage.name
			"app.kubernetes.io/part-of":  stage.project
		}
	}

	let RedisMetadata = {
		name:      Name + "-redis"
		namespace: stage.namespace
		labels: {
			"app.kubernetes.io/name":     name
			"app.kubernetes.io/instance": stage.name
			"app.kubernetes.io/part-of":  stage.project
		}
	}

	apiObjects: {
		// oauth2-proxy
		ExternalSecret: (Name): metadata: Metadata
		Deployment: (Name): #Deployment & {
			metadata: Metadata

			// project.dev.example.com, project.dev.k1.example.com, project.dev.k2.example.com
			let StageDomains = {
				for host in (#StageDomains & {project: Project, stage: Stage}).hosts {
					(host.name): host
				}
			}

			spec: {
				replicas: 1
				selector: matchLabels: Metadata.labels
				template: {
					metadata: labels: Metadata.labels
					metadata: labels: #IstioSidecar
					spec: securityContext: seccompProfile: type: "RuntimeDefault"
					spec: containers: [{
						image:           "quay.io/oauth2-proxy/oauth2-proxy:v7.4.0"
						imagePullPolicy: "IfNotPresent"
						name:            "oauth2-proxy"
						args: [
							"--proxy-prefix=/oauth2",
							"--email-domain=*",
							"--session-store-type=redis",
							"--redis-connection-url=redis://\(RedisMetadata.name):6379",
							"--cookie-refresh=12h",
							"--cookie-expire=2160h",
							"--cookie-secure=true",
							"--cookie-name=__Secure-\(Name)-\(stage.slug)",
							"--cookie-samesite=lax",
							for domain in StageDomains {"--cookie-domain=.\(domain.name)"},
							for domain in StageDomains {"--whitelist-domain=.\(domain.name)"},
							"--cookie-csrf-per-request=true",
							"--cookie-csrf-expire=120s",
							"--set-authorization-header=false",
							"--set-xauthrequest=true",
							"--pass-access-token=true",
							"--pass-authorization-header=true",
							"--upstream=static://200",
							"--reverse-proxy",
							"--real-client-ip-header=X-Forwarded-For",
							"--skip-provider-button=true",
							"--auth-logging",
							"--provider=oidc",
							"--client-id=" + stage.authProxyClientID,
							"--client-secret-file=/dev/null",
							"--oidc-issuer-url=https://login.\(#Platform.org.domain)",
							"--code-challenge-method=S256",
							"--http-address=0.0.0.0:4180",
							// "--allowed-group=\(project.resourceId):\(stage.name)-access",
						]
						env: [{
							name: "OAUTH2_PROXY_COOKIE_SECRET"
							// echo '{"cookiesecret":"'$(LC_ALL=C tr -dc "[:alpha:]" </dev/random | tr '[:upper:]' '[:lower:]' | head -c 32)'"}' | holos create secret -n dev-holos-system --append-hash=false --data-stdin authproxy
							valueFrom: secretKeyRef: {
								key:  "cookiesecret"
								name: Name
							}
						}]
						ports: [{
							containerPort: 4180
							protocol:      "TCP"
						}]
						securityContext: {
							seccompProfile: type: "RuntimeDefault"
							allowPrivilegeEscalation: false
							runAsNonRoot:             true
							runAsUser:                8192
							runAsGroup:               8192
							capabilities: drop: ["ALL"]
						}
					}]
				}
			}
		}
		Service: (Name): #Service & {
			metadata: Metadata
			spec: selector: Metadata.labels
			spec: ports: [
				{port: 80, targetPort: 4180, protocol: "TCP", name: "http"},
			]
		}

		// redis
		ConfigMap: (RedisMetadata.name): {
			metadata: RedisMetadata
			data: "redis.conf": """
				maxmemory 128mb
				maxmemory-policy allkeys-lru
				"""
		}
		Deployment: (RedisMetadata.name): {
			metadata: RedisMetadata
			spec: {
				selector: matchLabels: RedisMetadata.labels
				template: {
					metadata: labels: RedisMetadata.labels
					metadata: labels: #IstioSidecar
					spec: securityContext: seccompProfile: type: "RuntimeDefault"
					spec: {
						containers: [{
							command: [
								"redis-server",
								"/redis-master/redis.conf",
							]
							env: [{
								name:  "MASTER"
								value: "true"
							}]
							image: "quay.io/holos/redis:7.2.4"
							livenessProbe: {
								initialDelaySeconds: 15
								tcpSocket: port: "redis"
							}
							name: "redis"
							ports: [{
								containerPort: 6379
								name:          "redis"
							}]
							readinessProbe: {
								exec: command: [
									"redis-cli",
									"ping",
								]
								initialDelaySeconds: 5
							}
							resources: limits: cpu: "0.5"
							securityContext: {
								seccompProfile: type: "RuntimeDefault"
								allowPrivilegeEscalation: false
								runAsNonRoot:             true
								runAsUser:                8192
								runAsGroup:               8192
								capabilities: drop: ["ALL"]
							}
							volumeMounts: [{
								mountPath: "/redis-master-data"
								name:      "data"
							}, {
								mountPath: "/redis-master"
								name:      "config"
							}]
						}]
						volumes: [{
							emptyDir: {}
							name: "data"
						}, {
							configMap: name: RedisMetadata.name
							name: "config"
						}]
					}
				}
			}
		}
		Service: (RedisMetadata.name): #Service & {
			metadata: RedisMetadata
			spec: selector: RedisMetadata.labels
			spec: type:     "ClusterIP"
			spec: ports: [{
				name:       "redis"
				port:       6379
				protocol:   "TCP"
				targetPort: 6379
			}]
		}
	}
}
