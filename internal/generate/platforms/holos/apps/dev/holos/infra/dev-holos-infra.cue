package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

_AppInfo: spec: component: "infra"

let Objects = {
	Name:      _AppInfo.status.component
	Namespace: _AppInfo.metadata.namespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	let AWS_ACCOUNT = _Platform.Model.aws.accountNumber
	let Metadata = _AppInfo.metadata

	Resources: {
		ServiceAccount: holos: {
			metadata: Metadata
			// TODO(jeff): The ecr-creds-refresher name should be refactored to a root
			// level private var so we can update it in one place.
			// Refer to ecr-creds-refresher.cue
			imagePullSecrets: [{name: "ecr-creds-\(AWS_ACCOUNT)"}]
		}

		PostgresCluster: holos: {
			apiVersion: "postgres-operator.crunchydata.com/v1beta1"
			metadata: name: "holos"
			metadata: Metadata
			spec: {
				image: "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.1-0"
				instances: [{
					affinity: podAntiAffinity: preferredDuringSchedulingIgnoredDuringExecution: [{
						podAffinityTerm: {
							labelSelector: matchLabels: "postgres-operator.crunchydata.com/cluster": metadata.name
							topologyKey: "topology.kubernetes.io/zone"
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
	}
}
