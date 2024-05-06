package holos

let Namespace = "dev-holos"
let Holos = "holos"

// spec represents the output provided to holos
spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "dev-holos-infra"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let Metadata = {
	name:      Holos
	namespace: Namespace
	labels: app: Holos
}

// OBJECTS represents the kubernetes api objects to manage.
let OBJECTS = #APIObjects & {
	// Postgres
	// Deployment
	// VirtualService

	apiObjects: ServiceAccount: holos: {
		metadata: Metadata
		imagePullSecrets: [{name: "kube-system-ecr-image-pull-creds"}]
	}

	apiObjects: PostgresCluster: holos: {
		apiVersion: "postgres-operator.crunchydata.com/v1beta1"
		metadata:   Metadata
		spec: {
			image: "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.1-0"
			instances: [{
				affinity: podAntiAffinity: preferredDuringSchedulingIgnoredDuringExecution: [{
					podAffinityTerm: {
						labelSelector: matchLabels: "postgres-operator.crunchydata.com/cluster": "holos"
						topologyKey: "kubernetes.io/hostname"
					}
					weight: 1
				}]
				dataVolumeClaimSpec: {
					accessModes: ["ReadWriteOnce"]
					resources: requests: storage: "1Gi"
				}
				name:     "db"
				replicas: 1
			}]
			port:            5432
			postgresVersion: 16
			users: [{
				databases: ["holos"]
				name:    "holos"
				options: "SUPERUSER"
			}]
			backups: pgbackrest: {
				global: {
					"archive-async":          "y"
					"archive-push-queue-max": "100MiB"
					"spool-path":             "/pgdata/backups"
				}
				image: "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.47-2"
				repos: [{
					name: "repo1"
					volume: volumeClaimSpec: {
						accessModes: ["ReadWriteOnce"]
						resources: requests: storage: "1Gi"
					}
				}]
			}
		}
	}

	apiObjects: Service: holos: {
		apiVersion: "v1"
		metadata:   Metadata
		spec: {
			type: "ClusterIP"
			selector: app: "holos"
			ports: [{
				appProtocol: "http2"
				name:        "http"
				port:        3000
				protocol:    "TCP"
				targetPort:  3000
			}, {
				appProtocol: "http"
				name:        "metrics"
				port:        9090
				protocol:    "TCP"
				targetPort:  9090
			}]
		}
	}

	apiObjects: VirtualService: holos: {
		apiVersion: "networking.istio.io/v1beta1"
		metadata:   Metadata
		spec: {
			gateways: ["istio-ingress/default"]
			hosts: [
				"app.dev.holos.run",
				"app.dev.\(#ClusterName).holos.run",
			]
			http: [{
				match: [{
					uri: prefix: "/ui"
				}]
				name: "ui"
				route: [{
					destination: {
						host: "holos"
						port: number: 3000
					}
				}]
			}, {
				name: "api"
				route: [{
					destination: {
						host: "holos"
						port: number: 3000
					}
				}]
			}]
		}
	}
}
