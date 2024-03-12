package holos

#Values: {

	// Vault Helm Chart Holos Values
	global: {
		enabled: true
		// Istio handles this
		tlsDisable: true
	}

	injector: enabled: false

	server: {
		image: {
			// repository: "hashicorp/vault"
			repository: "quay.io/holos/hashicorp/vault"
			tag:        "1.14.10"
			// Overrides the default Image Pull Policy
			pullPolicy: "IfNotPresent"
		}

		extraLabels: "sidecar.istio.io/inject": "true"

		resources: requests: {
			memory: "256Mi"
			cpu:    "2000m"
		}
		// limits:
		//   memory: 1024Mi
		//   cpu: 2000m
		// For HA configuration and because we need to manually init the vault,
		// we need to define custom readiness/liveness Probe settings
		readinessProbe: {
			enabled: true
			path:    "/v1/sys/health?standbyok=true&sealedcode=204&uninitcode=204"
		}
		livenessProbe: {
			enabled:             true
			path:                "/v1/sys/health?standbyok=true"
			initialDelaySeconds: 60
		}

		// extraEnvironmentVars is a list of extra environment variables to set with
		// the stateful set. These could be used to include variables required for
		// auto-unseal.
		// Vault validates an incomplete chain:
		// https://github.com/hashicorp/vault/issues/11318
		extraEnvironmentVars: {
			GOMAXPROCS: "2"
		} // Set to cpu limit, see https://github.com/uber-go/automaxprocs
		// extraVolumes is a list of extra volumes to mount. These will be exposed
		// to Vault in the path `/vault/userconfig/<name>/`.
		extraVolumes: [{
			type: "secret"
			name: "gcpkms-creds"
		}]

		// This configures the Vault Statefulset to create a PVC for audit logs.
		// See https://www.vaultproject.io/docs/audit/index.html to know more
		auditStorage: {
			enabled:   true
			mountPath: "/var/log/vault"
		} // for compatibility with plain debian vm location.

		standalone: {
			enabled: false
		}

		ha: {
			enabled:  true
			replicas: 3
			raft: {
				enabled:   true
				setNodeId: true

				config: """
					ui = true
					listener \"tcp\" {
					  address = \"[::]:8200\"
					  cluster_address = \"[::]:8201\"
					  # mTLS is handled by the the istio sidecar
					  tls_disable = \"true\"
					  # Enable unauthenticated metrics access (necessary for Prometheus Operator)
					  telemetry {
					    unauthenticated_metrics_access = true
					  }
					}

					telemetry {
					  prometheus_retention_time = \"30s\"
					  disable_hostname = true
					}

					seal \"gcpckms\" {
					  credentials = \"/vault/userconfig/gcpkms-creds/credentials.json\"
					  project     = \"v6-vault-f15f\"
					  region      = \"us-west1\"
					  key_ring    = \"vault-core\"
					  crypto_key  = \"vault-core-unseal\"
					}

					# Note; the retry_join leader_api_address values come from the Stable
					# Network ID feature of a Statefulset.  See:
					# https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#stable-network-id
					storage \"raft\" {
					  path = \"/vault/data\"
					  retry_join {
					    leader_api_addr = \"http://vault-0.vault-internal:8200\"
					    leader_tls_servername = \"vault\"
					  }
					  retry_join {
					    leader_api_addr = \"http://vault-1.vault-internal:8200\"
					    leader_tls_servername = \"vault\"
					  }
					  retry_join {
					    leader_api_addr = \"http://vault-2.vault-internal:8200\"
					    leader_tls_servername = \"vault\"
					  }

					  autopilot {
					    cleanup_dead_servers = \"true\"
					    last_contact_threshold = \"200ms\"
					    last_contact_failure_threshold = \"10m\"
					    max_trailing_logs = 250000
					    min_quorum = 3
					    server_stabilization_time = \"10s\"
					  }

					}

					service_registration \"kubernetes\" {}

					"""

				// Vault UI (Will be exposed via the service mesh)
			} // Vault UI (Will be exposed via the service mesh)
		} // Vault UI (Will be exposed via the service mesh)
	} // Vault UI (Will be exposed via the service mesh)// Vault UI (Will be exposed via the service mesh)
	ui: {
		enabled:         true
		serviceType:     "ClusterIP"
		serviceNodePort: null
		externalPort:    8200
	}
}
