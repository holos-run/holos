package holos

// Imported with:
// holos cue import --package holos --path '#ComponentValues:' --outfile values.schema.cue vendor/6.6.2/podinfo/values.yaml
// null values need to be converted to optional fields in cue.
#ComponentValues: {
	// Default values for podinfo.
	replicaCount: 1
	logLevel:     "info"
	//0.0.0.0
	host?: string
	//http://backend-podinfo:9898/echo
	backend?: string
	backends: []
	image: {
		repository: string | *"ghcr.io/stefanprodan/podinfo"
		tag:        string
		pullPolicy: "IfNotPresent"
	}
	ui: {
		color:   "#34577c"
		message: string
		logo:    ""
	}

	// failure conditions
	faults: {
		delay:       false
		error:       false
		unhealthy:   false
		unready:     false
		testFail:    false
		testTimeout: false
	}

	// Kubernetes Service settings
	service: {
		enabled: true
		annotations: {}
		type:         "ClusterIP"
		metricsPort:  9797
		httpPort:     9898
		externalPort: 9898
		grpcPort:     9999
		grpcService:  "podinfo"
		nodePort:     31198
		// the port used to bind the http port to the host
		// NOTE: requires privileged container with NET_BIND_SERVICE capability -- this is useful for testing
		// in local clusters such as kind without port forwarding
		hostPort?: int
	}

	// enable h2c protocol (non-TLS version of HTTP/2)
	h2c: {
		enabled: false
	}

	// config file settings
	config: {
		// config file path
		path: ""
		// config file name
		name: ""
	}

	// Additional command line arguments to pass to podinfo container
	extraArgs: []

	// enable tls on the podinfo service
	tls: {
		enabled: false
		// the name of the secret used to mount the certificate key pair
		secretName?: string
		// the path where the certificate key pair will be mounted
		certPath: "/data/cert"
		// the port used to host the tls endpoint on the service
		port: 9899
		// the port used to bind the tls port to the host
		// NOTE: requires privileged container with NET_BIND_SERVICE capability -- this is useful for testing
		// in local clusters such as kind without port forwarding
		hostPort?: int
	}

	// create a certificate manager certificate (cert-manager required)
	certificate: {
		create: false
		// the issuer used to issue the certificate
		issuerRef: {
			kind: "ClusterIssuer"
			name: "self-signed"
		}
		// the hostname / subject alternative names for the certificate
		dnsNames: ["podinfo"]
	}

	// metrics-server add-on required
	hpa: {
		enabled:     false
		maxReplicas: 10
		// average total CPU usage per pod (1-100)
		cpu?: int
		// average memory usage per pod (100Mi-1Gi)
		memory?: string
		// average http requests per second per pod (k8s-prometheus-adapter)
		requests?: string
	}

	// Redis address in the format tcp://<host>:<port>
	cache: ""
	// Redis deployment
	redis: {
		enabled:    false
		repository: "redis"
		tag:        "7.0.7"
	}
	serviceAccount: {
		// Specifies whether a service account should be created
		enabled: false
		// The name of the service account to use.
		// If not set and create is true, a name is generated using the fullname template
		name?: string
		// List of image pull secrets if pulling from private registries
		imagePullSecrets: []
	}

	// set container security context
	securityContext: {}
	ingress: {
		enabled:   false
		className: ""
		additionalLabels: {}
		annotations: {}
		// kubernetes.io/ingress.class: nginx
		// kubernetes.io/tls-acme: "true"
		hosts: [{
			host: "podinfo.local"
			paths: [{
				path:     "/"
				pathType: "ImplementationSpecific"
			}]
		}]
		//  - secretName: chart-example-tls
		//    hosts:
		//      - chart-example.local
		tls: []
	}
	linkerd: profile: enabled: false

	// create Prometheus Operator monitor
	serviceMonitor: {
		enabled:  false
		interval: "15s"
		additionalLabels: {}
	}
	resources: {
		limits?: {}
		requests: {
			cpu:    "1m"
			memory: "16Mi"
		}
	}

	// Extra environment variables for the podinfo container
	// Example on how to configure extraEnvs
	//  - name: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
	//    value: "http://otel:4317"
	//  - name: MULTIPLE_VALUES
	//    value: TEST
	extraEnvs: []
	nodeSelector: {}
	tolerations: []
	affinity: {}
	podAnnotations: {}

	// https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
	topologySpreadConstraints: []

	// Disruption budget will be configured only when the replicaCount is greater than 1
	//  maxUnavailable: 1
	podDisruptionBudget: {}

	// https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	probes: {
		readiness: {
			initialDelaySeconds: 1
			timeoutSeconds:      5
			failureThreshold:    3
			successThreshold:    1
			periodSeconds:       10
		}
		liveness: {
			initialDelaySeconds: 1
			timeoutSeconds:      5
			failureThreshold:    3
			successThreshold:    1
			periodSeconds:       10
		}
		startup: {
			enable:              false
			initialDelaySeconds: 10
			timeoutSeconds:      5
			failureThreshold:    20
			successThreshold:    1
			periodSeconds:       10
		}
	}
}
