package holos

#InputKeys: component: "postgres"

let S3Secret = "pgo-s3-creds"
let ZitadelUser = _DBName
let ZitadelAdmin = "\(_DBName)-admin"

#KubernetesObjects & {
	apiObjects: {
		ExternalSecret: "pgo-s3-creds": _
		PostgresCluster: db: #PostgresCluster & HighlyAvailable & {
			metadata: name:      _DBName
			metadata: namespace: #TargetNamespace
			spec: {
				image:           "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0"
				postgresVersion: 16
				// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/references/crd/5.5.x/postgrescluster#postgresclusterspecusersindex
				users: [
					{name: ZitadelUser},
					// NOTE: Users with SUPERUSER role cannot log in through pgbouncer.  Use options that allow zitadel admin to use pgbouncer.
					// Refer to: https://github.com/CrunchyData/postgres-operator/issues/3095#issuecomment-1904712211
					{name: ZitadelAdmin, options: "CREATEDB CREATEROLE", databases: [_DBName, "postgres"]},
				]
				users: [...{databases: [_DBName, ...]}]
				instances: [{
					replicas: 2
					dataVolumeClaimSpec: {
						accessModes: ["ReadWriteOnce"]
						resources: requests: storage: string | *"1Gi"
					}
				}]
				// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/backups
				backups: pgbackrest: {
					configuration: [{secret: name: S3Secret}]
					manual: {
						// Note the repoName value must match the config keys in the S3Secret.
						repoName: "repo2"
						options: ["--type=full", ...]
					}
					global: {
						// Store only one full backup in the PV because it's more expensive than object storage.
						"\(repos[0].name)-retention-full": "1"
						// Store 14 days of full backups in the bucket.
						"\(manual.repoName)-retention-full":      string | *"14"
						"\(manual.repoName)-retention-full-type": "count" | *"time" // time in days
						// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/backups#encryption
						"\(manual.repoName)-cipher-type": "aes-256-cbc"
						// "The convention we recommend for setting this variable is /pgbackrest/$NAMESPACE/$CLUSTER_NAME/repoN"
						// Ref: https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/backups#understanding-backup-configuration-and-basic-operations
						"\(manual.repoName)-path": "/pgbackrest/\(#TargetNamespace)/\(metadata.name)/\(manual.repoName)"
					}
					repos: [
						{
							name: "repo1"
							volume: volumeClaimSpec: {
								accessModes: ["ReadWriteOnce"]
								resources: requests: storage: string | *"1Gi"
							}
						},
						{
							name: manual.repoName
							// Full backup weekly on Sunday at 1am, differntial daily at 1am every day except Sunday.
							schedules: full:         string | *"0 1 * * 0"
							schedules: differential: string | *"0 1 * * 1-6"
							s3: {
								bucket:   string | *"\(#Platform.org.name)-zitadel-backups"
								region:   string | *#Backups.s3.region
								endpoint: string | *"s3.dualstack.\(region).amazonaws.com"
							}
						},
					]
				}
			}
		}
	}
}

// Refer to https://github.com/holos-run/postgres-operator-examples/blob/main/kustomize/high-availability/ha-postgres.yaml
let HighlyAvailable = {
	apiVersion: "postgres-operator.crunchydata.com/v1beta1"
	kind:       "PostgresCluster"
	metadata: name: string
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
						"postgres-operator.crunchydata.com/cluster":      metadata.name
						"postgres-operator.crunchydata.com/instance-set": name
					}
				}
			}]
		}]
		backups: pgbackrest: {
			image: "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.49-0"
		}
		proxy: pgBouncer: {
			image:    "registry.developers.crunchydata.com/crunchydata/crunchy-pgbouncer:ubi8-1.21-3"
			replicas: 2
			affinity: podAntiAffinity: preferredDuringSchedulingIgnoredDuringExecution: [{
				weight: 1
				podAffinityTerm: {
					topologyKey: "kubernetes.io/hostname"
					labelSelector: matchLabels: {
						"postgres-operator.crunchydata.com/cluster": metadata.name
						"postgres-operator.crunchydata.com/role":    "pgbouncer"
					}
				}
			}]
		}
	}
}
