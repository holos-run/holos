package holos

#Values: {

	// Default values for zitadel.
	zitadel: {
		// The ZITADEL config under configmapConfig is written to a Kubernetes ConfigMap
		// See all defaults here:
		// https://github.com/zitadel/zitadel/blob/main/cmd/defaults.yaml
		configmapConfig: {
			ExternalSecure: true
			Machine: Identification: {
				Hostname: Enabled: true
				Webhook: Enabled:  false
			}
		}

		// The ZITADEL config under secretConfig is written to a Kubernetes Secret
		// See all defaults here:
		// https://github.com/zitadel/zitadel/blob/main/cmd/defaults.yaml
		secretConfig: null

		// Annotations set on secretConfig secret
		secretConfigAnnotations: {
			"helm.sh/hook":               "pre-install,pre-upgrade"
			"helm.sh/hook-delete-policy": "before-hook-creation"
			"helm.sh/hook-weight":        "0"
		}

		// Reference the name of a secret that contains ZITADEL configuration.
		configSecretName: null
		// The key under which the ZITADEL configuration is located in the secret.
		configSecretKey: "config-yaml"

		// ZITADEL uses the masterkey for symmetric encryption.
		// You can generate it for example with tr -dc A-Za-z0-9 </dev/urandom | head -c 32
		masterkey: ""
		// Reference the name of the secret that contains the masterkey. The key should be named "masterkey".
		// Note: Either zitadel.masterkey or zitadel.masterkeySecretName must be set
		masterkeySecretName: string | *""

		// Annotations set on masterkey secret
		masterkeyAnnotations: {
			"helm.sh/hook":               "pre-install,pre-upgrade"
			"helm.sh/hook-delete-policy": "before-hook-creation"
			"helm.sh/hook-weight":        "0"
		}

		// The CA Certificate needed for establishing secure database connections
		dbSslCaCrt: ""

		// The Secret containing the CA certificate at key ca.crt needed for establishing secure database connections
		dbSslCaCrtSecret: string | *""

		// The db admins secret containing the client certificate and key at tls.crt and tls.key needed for establishing secure database connections
		dbSslAdminCrtSecret: string | *""

		// The db users secret containing the client certificate and key at tls.crt and tls.key needed for establishing secure database connections
		dbSslUserCrtSecret: string | *""

		// Generate a self-signed certificate using an init container
		// This will also mount the generated files to /etc/tls/ so that you can reference them in the pod.
		// E.G. KeyPath: /etc/tls/tls.key CertPath: /etc/tls/tls.crt
		// By default, the SAN DNS names include, localhost, the POD IP address and the POD name. You may include one more by using additionalDnsName like "my.zitadel.fqdn".
		selfSignedCert: {
			enabled:           false
			additionalDnsName: null
		}
	}

	replicaCount: 3

	image: {
		repository: "ghcr.io/zitadel/zitadel"
		pullPolicy: "IfNotPresent"
		// Overrides the image tag whose default is the chart appVersion.
		tag: ""
	}

	chownImage: {
		repository: "alpine"
		pullPolicy: "IfNotPresent"
		tag:        "3.19"
	}

	imagePullSecrets: []
	nameOverride:     ""
	fullnameOverride: ""

	// Annotations to add to the deployment
	annotations: {}

	// Annotations to add to the configMap
	configMap: {
		annotations: {
			"helm.sh/hook":               "pre-install,pre-upgrade"
			"helm.sh/hook-delete-policy": "before-hook-creation"
			"helm.sh/hook-weight":        "0"
		}
	}

	serviceAccount: {
		// Specifies whether a service account should be created
		create: true
		// Annotations to add to the service account
		annotations: {
			"helm.sh/hook":               "pre-install,pre-upgrade"
			"helm.sh/hook-delete-policy": "before-hook-creation"
			"helm.sh/hook-weight":        "0"
		}
		// The name of the service account to use.
		// If not set and create is true, a name is generated using the fullname template
		name: ""
	}

	podAnnotations: {}

	podAdditionalLabels: {}

	podSecurityContext: {
		runAsNonRoot: true
		runAsUser:    1000
	}

	securityContext: {}

	// Additional environment variables
	env: []
	// - name: ZITADEL_DATABASE_POSTGRES_HOST
	//   valueFrom:
	//     secretKeyRef:
	//       name: postgres-pguser-postgres
	//       key: host

	service: {
		type: "ClusterIP"
		// If service type is "ClusterIP", this can optionally be set to a fixed IP address.
		clusterIP: ""
		port:      8080
		protocol:  "http2"
		annotations: {}
		scheme: "HTTP"
	}

	ingress: {
		enabled:   false
		className: ""
		annotations: {}
		hosts: [{
			host: "localhost"
			paths: [{
				path:     "/"
				pathType: "Prefix"
			}]
		}]
		tls: []
	}

	resources: {}

	nodeSelector: {}

	tolerations: []

	affinity: {}

	topologySpreadConstraints: []

	initJob: {
		// Once ZITADEL is installed, the initJob can be disabled.
		enabled: true
		annotations: {
			"helm.sh/hook":               "pre-install,pre-upgrade"
			"helm.sh/hook-delete-policy": "before-hook-creation"
			"helm.sh/hook-weight":        "1"
		}
		resources: {}
		backoffLimit:          5
		activeDeadlineSeconds: 300
		extraContainers: []
		podAnnotations: {}
		// Available init commands :
		// "": initialize ZITADEL instance (without skip anything)
		// database: initialize only the database
		// grant: set ALL grant to user
		// user: initialize only the database user
		// zitadel: initialize ZITADEL internals (skip "create user" and "create database")
		command: ""
	}

	setupJob: {
		annotations: {
			"helm.sh/hook":               "pre-install,pre-upgrade"
			"helm.sh/hook-delete-policy": "before-hook-creation"
			"helm.sh/hook-weight":        "2"
		}
		resources: {}
		activeDeadlineSeconds: 300
		extraContainers: []
		podAnnotations: {}
		additionalArgs: ["--init-projections=true"]
		machinekeyWriter: {
			image: {
				repository: "bitnami/kubectl"
				tag:        ""
			}
			resources: {}
		}
	}

	readinessProbe: {
		enabled:             true
		initialDelaySeconds: 0
		periodSeconds:       5
		failureThreshold:    3
	}

	livenessProbe: {
		enabled:             true
		initialDelaySeconds: 0
		periodSeconds:       5
		failureThreshold:    3
	}

	startupProbe: {
		enabled:          true
		periodSeconds:    1
		failureThreshold: 30
	}

	metrics: {
		enabled: false
		serviceMonitor: {
			// If true, the chart creates a ServiceMonitor that is compatible with Prometheus Operator
			// https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#monitoring.coreos.com/v1.ServiceMonitor.
			// The Prometheus community Helm chart installs this operator
			// https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack#kube-prometheus-stack
			enabled:         false
			honorLabels:     false
			honorTimestamps: true
		}
	}

	pdb: {
		enabled: false
		// these values are used for the PDB and are mutally exclusive
		minAvailable: 1
		// maxUnavailable: 1
		annotations: {}
	}
}
