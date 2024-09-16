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
	Name:      "bank-contacts"
	Namespace: #BankOfHolos.Backend.Namespace

	// Ensure resources go in the correct namespace
	Resources: [_]: [_]: metadata: namespace: Namespace
	Resources: [_]: [_]: metadata: labels:    CommonLabels

	// https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/release/v0.6.5/kubernetes-manifests
	Resources: {
		Service: contacts: {
			apiVersion: "v1"
			kind:       "Service"
			metadata: name: "contacts"
			spec: {
				ports: [{
					name:       "http"
					port:       8080
					targetPort: 8080
				}]
				selector: {
					app: "contacts"
					CommonLabels
				}
				type: "ClusterIP"
			}
		}

		Deployment: contacts: {
			apiVersion: "apps/v1"
			kind:       "Deployment"
			metadata: name: "contacts"
			spec: {
				selector: matchLabels: {
					app: "contacts"
					CommonLabels
				}
				template: {
					metadata: {
						labels: {
							app: "contacts"
							CommonLabels
						}
					}
					spec: {
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
							}]
							envFrom: [{
								configMapRef: name: "environment-config"
							}, {
								configMapRef: name: "accounts-db-config"
							}]
							image: "us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/contacts:v0.6.5@sha256:e451dcac7d34a7bde979c7f02d4c7ebd83a77aff373e1131ce3a2bba2f7fdc1a"
							name:  "contacts"
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
									cpu:                 "250m"
									"ephemeral-storage": "0.25Gi"
									memory:              "128Mi"
								}
								requests: {
									cpu:                 "100m"
									"ephemeral-storage": "0.25Gi"
									memory:              "64Mi"
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
								name:      "publickey"
								readOnly:  true
							}]
						}]
						securityContext: {
							fsGroup:      1000
							runAsGroup:   1000
							runAsNonRoot: true
							runAsUser:    1000
						}
						serviceAccountName:            BankName
						terminationGracePeriodSeconds: 5
						volumes: [{
							emptyDir: {}
							name: "tmp"
						}, {
							name: "publickey"
							secret: {
								items: [{
									key:  "jwtRS256.key.pub"
									path: "publickey"
								}]
								secretName: "jwt-key"
							}
						}]
					}
				}
			}
		}
	}
}
