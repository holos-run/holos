package holos

import "encoding/yaml"

#InstancePrefix: "prod-mesh"

#IstioVersion: "1.21.0"

// The ingress gateway auth proxy is used by multiple cue instances.
// AUTHPROXY configures one oauth2-proxy deployment for each host in each stage of a project.  Multiple deployments per stage are used to narrow down the cookie domain.
_IngressAuthProxy: {
	Name:          "authproxy"
	Namespace:     "istio-ingress"
	service:       "\(Name).\(Namespace).svc.cluster.local"
	AuthProxySpec: #AuthProxySpec & #Platform.authproxy

	Domains: [DOMAIN=string]: {name: DOMAIN}
	Domains: (#Platform.org.domain):                    _
	Domains: "\(#ClusterName).\(#Platform.org.domain)": _

	// TODO: This should be generated from ProjectHosts
	Domains: "holos.run":                 _
	Domains: "\(#ClusterName).holos.run": _

	let Metadata = {
		name:      string
		namespace: Namespace
		labels: "app.kubernetes.io/name":    name
		labels: "app.kubernetes.io/part-of": "istio-ingressgateway"
		...
	}

	let ProxyMetadata = Metadata & {name: Name}
	let RedisMetadata = Metadata & {name: Name + "-redis"}

	// Deployment represents the oauth2-proxy deployment
	Deployment: #APIObjects & {
		apiObjects: {
			// oauth2-proxy
			ExternalSecret: (Name): metadata: ProxyMetadata
			// Place the ID token in a header that does not conflict with the Authorization header.
			// Refer to: https://github.com/oauth2-proxy/oauth2-proxy/issues/1877#issuecomment-1364033723
			ConfigMap: (Name): {
				metadata: ProxyMetadata
				data: "config.yaml": yaml.Marshal(AuthProxyConfig)
				let AuthProxyConfig = {
					injectResponseHeaders: [{
						name: "x-oidc-id-token"
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
				metadata: ProxyMetadata

				spec: {
					replicas: 1
					selector: matchLabels: ProxyMetadata.labels
					template: {
						metadata: labels: ProxyMetadata.labels
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
									"--cookie-name=__Secure-\(#ClusterName)-ingress-\(Name)",
									"--cookie-samesite=lax",
									for domain in Domains {"--cookie-domain=.\(domain.name)"},
									for domain in Domains {"--cookie-domain=\(domain.name)"},
									for domain in Domains {"--whitelist-domain=.\(domain.name)"},
									for domain in Domains {"--whitelist-domain=\(domain.name)"},
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
									// echo '{"cookiesecret":"'$(LC_ALL=C tr -dc "[:alpha:]" </dev/random | tr '[:upper:]' '[:lower:]' | head -c 32)'"}' | holos create secret -n istio-ingress --append-hash=false --data-stdin authproxy
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
				metadata: ProxyMetadata
				spec: selector: ProxyMetadata.labels
				spec: ports: [
					{port: 4180, targetPort: 4180, protocol: "TCP", name: "http"},
				]
			}
			VirtualService: (Name): #VirtualService & {
				metadata: ProxyMetadata
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

	// Policy represents the AuthorizationPolicy and RequestAuthentication policy
	Policy: #APIObjects & {
		apiObjects: {
			RequestAuthentication: (Name): #RequestAuthentication & {
				metadata: Metadata & {name: Name}
				spec: jwtRules: [{
					audiences: ["\(AuthProxySpec.projectID)"]
					forwardOriginalToken: true
					fromHeaders: [{name: AuthProxySpec.idTokenHeader}]
					issuer: AuthProxySpec.issuer
				}]
				spec: selector: matchLabels: istio: "ingressgateway"
			}
			AuthorizationPolicy: "\(Name)-custom": {
				_description: "Route all requests through the auth proxy by default"

				metadata: Metadata & {name: "\(Name)-custom"}
				spec: {
					action: "CUSTOM"
					provider: name: AuthProxySpec.provider
					rules: [
						{
							to: [{
								// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
								operation: notHosts: [
									// Never send requests for the login service through the authorizer, would block login.
									AuthProxySpec.issuerHost,
									"\(AuthProxySpec.issuerHost):*",
									// Exclude hosts with specialized rules from the catch-all.
									for x in _AuthPolicyRules.hosts {x.name},
									for x in _AuthPolicyRules.hosts {"\(x.name):*"},
								]
							}]
							when: [
								{
									// bypass the external authorizer when the id token is already in the request.
									// the RequestAuthentication rule will verify the token.
									key: "request.headers[\(AuthProxySpec.idTokenHeader)]"
									notValues: ["*"]
								},
							]
						},
					]
					selector: matchLabels: istio: "ingressgateway"
				}
			}
		}
	}
}

_AuthPolicyRules: #AuthPolicyRules & {
	hosts: {
		let Vault = "vault.core.ois.run"
		"\(Vault)": {
			slug: "vault"
			// Rules for when to route requests through the auth proxy
			spec: rules: [
				{
					to: [{
						operation: hosts: [Vault]
						operation: paths: ["/ui", "/ui/*"]
					}]
				},
				{
					to: [{
						operation: hosts: [Vault]
					}]
					when: [{
						key: "request.headers[x-vault-request]"
						notValues: ["true"]
					}]
				},
			]
		}
	}
}

// Exclude project hosts from the auth proxy if configured to do so.  The
// intended effect is to exclude the host from the blanket `authproxy-custom`
// AuthorizationPolicy rule _without_ adding a specialized AuthorizationPolicy
// for the same host.  This has the effect of completely excluding the host from
// authorization policy.
for Project in _Projects {
	let ProjectHosts = (#ProjectHosts & {project: Project}).Hosts

	for FQDN, Host in ProjectHosts {
		if Host.NoAuthorizationPolicy {
			if Host.clusters[#ClusterName] != _|_ {
				_AuthPolicyRules: hosts: "\(Host.fqdn)": NoAuthorizationPolicy: true
			}
		}
	}
}
