package holos

#Values: {

	// Generated file, DO NOT EDIT. Source: build/templates/values.yaml
	// Overrides the chart name against the label "app.kubernetes.io/name: " placed on every resource this chart creates.
	nameOverride: ""

	// Override the resource names created by this chart which originally is generated using release and chart name.
	fullnameOverride: string | *""

	image: {
		repository: string | *"cockroachdb/cockroach"
		tag:        "v23.1.13"
		pullPolicy: "IfNotPresent"
		credentials: {}
	}
	// registry: docker.io
	// username: john_doe
	// password: changeme
	// Additional labels to apply to all Kubernetes resources created by this chart.
	labels: {}
	// app.kubernetes.io/part-of: my-app
	// Cluster's default DNS domain.
	// You should overwrite it if you're using a different one,
	// otherwise CockroachDB nodes discovery won't work.
	clusterDomain: "cluster.local"

	conf: {
		// An ordered list of CockroachDB node attributes.
		// Attributes are arbitrary strings specifying machine capabilities.
		// Machine capabilities might include specialized hardware or number of cores
		// (e.g. "gpu", "x16c").
		attrs: []
		// - x16c
		// - gpu
		// Total size in bytes for caches, shared evenly if there are multiple
		// storage devices. Size suffixes are supported (e.g. `1GB` and `1GiB`).
		// A percentage of physical memory can also be specified (e.g. `.25`).
		cache: "25%"

		// Sets a name to verify the identity of a cluster.
		// The value must match between all nodes specified via `conf.join`.
		// This can be used as an additional verification when either the node or
		// cluster, or both, have not yet been initialized and do not yet know their
		// cluster ID.
		// To introduce a cluster name into an already-initialized cluster, pair this
		// option with `conf.disable-cluster-name-verification: yes`.
		"cluster-name": ""

		// Tell the server to ignore `conf.cluster-name` mismatches.
		// This is meant for use when opting an existing cluster into starting to use
		// cluster name verification, or when changing the cluster name.
		// The cluster should be restarted once with `conf.cluster-name` and
		// `conf.disable-cluster-name-verification: yes` combined, and once all nodes
		// have been updated to know the new cluster name, the cluster can be restarted
		// again with `conf.disable-cluster-name-verification: no`.
		// This option has no effect if `conf.cluster-name` is not specified.
		"disable-cluster-name-verification": false

		// The addresses for connecting a CockroachDB nodes to an existing cluster.
		// If you are deploying a second CockroachDB instance that should join a first
		// one, use the below list to join to the existing instance.
		// Each item in the array should be a FQDN (and port if needed) resolvable by
		// new Pods.
		join: []

		// New logging configuration.
		log: {
			enabled: false
			// https://www.cockroachlabs.com/docs/v21.1/configure-logs
			config: {}
		}
		// file-defaults:
		//   dir: /custom/dir/path/
		// fluent-defaults:
		//   format: json-fluent
		// sinks:
		//   stderr:
		//     channels: [DEV]
		// Logs at or above this threshold to STDERR. Ignored when "log" is enabled
		logtostderr: "INFO"

		// Maximum storage capacity available to store temporary disk-based data for
		// SQL queries that exceed the memory budget (e.g. join, sorts, etc are
		// sometimes able to spill intermediate results to disk).
		// Accepts numbers interpreted as bytes, size suffixes (e.g. `32GB` and
		// `32GiB`) or a percentage of disk size (e.g. `10%`).
		// The location of the temporary files is within the first store dir.
		// If expressed as a percentage, `max-disk-temp-storage` is interpreted
		// relative to the size of the storage device on which the first store is
		// placed. The temp space usage is never counted towards any store usage
		// (although it does share the device with the first store) so, when
		// configuring this, make sure that the size of this temp storage plus the size
		// of the first store don't exceed the capacity of the storage device.
		// If the first store is an in-memory one (i.e. `type=mem`), then this
		// temporary "disk" data is also kept in-memory.
		// A percentage value is interpreted as a percentage of the available internal
		// memory.
		// max-disk-temp-storage: 0GB
		// Maximum allowed clock offset for the cluster. If observed clock offsets
		// exceed this limit, servers will crash to minimize the likelihood of
		// reading inconsistent data. Increasing this value will increase the time
		// to recovery of failures as well as the frequency of uncertainty-based
		// read restarts.
		// Note, that this value must be the same on all nodes in the cluster.
		// In order to change it, all nodes in the cluster must be stopped
		// simultaneously and restarted with the new value.
		// max-offset: 500ms
		// Maximum memory capacity available to store temporary data for SQL clients,
		// including prepared queries and intermediate data rows during query
		// execution. Accepts numbers interpreted as bytes, size suffixes
		// (e.g. `1GB` and `1GiB`) or a percentage of physical memory (e.g. `.25`).
		"max-sql-memory": "25%"

		// An ordered, comma-separated list of key-value pairs that describe the
		// topography of the machine. Topography might include country, datacenter
		// or rack designations. Data is automatically replicated to maximize
		// diversities of each tier. The order of tiers is used to determine
		// the priority of the diversity, so the more inclusive localities like
		// country should come before less inclusive localities like datacenter.
		// The tiers and order must be the same on all nodes. Including more tiers
		// is better than including fewer. For example:
		//   locality: country=us,region=us-west,datacenter=us-west-1b,rack=12
		//   locality: country=ca,region=ca-east,datacenter=ca-east-2,rack=4
		//   locality: planet=earth,province=manitoba,colo=secondary,power=3
		locality: ""

		// Run CockroachDB instances in standalone mode with replication disabled
		// (replication factor = 1).
		// Enabling this option makes the following values to be ignored:
		// - `conf.cluster-name`
		// - `conf.disable-cluster-name-verification`
		// - `conf.join`
		//
		// WARNING: Enabling this option makes each deployed Pod as a STANDALONE
		//          CockroachDB instance, so the StatefulSet does NOT FORM A CLUSTER.
		//          Don't use this option for production deployments unless you clearly
		//          understand what you're doing.
		//          Usually, this option is intended to be used in conjunction with
		//          `statefulset.replicas: 1` for temporary one-time deployments (like
		//          running E2E tests, for example).
		"single-node": false

		// If non-empty, create a SQL audit log in the specified directory.
		"sql-audit-dir": ""

		// CockroachDB's port to listen to inter-communications and client connections.
		port: 26257

		// CockroachDB's port to listen to HTTP requests.
		"http-port": 8080

		// CockroachDB's data mount path.
		path: "cockroach-data"

		// CockroachDB's storage configuration https://www.cockroachlabs.com/docs/v21.1/cockroach-start.html#storage
		// Uses --store flag
		store: {
			enabled: false
			// Should be empty or 'mem'
			type: null
			// Required for type=mem. If type and size is empty - storage.persistentVolume.size is used
			size: null
			// Arbitrary strings, separated by colons, specifying disk type or capability
			attrs: null
		}
	}

	statefulset: {
		replicas: 3
		updateStrategy: type: "RollingUpdate"
		podManagementPolicy: "Parallel"
		budget: maxUnavailable: 1

		// List of additional command-line arguments you want to pass to the
		// `cockroach start` command.
		args: []
		// - --disable-cluster-name-verification
		// List of extra environment variables to pass into container
		env: []
		// - name: COCKROACH_ENGINE_MAX_SYNC_DURATION
		//   value: "24h"
		// List of Secrets names in the same Namespace as the CockroachDB cluster,
		// which shall be mounted into `/etc/cockroach/secrets/` for every cluster
		// member.
		secretMounts: []

		// Additional labels to apply to this StatefulSet and all its Pods.
		labels: {
			"app.kubernetes.io/component": "cockroachdb"
		}

		// Additional annotations to apply to the Pods of this StatefulSet.
		annotations: {}

		// Affinity rules for scheduling Pods of this StatefulSet on Nodes.
		// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#node-affinity
		nodeAffinity: {}
		// Inter-Pod Affinity rules for scheduling Pods of this StatefulSet.
		// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity
		podAffinity: {}
		// Anti-affinity rules for scheduling Pods of this StatefulSet.
		// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity
		// You may either toggle options below for default anti-affinity rules,
		// or specify the whole set of anti-affinity rules instead of them.
		podAntiAffinity: {
			// The topologyKey to be used.
			// Can be used to spread across different nodes, AZs, regions etc.
			topologyKey: "kubernetes.io/hostname"
			// Type of anti-affinity rules: either `soft`, `hard` or empty value (which
			// disables anti-affinity rules).
			type: "soft"
			// Weight for `soft` anti-affinity rules.
			// Does not apply for other anti-affinity types.
			weight: 100
		}

		// Node selection constraints for scheduling Pods of this StatefulSet.
		// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
		nodeSelector: {}

		// PriorityClassName given to Pods of this StatefulSet
		// https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/#priorityclass
		priorityClassName: ""

		// Taints to be tolerated by Pods of this StatefulSet.
		// https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
		tolerations: []

		// https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
		topologySpreadConstraints: {
			maxSkew:           1
			topologyKey:       "topology.kubernetes.io/zone"
			whenUnsatisfiable: "ScheduleAnyway"
		}

		// Uncomment the following resources definitions or pass them from
		// command line to control the CPU and memory resources allocated
		// by Pods of this StatefulSet.
		resources: {}
		// limits:
		//   cpu: 100m
		//   memory: 512Mi
		// requests:
		//   cpu: 100m
		//   memory: 512Mi
		// Custom Liveness probe
		// https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-http-request
		customLivenessProbe: {}
		// httpGet:
		//   path: /health
		//   port: http
		//   scheme: HTTPS
		// initialDelaySeconds: 30
		// periodSeconds: 5
		// Custom Rediness probe
		// https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-readiness-probes
		customReadinessProbe: {}
		// httpGet:
		//   path: /health
		//   port: http
		//   scheme: HTTPS
		// initialDelaySeconds: 30
		// periodSeconds: 5

		securityContext: {
			enabled: true
		}

		serviceAccount: {
			// Specifies whether this ServiceAccount should be created.
			create: true
			// The name of this ServiceAccount to use.
			// If not set and `create` is `true`, then service account is auto-generated.
			// If not set and `create` is `false`, then it uses default service account.
			name: ""
			// Additional serviceAccount annotations (e.g. for attaching AWS IAM roles to pods)
			annotations: {}
		}
	}

	service: {
		ports: {
			// You can set a different external and internal gRPC ports and their name.
			grpc: {
				external: {
					port: 26257
					name: "grpc"
				}
				// If the port number is different than `external.port`, then it will be
				// named as `internal.name` in Service.
				internal: {
					port: 26257
					// If using Istio set it to `cockroach`.
					name: "grpc-internal"
				}
			}
			http: {
				port: 8080
				name: "http"
			}
		}

		// This Service is meant to be used by clients of the database.
		// It exposes a ClusterIP that will automatically load balance connections
		// to the different database Pods.
		public: {
			type: "ClusterIP"
			// Additional labels to apply to this Service.
			labels: {
				"app.kubernetes.io/component": "cockroachdb"
			}
			// Additional annotations to apply to this Service.
			annotations: {}
		}

		// This service only exists to create DNS entries for each pod in
		// the StatefulSet such that they can resolve each other's IP addresses.
		// It does not create a load-balanced ClusterIP and should not be used directly
		// by clients in most circumstances.
		discovery: {
			// Additional labels to apply to this Service.
			labels: {
				"app.kubernetes.io/component": "cockroachdb"
			}
			// Additional annotations to apply to this Service.
			annotations: {}
		}
	}

	// CockroachDB's ingress for web ui.
	ingress: {
		enabled: false
		labels: {}
		annotations: {}
		//   kubernetes.io/ingress.class: nginx
		//   cert-manager.io/cluster-issuer: letsencrypt
		paths: ["/"]
		hosts: []
		// - cockroachlabs.com
		tls: []
	}
	// - hosts: [cockroachlabs.com]
	//   secretName: cockroachlabs-tls

	prometheus: {
		enabled: true
	}

	securityContext: enabled: true

	// CockroachDB's Prometheus operator ServiceMonitor support
	serviceMonitor: {
		enabled: false
		labels: {}
		annotations: {}
		interval: "10s"
		// scrapeTimeout: 10s
		// Limits the ServiceMonitor to the current namespace if set to `true`.
		namespaced: false

		// tlsConfig: TLS configuration to use when scraping the endpoint.
		// Of type: https://github.com/coreos/prometheus-operator/blob/main/Documentation/api.md#tlsconfig
		tlsConfig: {}
	}

	// CockroachDB's data persistence.
	// If neither `persistentVolume` nor `hostPath` is used, then data will be
	// persisted in ad-hoc `emptyDir`.
	storage: {
		// Absolute path on host to store CockroachDB's data.
		// If not specified, then `emptyDir` will be used instead.
		// If specified, but `persistentVolume.enabled` is `true`, then has no effect.
		hostPath: ""

		// If `enabled` is `true` then a PersistentVolumeClaim will be created and
		// used to store CockroachDB's data, otherwise `hostPath` is used.
		persistentVolume: {
			enabled: true

			size: string | *"100Gi"

			// If defined, then `storageClassName: <storageClass>`.
			// If set to "-", then `storageClassName: ""`, which disables dynamic
			// provisioning.
			// If undefined or empty (default), then no `storageClassName` spec is set,
			// so the default provisioner will be chosen (gp2 on AWS, standard on
			// GKE, AWS & OpenStack).
			storageClass: ""

			// Additional labels to apply to the created PersistentVolumeClaims.
			labels: {}
			// Additional annotations to apply to the created PersistentVolumeClaims.
			annotations: {}
		}
	}

	// Kubernetes Job which initializes multi-node CockroachDB cluster.
	// It's not created if `statefulset.replicas` is `1`.
	init: {
		// Additional labels to apply to this Job and its Pod.
		labels: {
			"app.kubernetes.io/component": "init"
		}

		// Additional annotations to apply to this Job.
		jobAnnotations: {}

		// Additional annotations to apply to the Pod of this Job.
		annotations: {}

		// Affinity rules for scheduling the Pod of this Job.
		// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#node-affinity
		affinity: {}

		// Node selection constraints for scheduling the Pod of this Job.
		// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
		nodeSelector: {}

		// Taints to be tolerated by the Pod of this Job.
		// https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
		tolerations: []

		// The init Pod runs at cluster creation to initialize CockroachDB. It finishes
		// quickly and doesn't continue to consume resources in the Kubernetes
		// cluster. Normally, you should leave this section commented out, but if your
		// Kubernetes cluster uses Resource Quotas and requires all pods to specify
		// resource requests or limits, you can set those here.
		resources: {}
		// requests:
		//   cpu: "10m"
		//   memory: "128Mi"
		// limits:
		//   cpu: "10m"
		//   memory: "128Mi"

		securityContext: {
			enabled: true
		}

		provisioning: {
			enabled: false
			// https://www.cockroachlabs.com/docs/stable/cluster-settings.html
			clusterSettings: null
			// cluster.organization: "'FooCorp - Local Testing'"
			// enterprise.license: "'xxxxx'"
			users: []
			// - name:
			//   password:
			//   # https://www.cockroachlabs.com/docs/stable/create-user.html#parameters
			//   options: [LOGIN]
			databases: []
		}
	}
	// - name:
	//   # https://www.cockroachlabs.com/docs/stable/create-database.html#parameters
	//   options: [encoding='utf-8']
	//   owners: []
	//   # https://www.cockroachlabs.com/docs/stable/grant.html#parameters
	//   owners_with_grant_option: []
	//   # Backup schedules are not idemponent for now and will fail on next run
	//   # https://github.com/cockroachdb/cockroach/issues/57892
	//   backup:
	//     into: s3://
	//     # Enterprise-only option (revision_history)
	//     # https://www.cockroachlabs.com/docs/stable/create-schedule-for-backup.html#backup-options
	//     options: [revision_history]
	//     recurring: '@always'
	//     # Enterprise-only feature. Remove this value to use `FULL BACKUP ALWAYS`
	//     fullBackup: '@daily'
	//     schedule:
	//       # https://www.cockroachlabs.com/docs/stable/create-schedule-for-backup.html#schedule-options
	//       options: [first_run = 'now']
	// Whether to run securely using TLS certificates.
	tls: {
		enabled: true
		copyCerts: image: "busybox"
		certs: {
			// Bring your own certs scenario. If provided, tls.init section will be ignored.
			provided: false
			// Secret name for the client root cert.
			clientRootSecret: "cockroachdb-root"
			// Secret name for node cert.
			nodeSecret: "cockroachdb-node"
			// Secret name for CA cert
			caSecret: "cockroach-ca"
			// Enable if the secret is a dedicated TLS.
			// TLS secrets are created by cert-mananger, for example.
			tlsSecret: false
			// Enable if the you want cockroach db to create its own certificates
			selfSigner: {
				// If set, the cockroach db will generate its own certificates
				enabled: false | *true
				// Run selfSigner as non-root
				securityContext: {
					enabled: true
				}
				// If set, the user should provide the CA certificate to sign other certificates.
				caProvided: false
				// It holds the name of the secret with caCerts. If caProvided is set, this can not be empty.
				caSecret: ""
				// Minimum Certificate duration for all the certificates, all certs duration will be validated against this.
				minimumCertDuration: "624h"
				// Duration of CA certificates in hour
				caCertDuration: "43800h"
				// Expiry window of CA certificates means a window before actual expiry in which CA certs should be rotated.
				caCertExpiryWindow: "648h"
				// Duration of Client certificates in hour
				clientCertDuration: "672h"
				// Expiry window of client certificates means a window before actual expiry in which client certs should be rotated.
				clientCertExpiryWindow: "48h"
				// Duration of node certificates in hour
				nodeCertDuration: "8760h"
				// Expiry window of node certificates means a window before actual expiry in which node certs should be rotated.
				nodeCertExpiryWindow: "168h"
				// If set, the cockroachdb cert selfSigner will rotate the certificates before expiry.
				rotateCerts: true
				// Wait time for each cockroachdb replica to become ready once it comes in running state. Only considered when rotateCerts is set to true
				readinessWait: "30s"
				// Wait time for each cockroachdb replica to get to running state. Only considered when rotateCerts is set to true
				podUpdateTimeout: "2m"
				// ServiceAccount annotations for selfSigner jobs (e.g. for attaching AWS IAM roles to pods)
				svcAccountAnnotations: {}
			}

			// Use cert-manager to issue certificates for mTLS.
			certManager: true | *false
			// Specify an Issuer or a ClusterIssuer to use, when issuing
			// node and client certificates. The values correspond to the
			// issuerRef specified in the certificate.
			certManagerIssuer: {
				group: "cert-manager.io"
				kind:  "Issuer"
				name:  string | *"cockroachdb"
				// Make it false when you are providing your own CA issuer
				isSelfSignedIssuer: true
				// Duration of Client certificates in hours
				clientCertDuration: "672h"
				// Expiry window of client certificates means a window before actual expiry in which client certs should be rotated.
				clientCertExpiryWindow: "48h"
				// Duration of node certificates in hours
				nodeCertDuration: "8760h"
				// Expiry window of node certificates means a window before actual expiry in which node certs should be rotated.
				nodeCertExpiryWindow: "168h"
			}
		}

		selfSigner: {
			// Additional annotations to apply to the Pod of this Job.
			annotations: {}

			// Affinity rules for scheduling the Pod of this Job.
			// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#node-affinity
			affinity: {}

			// Node selection constraints for scheduling the Pod of this Job.
			// https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
			nodeSelector: {}

			// Taints to be tolerated by the Pod of this Job.
			// https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
			tolerations: []

			// Image Placeholder for the selfSigner utility. This will be changed once the CI workflows for the image is in place.
			image: {
				repository: "cockroachlabs-helm-charts/cockroach-self-signer-cert"
				tag:        "1.5"
				pullPolicy: "IfNotPresent"
				credentials: {}
				registry: "gcr.io"
			}
		}
	}
	// username: john_doe
	// password: changeme

	networkPolicy: {
		enabled: false

		ingress: {
			// List of sources which should be able to access the CockroachDB Pods via
			// gRPC port. Items in this list are combined using a logical OR operation.
			// Rules for allowing inter-communication are applied automatically.
			// If empty, then connections from any Pod is allowed.
			grpc: []
			// - podSelector:
			//     matchLabels:
			//       app.kubernetes.io/name: my-app-django
			//       app.kubernetes.io/instance: my-app
			// List of sources which should be able to access the CockroachDB Pods via
			// HTTP port. Items in this list are combined using a logical OR operation.
			// If empty, then connections from any Pod is allowed.
			http: []
		}
	}
	// - namespaceSelector:
	//     matchLabels:
	//       project: my-project
	// To put the admin interface behind Identity Aware Proxy (IAP) on Google Cloud Platform
	// make sure to set ingress.paths: ['/*']
	iap: {
		enabled: false
	}

}
