package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

// Restore from backup.  Flip this to true after the database is provisioned and
// a backup has been taken.
let RestoreFromBackup = false

// The Secret containing the pgbackrest s3.conf file.
let S3Secret = "pgbackrest"

let Cluster = _Clusters[_ClusterName]

let DatabaseUser = _DBName
let DatabaseAdmin = "\(_DBName)-admin"

// This must be an external storage bucket for our architecture.
let BucketRepoName = "repo2"

// Restore options.  Set the timestamp to a known good point in time.
// time="2024-03-11T17:08:58Z" level=info msg="crunchy-pgbackrest ends"
// let RestoreOptions = ["--type=time", "--target=\"2024-03-11 17:10:00+00\""]

// Restore the most recent backup.
let RestoreOptions = []

let Objects = {
	Name:      "backstage-database"
	Namespace: "backstage"

	Resources: {
		// All resources go into the same namespace
		[_]: [_]: metadata: namespace: Namespace

		PostgresCluster: db: HighlyAvailable & {
			metadata: name: _DBName
			spec: {
				image:           "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0"
				postgresVersion: 16
				// Custom certs are necessary for streaming standby replication which we use to replicate between two regions.
				// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/disaster-recovery#streaming-standby
				customTLSSecret: name:            "\(_DBName)-primary-tls"
				customReplicationTLSSecret: name: "\(_DBName)-repl-tls"
				// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/references/crd/5.5.x/postgrescluster#postgresclusterspecusersindex
				users: [
					{name: DatabaseUser},
					// NOTE: Users with SUPERUSER role cannot log in through pgbouncer.  Use options that allow zitadel admin to use pgbouncer.
					// Refer to: https://github.com/CrunchyData/postgres-operator/issues/3095#issuecomment-1904712211
					{name: DatabaseAdmin, options: "CREATEDB CREATEROLE", databases: [_DBName, "postgres"]},
				]
				users: [...{databases: [_DBName, ...]}]
				instances: [{
					replicas: 2
					dataVolumeClaimSpec: {
						accessModes: ["ReadWriteOnce"]
						resources: requests: storage: "20Gi"
					}
				}]
				standby: {
					repoName: BucketRepoName
					if Cluster.primary {
						enabled: false
					}
					if !Cluster.primary {
						enabled: true
					}
				}
				// Monitoring configuration
				monitoring: pgmonitor: exporter: image: "registry.developers.crunchydata.com/crunchydata/crunchy-postgres-exporter:ubi8-5.5.1-0"
				// Restore from backup if and only if the cluster is primary and
				// RestoreFromBackup has transitioned from false to true.
				if Cluster.primary && RestoreFromBackup {
					dataSource: pgbackrest: {
						stanza:        "db"
						configuration: backups.pgbackrest.configuration
						// Restore from known good full backup taken
						options: RestoreOptions
						global: {
							"\(BucketRepoName)-path":        "/pgbackrest/\(metadata.namespace)/\(metadata.name)/\(BucketRepoName)"
							"\(BucketRepoName)-cipher-type": "aes-256-cbc"
						}
						repo: {
							name: BucketRepoName
							s3:   backups.pgbackrest.repos[1].s3
						}
					}
				}

				// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/backups
				backups: pgbackrest: {
					configuration: [{secret: name: S3Secret}]
					// Defines details for manual pgBackRest backup Jobs
					manual: {
						// Note: the repoName value must match the config keys in the S3Secret.
						// This must be an external repository for backup / restore / regional failovers.
						repoName: BucketRepoName
						options: ["--type=full", ...]
					}
					// Defines details for performing an in-place restore using pgBackRest
					restore: {
						// Enables triggering a restore by annotating the postgrescluster with postgres-operator.crunchydata.com/pgbackrest-restore="$(date)"
						enabled:  true
						repoName: BucketRepoName
					}
					global: {
						// Store only one full backup in the PV because it's more expensive than object storage.
						"\(repos[0].name)-retention-full": "1"
						// Store 14 days of full backups in the bucket.
						"\(BucketRepoName)-retention-full":      string | *"14"
						"\(BucketRepoName)-retention-full-type": "count" | *"time" // time in days
						// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/backups#encryption
						"\(BucketRepoName)-cipher-type": "aes-256-cbc"
						// "The convention we recommend for setting this variable is /pgbackrest/$NAMESPACE/$CLUSTER_NAME/repoN"
						// Ref: https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/backups#understanding-backup-configuration-and-basic-operations
						"\(BucketRepoName)-path": "/pgbackrest/\(metadata.namespace)/\(metadata.name)/\(manual.repoName)"
					}
					repos: [
						{
							name: "repo1"
							volume: volumeClaimSpec: {
								accessModes: ["ReadWriteOnce"]
								resources: requests: storage: string | *"4Gi"
							}
						},
						{
							name: BucketRepoName
							// Full backup weekly on Sunday at 1am, differntial daily at 1am every day except Sunday.
							schedules: full:         string | *"0 1 * * 0"
							schedules: differential: string | *"0 1 * * 1-6"
							s3: {
								bucket:   _BackupBucket.metadata.name
								region:   _BackupBucket.spec.region
								endpoint: "s3.dualstack.\(region).amazonaws.com"
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
				resources: requests: storage: string | *"20Gi"
			}
			affinity: podAntiAffinity: preferredDuringSchedulingIgnoredDuringExecution: [{
				weight: 1
				podAffinityTerm: {
					topologyKey: "topology.kubernetes.io/zone"
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
					topologyKey: "topology.kubernetes.io/zone"
					labelSelector: matchLabels: {
						"postgres-operator.crunchydata.com/cluster": metadata.name
						"postgres-operator.crunchydata.com/role":    "pgbouncer"
					}
				}
			}]
		}
	}
}
