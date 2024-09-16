package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let BankName = #BankOfHolos.Name

let CommonLabels = {
	application: BankName
	environment: "development"
	team:        "accounts"
	tier:        "backend"
}

let Objects = {
	Name:      "bank-userservice"
	Namespace: #BankOfHolos.Backend.Namespace

	// Ensure resources go in the correct namespace
	Resources: [_]: [_]: metadata: namespace: Namespace

	// https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/release/v0.6.5/kubernetes-manifests/userservice.yaml
	Resources: {
		Service: userservice: {
			metadata: name:   "userservice"
			metadata: labels: CommonLabels
			spec: {
				selector: {
					app: "userservice"
					CommonLabels
				}
				_ports: http: {
					name:       "http"
					port:       8080
					targetPort: 8080
					protocol:   "TCP"
				}
				ports: [for x in _ports {x}]
			}
		}

		Deployment: userservice: {
			metadata: name:   "userservice"
			metadata: labels: CommonLabels
			spec: {
				selector: matchLabels: {
					app: "userservice"
					CommonLabels
				}
				template: {
					metadata: labels: {
						app: "userservice"
						CommonLabels
					}
					spec: {
						serviceAccountName:            BankName
						terminationGracePeriodSeconds: 5
						containers: [{
							env: [{
								name:  "VERSION"
								value: "v0.6.5"
							}, {
								name:  "PORT"
								value: "8080"
							}, {
								name:  "ENABLE_TRACING"
								value: "false"
							}, {
								name:  "LOG_LEVEL"
								value: "info"
							}, {
								name:  "TOKEN_EXPIRY_SECONDS"
								value: "3600"
							}, {
								name:  "PRIV_KEY_PATH"
								value: "/tmp/.ssh/privatekey"
							}]
							envFrom: [{
								configMapRef: name: "environment-config"
							}, {
								configMapRef: name: "accounts-db-config"
							}]
							image: "us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/userservice:v0.6.5@sha256:f91e0e5bd6cdb16f6b867b2e3e874b23dd01f11592de006776f1dfb136702941"
							name:  "userservice"
							ports: [{
								containerPort: 8080
								name:          "http-server"
							}]
							readinessProbe: {
								httpGet: {
									path: "/ready"
									port: 8080
								}
								initialDelaySeconds: 10
								periodSeconds:       5
								timeoutSeconds:      10
							}
							resources: {
								limits: {
									cpu:                 "500m"
									"ephemeral-storage": "0.25Gi"
									memory:              "256Mi"
								}
								requests: {
									cpu:                 "260m"
									"ephemeral-storage": "0.25Gi"
									memory:              "128Mi"
								}
							}
							securityContext: {
								allowPrivilegeEscalation: false
								capabilities: drop: ["all"]
								privileged:             false
								readOnlyRootFilesystem: true
							}
							volumeMounts: [{
								mountPath: "/tmp"
								name:      "tmp"
							}, {
								mountPath: "/tmp/.ssh"
								name:      "keys"
								readOnly:  true
							}]
						}]
						volumes: [{
							emptyDir: {}
							name: "tmp"
						}, {
							name: "keys"
							secret: {
								secretName: "jwt-key"
								items: [
									{
										key:  "jwtRS256.key"
										path: "privatekey"
									},
									{
										key:  "jwtRS256.key.pub"
										path: "publickey"
									},
								]
							}
						}]
						securityContext: {
							seccompProfile: type: "RuntimeDefault"
							fsGroup:      1000
							runAsGroup:   1000
							runAsNonRoot: true
							runAsUser:    1000
						}
					}
				}
			}
		}

		// Allow HTTPRoutes in the ingress gateway namespace to reference Services
		// in this namespace.
		ReferenceGrant: grant: #ReferenceGrant & {
			metadata: namespace: Namespace
		}

		// Include shared resources
		#BankOfHolos.Resources
	}
}
