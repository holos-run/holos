package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let BankName = #BankOfHolos.Name

let CommonLabels = {
	application: BankName
	environment: "development"
	team:        "ledger"
	tier:        "backend"
}

let Objects = {
	Name:      "bank-ledger-writer"
	Namespace: #BankOfHolos.Backend.Namespace

	// Ensure resources go in the correct namespace
	Resources: [_]: [_]: metadata: namespace: Namespace
	Resources: [_]: [_]: metadata: labels:    CommonLabels

	// https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/release/v0.6.5/kubernetes-manifests
	Resources: {
		Service: ledgerwriter: {
			apiVersion: "v1"
			kind:       "Service"
			metadata: {
				name:   "ledgerwriter"
				labels: CommonLabels
			}
			spec: {
				ports: [{
					name:       "http"
					port:       8080
					targetPort: 8080
				}]
				selector: {
					app: "ledgerwriter"
					CommonLabels
				}
				type: "ClusterIP"
			}
		}

		Deployment: ledgerwriter: {
			apiVersion: "apps/v1"
			kind:       "Deployment"
			metadata: {
				name:   "ledgerwriter"
				labels: CommonLabels
			}
			spec: {
				selector: matchLabels: {
					app: "ledgerwriter"
					CommonLabels
				}
				template: {
					metadata: {
						labels: {
							app: "ledgerwriter"
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
								name:  "ENABLE_METRICS"
								value: "false"
							}, {
								name:  "JVM_OPTS"
								value: "-XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -Xms256m -Xmx512m"
							}, {
								name:  "LOG_LEVEL"
								value: "info"
							}, {
								name: "NAMESPACE"
								valueFrom: fieldRef: fieldPath: "metadata.namespace"
							}]
							envFrom: [{
								configMapRef: name: "environment-config"
							}, {
								configMapRef: name: "service-api-config"
							}, {
								configMapRef: name: "ledger-db-config"
							}]
							image: "us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/ledgerwriter:v0.6.5@sha256:5b66d6888b87993c8ebe260fe33005c4e4bc2bdae4b5682874e1a078d37ff3b2"
							name:  "ledgerwriter"
							readinessProbe: {
								httpGet: {
									path: "/ready"
									port: 8080
								}
								initialDelaySeconds: 60
								periodSeconds:       5
								timeoutSeconds:      10
							}
							resources: {
								limits: {
									cpu:                 "500m"
									"ephemeral-storage": "0.5Gi"
									memory:              "512Mi"
								}
								requests: {
									cpu:                 "100m"
									"ephemeral-storage": "0.5Gi"
									memory:              "256Mi"
								}
							}
							securityContext: {
								allowPrivilegeEscalation: false
								capabilities: drop: ["all"]
								privileged:             false
								readOnlyRootFilesystem: true
							}
							startupProbe: {
								failureThreshold: 30
								httpGet: {
									path: "/ready"
									port: 8080
								}
								periodSeconds: 10
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
