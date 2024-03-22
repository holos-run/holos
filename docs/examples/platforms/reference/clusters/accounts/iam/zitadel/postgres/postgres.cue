package holos

#InputKeys: component:        "postgres"
#DependsOn: "postgres-certs": _

let Cluster = #Platform.clusters[#ClusterName]
let S3Secret = "pgo-s3-creds"
let ZitadelUser = _DBName
let ZitadelAdmin = "\(_DBName)-admin"

// This must be an external storage bucket for our architecture.
let BucketRepoName = "repo2"

// Restore options.  Set the timestamp to a known good point in time.
// time="2024-03-11T17:08:58Z" level=info msg="crunchy-pgbackrest ends"
// let RestoreOptions = ["--type=time", "--target=\"2024-03-11 17:10:00+00\""]

// Restore the most recent backup.
let RestoreOptions = []

#Kustomization: spec: healthChecks: [
	{
		apiVersion: "external-secrets.io/v1beta1"
		kind:       "ExternalSecret"
		name:       S3Secret
		namespace:  #TargetNamespace
	},
	{
		apiVersion: "postgres-operator.crunchydata.com/v1beta1"
		kind:       "PostgresCluster"
		name:       _DBName
		namespace:  #TargetNamespace
	},
]

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "prod-iam-postgres"

		_dependsOn: "prod-secrets-namespaces": _
		_dependsOn: "prod-iam-postgres-certs": _
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: "\(S3Secret)": _
		PostgresCluster: db: #PostgresCluster & HighlyAvailable & {
			metadata: name:      _DBName
			metadata: namespace: #TargetNamespace
			spec: {
				image:           "registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-16.2-0"
				postgresVersion: 16
				// Custom certs are necessary for streaming standby replication which we use to replicate between two regions.
				// Refer to https://access.crunchydata.com/documentation/postgres-operator/latest/tutorials/backups-disaster-recovery/disaster-recovery#streaming-standby
				customTLSSecret: name:            "\(_DBName)-primary-tls"
				customReplicationTLSSecret: name: "\(_DBName)-repl-tls"
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
						resources: requests: storage: "10Gi"
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
				// Restore from backup if and only if the cluster is primary
				if Cluster.primary {
					dataSource: pgbackrest: {
						stanza:        "db"
						configuration: backups.pgbackrest.configuration
						// Restore from known good full backup taken
						options: RestoreOptions
						global: {
							"\(BucketRepoName)-path":        "/pgbackrest/\(#TargetNamespace)/\(metadata.name)/\(BucketRepoName)"
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
						"\(BucketRepoName)-path": "/pgbackrest/\(#TargetNamespace)/\(metadata.name)/\(manual.repoName)"
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
				resources: requests: storage: string | *"10Gi"
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
