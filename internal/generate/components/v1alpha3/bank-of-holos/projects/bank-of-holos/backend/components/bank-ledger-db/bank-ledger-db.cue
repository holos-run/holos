package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let BankName = #BankOfHolos.Name

let CommonLabels = {
	application: BankName
	environment: "development"
	team:        "ledger"
	tier:        "db"
}

let Objects = {
	Name:      "bank-ledger-db"
	Namespace: #BankOfHolos.Backend.Namespace

	// Ensure resources go in the correct namespace
	Resources: [_]: [_]: metadata: namespace: Namespace

	// https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/release/v0.6.5/kubernetes-manifests
	Resources: {
		ConfigMap: "ledger-db-config": {
			apiVersion: "v1"
			metadata: {
				name: "ledger-db-config"
				labels: {
					app: "ledger-db"
					CommonLabels
				}
			}
			data: {
				POSTGRES_DB:                "postgresdb"
				POSTGRES_PASSWORD:          "password"
				POSTGRES_USER:              "admin"
				SPRING_DATASOURCE_PASSWORD: "password"
				SPRING_DATASOURCE_URL:      "jdbc:postgresql://ledger-db:5432/postgresdb"
				SPRING_DATASOURCE_USERNAME: "admin"
			}
		}

		Service: "ledger-db": {
			apiVersion: "v1"
			kind:       "Service"
			metadata: {
				name:   "ledger-db"
				labels: CommonLabels
			}
			spec: {
				ports: [{
					name:       "tcp"
					port:       5432
					targetPort: 5432
				}]
				selector: {
					app: "ledger-db"
					CommonLabels
				}
				type: "ClusterIP"
			}
		}

		StatefulSet: "ledger-db": {
			apiVersion: "apps/v1"
			kind:       "StatefulSet"
			metadata: {
				name:   "ledger-db"
				labels: CommonLabels
			}
			spec: {
				replicas: 1
				selector: matchLabels: {
					app: "ledger-db"
					CommonLabels
				}
				serviceName: "ledger-db"
				template: {
					metadata: labels: {
						app: "ledger-db"
						CommonLabels
					}
					spec: {
						containers: [{
							envFrom: [{
								configMapRef: name: "environment-config"
							}, {
								configMapRef: name: "ledger-db-config"
							}, {
								configMapRef: name: "demo-data-config"
							}]
							image: "us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/ledger-db:v0.6.5@sha256:cc4fd25f301ab6d46b1312244d6931babc4c6cb66c5cb6d31d4a1adfa318a321"
							name:  "postgres"
							ports: [{containerPort: 5432}]
							resources: {
								limits: {
									cpu:    "250m"
									memory: "1Gi"
								}
								requests: {
									cpu:    "100m"
									memory: "512Mi"
								}
							}
							volumeMounts: [{
								mountPath: "/var/lib/postgresql/data"
								name:      "postgresdb"
								subPath:   "postgres"
							}]
						}]
						serviceAccountName: BankName
						volumes: [{
							emptyDir: {}
							name: "postgresdb"
						}]
					}
				}
			}
		}
	}
}
