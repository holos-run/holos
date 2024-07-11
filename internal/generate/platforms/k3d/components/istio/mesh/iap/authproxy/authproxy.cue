package holos

import "encoding/yaml"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let AuthProxyPrefix = _AuthProxy.pathPrefix

// Auth Proxy
// apiObjectMap: _IngressAuthProxy.Deployment.apiObjectMap
// Auth Policy
// apiObjectMap: _IngressAuthProxy.Policy.apiObjectMap
// Auth Policy Exclusions
//apiObjectMap: _AuthPolicyRules.objects.apiObjectMap

let Objects = {
	Name:      _AuthProxy.metadata.name
	Namespace: _AuthProxy.metadata.namespace

	let Metadata = _IAP.metadata

	let ProxyMetadata = Metadata & {name: Name}
	let RedisMetadata = Metadata & {name: Name + "-redis"}

	Resources: {
		// echo '{"cookiesecret":"'$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)'"}' \
		//   | holos create secret -n istio-gateways --append-hash=false --data-stdin authproxy
		// ExternalSecret: (Name): #ExternalSecret & {metadata: ProxyMetadata}
		// Place the ID token in a header that does not conflict with the Authorization header.
		// Refer to: https://github.com/oauth2-proxy/oauth2-proxy/issues/1877#issuecomment-1364033723
		ConfigMap: (Name): {
			metadata: ProxyMetadata
			data: "config.yaml": yaml.Marshal(AuthProxyConfig)
			let AuthProxyConfig = {
				injectResponseHeaders: [
					{
						name: _AuthProxy.idTokenHeader
						values: [{claim: "id_token"}]
					},
					{
						name: "x-auth-request-email"
						values: [{claim: "email"}]
					},
					{
						name: "x-auth-request-groups"
						values: [{claim: "groups"}]
					},
					{
						name: "x-forwarded-email"
						values: [{claim: "email"}]
					},
					{
						name: "x-forwarded-user"
						values: [{claim: "email"}]
					},
					{
						name: "x-forwarded-preferred-username"
						values: [{claim: "preferred_username"}]
					},
				]
				providers: [{
					id:       "Holos Platform"
					name:     "Holos Platform"
					provider: "oidc"
					// A specific org may be required with urn:zitadel:iam:org:domain:primary:\(_Platform.Model.org.domain)
					scope:                 "openid profile email groups offline_access"
					clientID:              _AuthProxy.clientID
					clientSecretFile:      "/dev/null"
					code_challenge_method: "S256"
					loginURLParameters: [{
						default: ["force"]
						name: "approval_prompt"
					}]
					oidcConfig: {
						issuerURL: _AuthProxy.issuerURL
						audienceClaims: ["aud"]
						emailClaim:  "email"
						groupsClaim: "groups"
						userIDClaim: "email"
					}
				}]
				server: BindAddress: ":\(_AuthProxy.servicePort)"
				upstreamConfig: upstreams: [{
					id:         "static://200"
					path:       "/"
					static:     true
					staticCode: 200
				}]
			}
		}

		Deployment: (Name): {
			metadata: ProxyMetadata

			spec: {
				replicas: 1
				selector: matchLabels: ProxyMetadata.labels
				template: {
					metadata: labels: ProxyMetadata.labels
					metadata: labels: "sidecar.istio.io/inject": "true"
					spec: {
						securityContext: seccompProfile: type: "RuntimeDefault"
						containers: [{
							// image: "quay.io/oauth3-proxy/oauth2-proxy:v7.6.0"
							image:           "quay.io/holos/oauth2-proxy:v7.6.0-1-g77a03ae2"
							imagePullPolicy: "IfNotPresent"
							name:            "oauth2-proxy"
							volumeMounts: [{
								name:      "config"
								mountPath: "/config"
								readOnly:  true
							}]
							args: [
								// callback url is Host + proxyPrefix + /callback
								"--proxy-prefix=" + AuthProxyPrefix,
								"--email-domain=*",
								"--session-store-type=redis",
								"--redis-connection-url=redis://authproxy-redis:6379",
								"--cookie-refresh=12h",
								"--cookie-expire=2160h",
								"--cookie-secure=true",
								"--cookie-name=__Secure-\(_ClusterName)-authproxy",
								"--cookie-samesite=lax",
								"--cookie-domain=.\(_ClusterName).\(_Platform.Model.org.domain)",
								"--cookie-domain=\(_ClusterName).\(_Platform.Model.org.domain)",
								"--cookie-domain=.\(_Platform.Model.org.domain)",
								"--cookie-domain=\(_Platform.Model.org.domain)",
								"--whitelist-domain=.\(_ClusterName).\(_Platform.Model.org.domain)",
								"--whitelist-domain=\(_ClusterName).\(_Platform.Model.org.domain)",
								"--whitelist-domain=.\(_Platform.Model.org.domain)",
								"--whitelist-domain=\(_Platform.Model.org.domain)",
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
								valueFrom: secretKeyRef: {
									key:  "cookiesecret"
									name: Name
								}
							}]
							ports: [{
								containerPort: _AuthProxy.servicePort
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

		Service: (Name): {
			metadata: ProxyMetadata
			spec: selector: ProxyMetadata.labels
			spec: ports: [
				{
					port:       _AuthProxy.servicePort
					targetPort: port
					protocol:   "TCP"
					name:       "http"
				},
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
					metadata: labels: "sidecar.istio.io/inject": "true"
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

		Service: (RedisMetadata.name): {
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
