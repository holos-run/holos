package holos

#InputKeys: component: "postgres"

#KubernetesObjects & {
	apiObjects: {
		PostgresCluster: db: #PostgresCluster & HighlyAvailable & {
			metadata: name:      _DBName
			metadata: namespace: #TargetNamespace
			spec: {
				image:           "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0"
				postgresVersion: 16
				users: [
					{name: "postgres"},
					{name: _DBName},
				]
				users: [...{databases: [_DBName]}]
				instances: [{
					replicas: 2
					dataVolumeClaimSpec: {
						accessModes: ["ReadWriteOnce"]
						resources: requests: storage: "1Gi"
					}
				}]
				backups: pgbackrest: {
					image: "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.49-0"
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
	}
}

// Refer to https://github.com/holos-run/postgres-operator-examples/blob/main/kustomize/high-availability/ha-postgres.yaml
let HighlyAvailable = {
	apiVersion: "postgres-operator.crunchydata.com/v1beta1"
	kind:       "PostgresCluster"
	metadata: name: string | *"hippo-ha"
	spec: {
		image:           "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0"
		postgresVersion: 16
		instances: [{
			name:     "pgha1"
			replicas: 2
			dataVolumeClaimSpec: {
				accessModes: ["ReadWriteOnce"]
				resources: requests: storage: "1Gi"
			}
			affinity: podAntiAffinity: preferredDuringSchedulingIgnoredDuringExecution: [{
				weight: 1
				podAffinityTerm: {
					topologyKey: "kubernetes.io/hostname"
					labelSelector: matchLabels: {
						"postgres-operator.crunchydata.com/cluster":      "hippo-ha"
						"postgres-operator.crunchydata.com/instance-set": "pgha1"
					}
				}
			}]
		}]
		backups: pgbackrest: {
			image: "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.49-0"
			repos: [{
				name: "repo1"
				volume: volumeClaimSpec: {
					accessModes: ["ReadWriteOnce"]
					resources: requests: storage: "1Gi"
				}
			}]
		}
		proxy: pgBouncer: {
			image:    "registry.developers.crunchydata.com/crunchydata/crunchy-pgbouncer:ubi8-1.21-3"
			replicas: 2
			affinity: podAntiAffinity: preferredDuringSchedulingIgnoredDuringExecution: [{
				weight: 1
				podAffinityTerm: {
					topologyKey: "kubernetes.io/hostname"
					labelSelector: matchLabels: {
						"postgres-operator.crunchydata.com/cluster": "hippo-ha"
						"postgres-operator.crunchydata.com/role":    "pgbouncer"
					}
				}
			}]
		}
	}
}
