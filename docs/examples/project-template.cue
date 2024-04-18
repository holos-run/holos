package holos

import (
	h "github.com/holos-run/holos/api/v1alpha1"
	"encoding/yaml"
)

// let SourceLoc = "project-template.cue"

#ProjectTemplate: {
	project: #Project

	// workload cluster resources
	workload: resources: [Name=_]: h.#KubernetesObjects & {
		metadata: name: Name
	}

	// provisioner cluster resources
	provisioner: resources: [Name=_]: h.#KubernetesObjects & {
		metadata: name: Name
	}
}

// Reference Platform Project Template
#ProjectTemplate: {
	project: #Project
	let Project = project

	ProjectHosts: (#ProjectHosts & {project: Project}).Hosts

	// GatewayServers maps Gateway spec.servers #GatewayServer values indexed by stage then name.
	GatewayServers: {
		for FQDN, Host in ProjectHosts {
			// If the host is valid on the cluster being rendered
			if Host.clusters[#ClusterName] != _|_ {
				"\(FQDN)": #GatewayServer & {
					_CertInfo: Host
					hosts: [
						"\(Host.env.namespace)/\(FQDN)",
						// Allow the authproxy VirtualService to match the project.authProxyPrefix path.
						"\(Host.stage.namespace)/\(FQDN)",
					]
					port: {
						name:     "https"
						number:   443
						protocol: "HTTPS"
					}
					tls: credentialName: Host.canonical
					tls: mode:           "SIMPLE"
				}
			}
		}
	}

	// ClusterDefaultGatewayServers provides a struct of Gateway servers for the current cluster.
	// This is intended for Gateway/default to add all servers to the default gateway.
	ClusterDefaultGatewayServers: {
		if project.clusters[#ClusterName] != _|_ {
			GatewayServers
		}
	}

	workload: resources: {
		// Provide resources only if the project is managed on the cluster specified
		// by --cluster-name
		if project.clusters[#ClusterName] != _|_ {
			for stage in project.stages {
				let Stage = stage

				// Manage auth-proxy in each stage
				if project.features.authproxy.enabled {
					"\(stage.slug)-authproxy": #KubernetesObjects & {
						apiObjectMap: (#APIObjects & {
							apiObjects: (AUTHPROXY & {stage: Stage, project: Project, servers: GatewayServers[stage.name]}).apiObjects
						}).apiObjectMap
					}

					for Env in project.environments if Env.stage == stage.name {
						"\(Env.slug)-authpolicy": #KubernetesObjects & {
							// Manage auth policy in each env
							apiObjectMap: (#APIObjects & {
								apiObjects: (AUTHPOLICY & {env: Env, project: Project, servers: GatewayServers[stage.name]}).apiObjects
							}).apiObjectMap
						}
					}
				}

				// Manage httpbin in each environment
				if project.features.httpbin.enabled {
					for Env in project.environments if Env.stage == stage.name {
						"\(Env.slug)-httpbin": #KubernetesObjects & {
							let Project = project
							apiObjectMap: (#APIObjects & {
								apiObjects: (HTTPBIN & {env: Env, project: Project}).apiObjects
							}).apiObjectMap
						}
					}
				}
			}
		}
	}
}

// AUTHPROXY configures one oauth2-proxy deployment for each host in each stage of a project.  Multiple deployments per stage are used to narrow down the cookie domain.
let AUTHPROXY = {
	name:    string | *"authproxy"
	project: #Project
	stage:   #Stage
	servers: {}
	let Name = name
	let Project = project
	let Stage = stage

	let AuthProxySpec = #AuthProxySpec & {
		namespace: stage.namespace
		projectID: project.resourceId
		clientID:  stage.authProxyClientID
		orgDomain: project.authProxyOrgDomain
		provider:  stage.extAuthzProviderName
	}

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
		// Place the ID token in a header that does not conflict with the Authorization header.
		// Refer to: https://github.com/oauth2-proxy/oauth2-proxy/issues/1877#issuecomment-1364033723
		ConfigMap: (Name): {
			metadata: Metadata
			data: "config.yaml": yaml.Marshal(AuthProxyConfig)
			let AuthProxyConfig = {
				injectResponseHeaders: [{
					name: AuthProxySpec.idTokenHeader
					values: [{claim: "id_token"}]
				}]
				providers: [{
					id:                    "Holos Platform"
					name:                  "Holos Platform"
					provider:              "oidc"
					scope:                 "openid profile email groups offline_access urn:zitadel:iam:org:domain:primary:\(AuthProxySpec.orgDomain)"
					clientID:              AuthProxySpec.clientID
					clientSecretFile:      "/dev/null"
					code_challenge_method: "S256"
					loginURLParameters: [{
						default: ["force"]
						name: "approval_prompt"
					}]
					oidcConfig: {
						issuerURL: AuthProxySpec.issuer
						audienceClaims: ["aud"]
						emailClaim:  "email"
						groupsClaim: "groups"
						userIDClaim: "sub"
					}
				}]
				server: BindAddress: ":4180"
				upstreamConfig: upstreams: [{
					id:         "static://200"
					path:       "/"
					static:     true
					staticCode: 200
				}]
			}
		}
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
					spec: {
						securityContext: seccompProfile: type: "RuntimeDefault"
						containers: [{
							image:           "quay.io/oauth2-proxy/oauth2-proxy:v7.6.0"
							imagePullPolicy: "IfNotPresent"
							name:            "oauth2-proxy"
							volumeMounts: [{
								name:      "config"
								mountPath: "/config"
								readOnly:  true
							}]
							args: [
								// callback url is proxy prefix + /callback
								"--proxy-prefix=" + AuthProxySpec.proxyPrefix,
								"--email-domain=*",
								"--session-store-type=redis",
								"--redis-connection-url=redis://\(RedisMetadata.name):6379",
								"--cookie-refresh=12h",
								"--cookie-expire=2160h",
								"--cookie-secure=true",
								"--cookie-name=__Secure-\(stage.slug)-\(Name)",
								"--cookie-samesite=lax",
								for domain in StageDomains {"--cookie-domain=.\(domain.name)"},
								for domain in StageDomains {"--cookie-domain=\(domain.name)"},
								for domain in StageDomains {"--whitelist-domain=.\(domain.name)"},
								for domain in StageDomains {"--whitelist-domain=\(domain.name)"},
								"--cookie-csrf-per-request=true",
								"--cookie-csrf-expire=120s",
								// will skip authentication for OPTIONS requests
								"--skip-auth-preflight=true",
								"--real-client-ip-header=X-Forwarded-For",
								"--skip-provider-button=true",
								"--auth-logging",
								"--alpha-config=/config/config.yaml",
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
						volumes: [{name: "config", configMap: name: Name}]
					}
				}
			}
		}
		Service: (Name): #Service & {
			metadata: Metadata
			spec: selector: Metadata.labels
			spec: ports: [
				{port: 4180, targetPort: 4180, protocol: "TCP", name: "http"},
			]
		}
		VirtualService: (Name): #VirtualService & {
			metadata: Metadata
			spec: hosts: ["*"]
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{
				match: [{uri: prefix: AuthProxySpec.proxyPrefix}]
				route: [{
					destination: host: Name
					destination: port: number: 4180
				}]
			}]
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
								capabilities: drop: ["ALL"]
								runAsNonRoot: true
								runAsUser:    999
								runAsGroup:   999
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

let HTTPBIN = {
	name:    string | *"httpbin"
	project: #Project
	env:     #Environment
	let Name = name
	let Stage = project.stages[env.stage]

	let Metadata = {
		name:      Name
		namespace: env.namespace
		labels: app: name
	}
	let Labels = {
		"app.kubernetes.io/name":       Name
		"app.kubernetes.io/instance":   env.slug
		"app.kubernetes.io/part-of":    env.project
		"security.holos.run/authproxy": Stage.extAuthzProviderName
	}

	apiObjects: {
		Deployment: (Name): #Deployment & {
			metadata: Metadata

			spec: selector: matchLabels: Metadata.labels
			spec: template: {
				metadata: labels: Metadata.labels & #IstioSidecar & Labels
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
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{route: [{destination: host: Name}]}]
		}
	}
}

// AUTHPOLICY configures the baseline AuthorizationPolicy and RequestAuthentication policy for each stage of each project.
let AUTHPOLICY = {
	project: #Project
	env:     #Environment
	let Name = "\(stage.slug)-authproxy"
	let Project = project
	let stage = project.stages[env.stage]
	let Env = env

	let AuthProxySpec = #AuthProxySpec & {
		namespace: stage.namespace
		projectID: project.resourceId
		clientID:  stage.authProxyClientID
		orgDomain: project.authProxyOrgDomain
		provider:  stage.extAuthzProviderName
	}

	let Metadata = {
		name:      string
		namespace: env.namespace
		labels: {
			"app.kubernetes.io/name":     name
			"app.kubernetes.io/instance": stage.name
			"app.kubernetes.io/part-of":  stage.project
		}
	}

	// Collect all the hosts associated with the stage
	let Hosts = {
		for HOST in (#EnvHosts & {project: Project, env: Env}).hosts {
			(HOST.name): HOST
		}
	}

	// HostList is a list of hosts for AuthorizationPolicy rules
	let HostList = [
		for host in Hosts {host.name},
		for host in Hosts {host.name + ":*"},
	]
	let MatchLabels = {"security.holos.run/authproxy": AuthProxySpec.provider}

	apiObjects: {
		RequestAuthentication: (Name): #RequestAuthentication & {
			metadata: Metadata & {name: Name}
			spec: jwtRules: [{
				audiences: [AuthProxySpec.clientID]
				forwardOriginalToken: true
				fromHeaders: [{name: AuthProxySpec.idTokenHeader}]
				issuer: AuthProxySpec.issuer
			}]
			spec: selector: matchLabels: MatchLabels
		}
		AuthorizationPolicy: "\(Name)-custom": {
			metadata: Metadata & {name: "\(Name)-custom"}
			spec: {
				action: "CUSTOM"
				// send the request to the auth proxy
				provider: name: AuthProxySpec.provider
				rules: [{
					to: [{operation: hosts: HostList}]
					when: [
						{
							key: "request.headers[\(AuthProxySpec.idTokenHeader)]"
							notValues: ["*"]
						},
						{
							key: "request.headers[host]"
							notValues: [AuthProxySpec.issuerHost]
						},
					]}]
				selector: matchLabels: MatchLabels
			}
		}
	}
}
