package holos

#Values: {

	// +docs:section=Global
	// Default values for cert-manager.
	// This is a YAML-formatted file.
	// Declare variables to be passed into your templates.
	global: {
		// Reference to one or more secrets to be used when pulling images
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
		//
		// For example:
		//  imagePullSecrets:
		//    - name: "image-pull-secret"
		imagePullSecrets: []

		// Labels to apply to all resources
		// Please note that this does not add labels to the resources created dynamically by the controllers.
		// For these resources, you have to add the labels in the template in the cert-manager custom resource:
		// eg. podTemplate/ ingressTemplate in ACMEChallengeSolverHTTP01Ingress
		//    ref: https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEChallengeSolverHTTP01Ingress
		// eg. secretTemplate in CertificateSpec
		//    ref: https://cert-manager.io/docs/reference/api-docs/#cert-manager.io/v1.CertificateSpec
		commonLabels: {}

		// The number of old ReplicaSets to retain to allow rollback (If not set, default Kubernetes value is set to 10)
		// +docs:property
		// revisionHistoryLimit: 1
		// Optional priority class to be used for the cert-manager pods
		priorityClassName: ""

		rbac: {
			// Create required ClusterRoles and ClusterRoleBindings for cert-manager
			create: true
			// Aggregate ClusterRoles to Kubernetes default user-facing roles. Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles
			aggregateClusterRoles: true
		}

		podSecurityPolicy: {
			// Create PodSecurityPolicy for cert-manager
			//
			// NOTE: PodSecurityPolicy was deprecated in Kubernetes 1.21 and removed in 1.25
			enabled: false
			// Configure the PodSecurityPolicy to use AppArmor
			useAppArmor: true
		}

		// Set the verbosity of cert-manager. Range of 0 - 6 with 6 being the most verbose.
		logLevel: 2

		leaderElection: {
			// Override the namespace used for the leader election lease
			namespace: string | *"kube-system"
		}
	}

	// The duration that non-leader candidates will wait after observing a
	// leadership renewal until attempting to acquire leadership of a led but
	// unrenewed leader slot. This is effectively the maximum duration that a
	// leader can be stopped before it is replaced by another candidate.
	// +docs:property
	// leaseDuration: 60s
	// The interval between attempts by the acting master to renew a leadership
	// slot before it stops leading. This must be less than or equal to the
	// lease duration.
	// +docs:property
	// renewDeadline: 40s
	// The duration the clients should wait between attempting acquisition and
	// renewal of a leadership.
	// +docs:property
	// retryPeriod: 15s
	// Install the cert-manager CRDs, it is recommended to not use Helm to manage
	// the CRDs
	installCRDs: true | *false

	// +docs:section=Controller
	// Number of replicas of the cert-manager controller to run.
	//
	// The default is 1, but in production you should set this to 2 or 3 to provide high
	// availability.
	//
	// If `replicas > 1` you should also consider setting `podDisruptionBudget.enabled=true`.
	//
	// Note: cert-manager uses leader election to ensure that there can
	// only be a single instance active at a time.
	replicaCount: 1

	// Deployment update strategy for the cert-manager controller deployment.
	// See https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy
	//
	// For example:
	//  strategy:
	//    type: RollingUpdate
	//    rollingUpdate:
	//      maxSurge: 0
	//      maxUnavailable: 1
	strategy: {}

	podDisruptionBudget: {
		// Enable or disable the PodDisruptionBudget resource
		//
		// This prevents downtime during voluntary disruptions such as during a Node upgrade.
		// For example, the PodDisruptionBudget will block `kubectl drain`
		// if it is used on the Node where the only remaining cert-manager
		// Pod is currently running.
		enabled: false
	}

	// Configures the minimum available pods for disruptions. Can either be set to
	// an integer (e.g. 1) or a percentage value (e.g. 25%).
	// Cannot be used if `maxUnavailable` is set.
	// +docs:property
	// minAvailable: 1
	// Configures the maximum unavailable pods for disruptions. Can either be set to
	// an integer (e.g. 1) or a percentage value (e.g. 25%).
	// Cannot be used if `minAvailable` is set.
	// +docs:property
	// maxUnavailable: 1
	// Comma separated list of feature gates that should be enabled on the
	// controller pod.
	featureGates: ""

	// The maximum number of challenges that can be scheduled as 'processing' at once
	maxConcurrentChallenges: 60

	image: {
		// The container registry to pull the manager image from
		// +docs:property
		// registry: quay.io
		// The container image for the cert-manager controller
		// +docs:property
		repository: "quay.io/jetstack/cert-manager-controller"

		// Override the image tag to deploy by setting this variable.
		// If no value is set, the chart's appVersion will be used.
		// +docs:property
		// tag: vX.Y.Z
		// Setting a digest will override any tag
		// +docs:property
		// digest: sha256:0e072dddd1f7f8fc8909a2ca6f65e76c5f0d2fcfb8be47935ae3457e8bbceb20
		// Kubernetes imagePullPolicy on Deployment.
		pullPolicy: "IfNotPresent"
	}

	// Override the namespace used to store DNS provider credentials etc. for ClusterIssuer
	// resources. By default, the same namespace as cert-manager is deployed within is
	// used. This namespace will not be automatically created by the Helm chart.
	clusterResourceNamespace: ""

	// This namespace allows you to define where the services will be installed into
	// if not set then they will use the namespace of the release
	// This is helpful when installing cert manager as a chart dependency (sub chart)
	namespace: ""

	serviceAccount: {
		// Specifies whether a service account should be created
		create: true

		// The name of the service account to use.
		// If not set and create is true, a name is generated using the fullname template
		// +docs:property
		// name: ""
		// Optional additional annotations to add to the controller's ServiceAccount
		// +docs:property
		// annotations: {}
		// Optional additional labels to add to the controller's ServiceAccount
		// +docs:property
		// labels: {}
		// Automount API credentials for a Service Account.
		automountServiceAccountToken: true
	}

	// Automounting API credentials for a particular pod
	// +docs:property
	// automountServiceAccountToken: true
	// When this flag is enabled, secrets will be automatically removed when the certificate resource is deleted
	enableCertificateOwnerRef: false

	// Used to configure options for the controller pod.
	// This allows setting options that'd usually be provided via flags.
	// An APIVersion and Kind must be specified in your values.yaml file.
	// Flags will override options that are set here.
	//
	// For example:
	//  config:
	//    apiVersion: controller.config.cert-manager.io/v1alpha1
	//    kind: ControllerConfiguration
	//    logging:
	//      verbosity: 2
	//      format: text
	//    leaderElectionConfig:
	//      namespace: kube-system
	//    kubernetesAPIQPS: 9000
	//    kubernetesAPIBurst: 9000
	//    numberOfConcurrentWorkers: 200
	//    featureGates:
	//      AdditionalCertificateOutputFormats: true
	//      DisallowInsecureCSRUsageDefinition: true
	//      ExperimentalCertificateSigningRequestControllers: true
	//      ExperimentalGatewayAPISupport: true
	//      LiteralCertificateSubject: true
	//      SecretsFilteredCaching: true
	//      ServerSideApply: true
	//      StableCertificateRequestName: true
	//      UseCertificateRequestBasicConstraints: true
	//      ValidateCAA: true
	//    metricsTLSConfig:
	//      dynamic:
	//        secretNamespace: "cert-manager"
	//        secretName: "cert-manager-metrics-ca"
	//        dnsNames:
	//        - cert-manager-metrics
	//        - cert-manager-metrics.cert-manager
	//        - cert-manager-metrics.cert-manager.svc
	config: {}

	// Setting Nameservers for DNS01 Self Check
	// See: https://cert-manager.io/docs/configuration/acme/dns01/#setting-nameservers-for-dns01-self-check
	// Comma separated string with host and port of the recursive nameservers cert-manager should query
	dns01RecursiveNameservers: ""

	// Forces cert-manager to only use the recursive nameservers for verification.
	// Enabling this option could cause the DNS01 self check to take longer due to caching performed by the recursive nameservers
	dns01RecursiveNameserversOnly: false

	// Additional command line flags to pass to cert-manager controller binary.
	// To see all available flags run docker run quay.io/jetstack/cert-manager-controller:<version> --help
	//
	// Use this flag to enable or disable arbitrary controllers, for example, disable the CertificiateRequests approver
	//
	// For example:
	//  extraArgs:
	//    - --controllers=*,-certificaterequests-approver
	extraArgs: []

	// Additional environment variables to pass to cert-manager controller binary.
	extraEnv: []
	// - name: SOME_VAR
	//   value: 'some value'
	// Resources to provide to the cert-manager controller pod
	//
	// For example:
	//  requests:
	//    cpu: 10m
	//    memory: 32Mi
	//
	// ref: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	resources: #PodResources

	// Pod Security Context
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +docs:property
	securityContext: {
		runAsNonRoot: true
		seccompProfile: type: "RuntimeDefault"
	}

	// Container Security Context to be set on the controller component container
	// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +docs:property
	containerSecurityContext: {
		allowPrivilegeEscalation: false
		capabilities: drop: ["ALL"]
		readOnlyRootFilesystem: true
	}

	// Additional volumes to add to the cert-manager controller pod.
	volumes: []

	// Additional volume mounts to add to the cert-manager controller container.
	volumeMounts: []

	// Optional additional annotations to add to the controller Deployment
	// +docs:property
	// deploymentAnnotations: {}
	// Optional additional annotations to add to the controller Pods
	// +docs:property
	// podAnnotations: {}
	// Optional additional labels to add to the controller Pods
	podLabels: {}

	// Optional annotations to add to the controller Service
	// +docs:property
	// serviceAnnotations: {}
	// Optional additional labels to add to the controller Service
	// +docs:property
	// serviceLabels: {}
	// Optional DNS settings, useful if you have a public and private DNS zone for
	// the same domain on Route 53. What follows is an example of ensuring
	// cert-manager can access an ingress or DNS TXT records at all times.
	// NOTE: This requires Kubernetes 1.10 or `CustomPodDNS` feature gate enabled for
	// the cluster to work.
	// Pod DNS policy
	// ref: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy
	// +docs:property
	// podDnsPolicy: "None"
	// Pod DNS config, podDnsConfig field is optional and it can work with any podDnsPolicy
	// settings. However, when a Pod's dnsPolicy is set to "None", the dnsConfig field has to be specified.
	// ref: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-dns-config
	// +docs:property
	// podDnsConfig:
	//   nameservers:
	//     - "1.1.1.1"
	//     - "8.8.8.8"
	// The nodeSelector on Pods tells Kubernetes to schedule Pods on the nodes with
	// matching labels.
	// See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/
	//
	// This default ensures that Pods are only scheduled to Linux nodes.
	// It prevents Pods being scheduled to Windows nodes in a mixed OS cluster.
	// +docs:property
	nodeSelector: #NodeSelector

	// +docs:ignore
	ingressShim: {}

	// Optional default issuer to use for ingress resources
	// +docs:property=ingressShim.defaultIssuerName
	// defaultIssuerName: ""
	// Optional default issuer kind to use for ingress resources
	// +docs:property=ingressShim.defaultIssuerKind
	// defaultIssuerKind: ""
	// Optional default issuer group to use for ingress resources
	// +docs:property=ingressShim.defaultIssuerGroup
	// defaultIssuerGroup: ""
	// Use these variables to configure the HTTP_PROXY environment variables
	// Configures the HTTP_PROXY environment variable for where a HTTP proxy is required
	// +docs:property
	// http_proxy: "http://proxy:8080"
	// Configures the HTTPS_PROXY environment variable for where a HTTP proxy is required
	// +docs:property
	// https_proxy: "https://proxy:8080"
	// Configures the NO_PROXY environment variable for where a HTTP proxy is required,
	// but certain domains should be excluded
	// +docs:property
	// no_proxy: 127.0.0.1,localhost
	// A Kubernetes Affinity, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#affinity-v1-core
	//
	// For example:
	//   affinity:
	//     nodeAffinity:
	//      requiredDuringSchedulingIgnoredDuringExecution:
	//        nodeSelectorTerms:
	//        - matchExpressions:
	//          - key: foo.bar.com/role
	//            operator: In
	//            values:
	//            - master
	affinity: {}

	// A list of Kubernetes Tolerations, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#toleration-v1-core
	//
	// For example:
	//   tolerations:
	//   - key: foo.bar.com/role
	//     operator: Equal
	//     value: master
	//     effect: NoSchedule
	tolerations: []

	// A list of Kubernetes TopologySpreadConstraints, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#topologyspreadconstraint-v1-core
	//
	// For example:
	//   topologySpreadConstraints:
	//   - maxSkew: 2
	//     topologyKey: topology.kubernetes.io/zone
	//     whenUnsatisfiable: ScheduleAnyway
	//     labelSelector:
	//       matchLabels:
	//         app.kubernetes.io/instance: cert-manager
	//         app.kubernetes.io/component: controller
	topologySpreadConstraints: []

	// LivenessProbe settings for the controller container of the controller Pod.
	//
	// Enabled by default, because we want to enable the clock-skew liveness probe that
	// restarts the controller in case of a skew between the system clock and the monotonic clock.
	// LivenessProbe durations and thresholds are based on those used for the Kubernetes
	// controller-manager. See:
	// https://github.com/kubernetes/kubernetes/blob/806b30170c61a38fedd54cc9ede4cd6275a1ad3b/cmd/kubeadm/app/util/staticpod/utils.go#L241-L245
	// +docs:property
	livenessProbe: {
		enabled:             true
		initialDelaySeconds: 10
		periodSeconds:       10
		timeoutSeconds:      15
		successThreshold:    1
		failureThreshold:    8
	}

	// enableServiceLinks indicates whether information about services should be
	// injected into pod's environment variables, matching the syntax of Docker
	// links.
	enableServiceLinks: false

	// +docs:section=Prometheus

	prometheus: {
		// Enable Prometheus monitoring for the cert-manager controller to use with the
		// Prometheus Operator. If this option is enabled without enabling `prometheus.servicemonitor.enabled` or
		// `prometheus.podmonitor.enabled`, 'prometheus.io' annotations are added to the cert-manager Deployment
		// resources. Additionally, a service is created which can be used together
		// with your own ServiceMonitor (managed outside of this Helm chart).
		// Otherwise, a ServiceMonitor/ PodMonitor is created.
		enabled: true
		servicemonitor: {
			// Create a ServiceMonitor to add cert-manager to Prometheus
			enabled: true | *false

			// Specifies the `prometheus` label on the created ServiceMonitor, this is
			// used when different Prometheus instances have label selectors matching
			// different ServiceMonitors.
			prometheusInstance: "default"

			// The target port to set on the ServiceMonitor, should match the port that
			// cert-manager controller is listening on for metrics
			targetPort: 9402

			// The path to scrape for metrics
			path: "/metrics"

			// The interval to scrape metrics
			interval: "60s"

			// The timeout before a metrics scrape fails
			scrapeTimeout: "30s"

			// Additional labels to add to the ServiceMonitor
			labels: {}

			// Additional annotations to add to the ServiceMonitor
			annotations: {}

			// Keep labels from scraped data, overriding server-side labels.
			honorLabels: false

			// EndpointAdditionalProperties allows setting additional properties on the
			// endpoint such as relabelings, metricRelabelings etc.
			//
			// For example:
			//  endpointAdditionalProperties:
			//   relabelings:
			//   - action: replace
			//     sourceLabels:
			//     - __meta_kubernetes_pod_node_name
			//     targetLabel: instance
			//
			// +docs:property
			endpointAdditionalProperties: {}
		}

		// Note: Enabling both PodMonitor and ServiceMonitor is mutually exclusive, enabling both will result in a error.
		podmonitor: {
			// Create a PodMonitor to add cert-manager to Prometheus
			enabled: false

			// Specifies the `prometheus` label on the created PodMonitor, this is
			// used when different Prometheus instances have label selectors matching
			// different PodMonitor.
			prometheusInstance: "default"

			// The path to scrape for metrics
			path: "/metrics"

			// The interval to scrape metrics
			interval: "60s"

			// The timeout before a metrics scrape fails
			scrapeTimeout: "30s"

			// Additional labels to add to the PodMonitor
			labels: {}

			// Additional annotations to add to the PodMonitor
			annotations: {}

			// Keep labels from scraped data, overriding server-side labels.
			honorLabels: false

			// EndpointAdditionalProperties allows setting additional properties on the
			// endpoint such as relabelings, metricRelabelings etc.
			//
			// For example:
			//  endpointAdditionalProperties:
			//   relabelings:
			//   - action: replace
			//     sourceLabels:
			//     - __meta_kubernetes_pod_node_name
			//     targetLabel: instance
			//
			// +docs:property
			endpointAdditionalProperties: {}
		}
	}

	// +docs:section=Webhook

	webhook: {
		// Number of replicas of the cert-manager webhook to run.
		//
		// The default is 1, but in production you should set this to 2 or 3 to provide high
		// availability.
		//
		// If `replicas > 1` you should also consider setting `webhook.podDisruptionBudget.enabled=true`.
		replicaCount: 1

		// Seconds the API server should wait for the webhook to respond before treating the call as a failure.
		// Value must be between 1 and 30 seconds. See:
		// https://kubernetes.io/docs/reference/kubernetes-api/extend-resources/validating-webhook-configuration-v1/
		//
		// We set the default to the maximum value of 30 seconds. Here's why:
		// Users sometimes report that the connection between the K8S API server and
		// the cert-manager webhook server times out.
		// If *this* timeout is reached, the error message will be "context deadline exceeded",
		// which doesn't help the user diagnose what phase of the HTTPS connection timed out.
		// For example, it could be during DNS resolution, TCP connection, TLS
		// negotiation, HTTP negotiation, or slow HTTP response from the webhook
		// server.
		// So by setting this timeout to its maximum value the underlying timeout error
		// message has more chance of being returned to the end user.
		timeoutSeconds: 30

		// Used to configure options for the webhook pod.
		// This allows setting options that'd usually be provided via flags.
		// An APIVersion and Kind must be specified in your values.yaml file.
		// Flags will override options that are set here.
		//
		// For example:
		//  apiVersion: webhook.config.cert-manager.io/v1alpha1
		//  kind: WebhookConfiguration
		//  # The port that the webhook should listen on for requests.
		//  # In GKE private clusters, by default kubernetes apiservers are allowed to
		//  # talk to the cluster nodes only on 443 and 10250. so configuring
		//  # securePort: 10250, will work out of the box without needing to add firewall
		//  # rules or requiring NET_BIND_SERVICE capabilities to bind port numbers < 1000.
		//  # This should be uncommented and set as a default by the chart once we graduate
		//  # the apiVersion of WebhookConfiguration past v1alpha1.
		//  securePort: 10250
		config: {}

		// Deployment update strategy for the cert-manager webhook deployment.
		// See https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy
		//
		// For example:
		//  strategy:
		//    type: RollingUpdate
		//    rollingUpdate:
		//      maxSurge: 0
		//      maxUnavailable: 1
		strategy: {}

		// Pod Security Context to be set on the webhook component Pod
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
		// +docs:property
		securityContext: {
			runAsNonRoot: true
			seccompProfile: type: "RuntimeDefault"
		}

		// Container Security Context to be set on the webhook component container
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
		// +docs:property
		containerSecurityContext: {
			allowPrivilegeEscalation: false
			capabilities: drop: ["ALL"]
			readOnlyRootFilesystem: true
		}

		podDisruptionBudget: {
			// Enable or disable the PodDisruptionBudget resource
			//
			// This prevents downtime during voluntary disruptions such as during a Node upgrade.
			// For example, the PodDisruptionBudget will block `kubectl drain`
			// if it is used on the Node where the only remaining cert-manager
			// Pod is currently running.
			enabled: false
		}

		// Configures the minimum available pods for disruptions. Can either be set to
		// an integer (e.g. 1) or a percentage value (e.g. 25%).
		// Cannot be used if `maxUnavailable` is set.
		// +docs:property
		// minAvailable: 1

		// Configures the maximum unavailable pods for disruptions. Can either be set to
		// an integer (e.g. 1) or a percentage value (e.g. 25%).
		// Cannot be used if `minAvailable` is set.
		// +docs:property
		// maxUnavailable: 1

		// Optional additional annotations to add to the webhook Deployment
		// +docs:property
		// deploymentAnnotations: {}

		// Optional additional annotations to add to the webhook Pods
		// +docs:property
		// podAnnotations: {}

		// Optional additional annotations to add to the webhook Service
		// +docs:property
		// serviceAnnotations: {}

		// Optional additional annotations to add to the webhook MutatingWebhookConfiguration
		// +docs:property
		// mutatingWebhookConfigurationAnnotations: {}

		// Optional additional annotations to add to the webhook ValidatingWebhookConfiguration
		// +docs:property
		// validatingWebhookConfigurationAnnotations: {}

		validatingWebhookConfiguration: {
			// Configure spec.namespaceSelector for validating webhooks.
			// +docs:property
			namespaceSelector: {
				matchExpressions: [{
					key:      "cert-manager.io/disable-validation"
					operator: "NotIn"
					values: ["true"]
				}]
			}
		}

		mutatingWebhookConfiguration: {
			// Configure spec.namespaceSelector for mutating webhooks.
			// +docs:property
			namespaceSelector: {}
		}
		//  matchLabels:
		//    key: value
		//  matchExpressions:
		//    - key: kubernetes.io/metadata.name
		//      operator: NotIn
		//      values:
		//        - kube-system
		// Additional command line flags to pass to cert-manager webhook binary.
		// To see all available flags run docker run quay.io/jetstack/cert-manager-webhook:<version> --help
		extraArgs: []
		// Path to a file containing a WebhookConfiguration object used to configure the webhook
		// - --config=<path-to-config-file>
		// Comma separated list of feature gates that should be enabled on the
		// webhook pod.
		featureGates: ""

		// Resources to provide to the cert-manager webhook pod
		//
		// For example:
		//  requests:
		//    cpu: 10m
		//    memory: 32Mi
		//
		// ref: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
		resources: #PodResources

		// Liveness probe values
		// ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes
		//
		// +docs:property
		livenessProbe: {
			failureThreshold:    3
			initialDelaySeconds: 60
			periodSeconds:       10
			successThreshold:    1
			timeoutSeconds:      1
		}

		// Readiness probe values
		// ref: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes
		//
		// +docs:property
		readinessProbe: {
			failureThreshold:    3
			initialDelaySeconds: 5
			periodSeconds:       5
			successThreshold:    1
			timeoutSeconds:      1
		}

		// The nodeSelector on Pods tells Kubernetes to schedule Pods on the nodes with
		// matching labels.
		// See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/
		//
		// This default ensures that Pods are only scheduled to Linux nodes.
		// It prevents Pods being scheduled to Windows nodes in a mixed OS cluster.
		// +docs:property
		nodeSelector: #NodeSelector

		// A Kubernetes Affinity, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#affinity-v1-core
		//
		// For example:
		//   affinity:
		//     nodeAffinity:
		//      requiredDuringSchedulingIgnoredDuringExecution:
		//        nodeSelectorTerms:
		//        - matchExpressions:
		//          - key: foo.bar.com/role
		//            operator: In
		//            values:
		//            - master
		affinity: {}

		// A list of Kubernetes Tolerations, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#toleration-v1-core
		//
		// For example:
		//   tolerations:
		//   - key: foo.bar.com/role
		//     operator: Equal
		//     value: master
		//     effect: NoSchedule
		tolerations: []

		// A list of Kubernetes TopologySpreadConstraints, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#topologyspreadconstraint-v1-core
		//
		// For example:
		//   topologySpreadConstraints:
		//   - maxSkew: 2
		//     topologyKey: topology.kubernetes.io/zone
		//     whenUnsatisfiable: ScheduleAnyway
		//     labelSelector:
		//       matchLabels:
		//         app.kubernetes.io/instance: cert-manager
		//         app.kubernetes.io/component: controller
		topologySpreadConstraints: []

		// Optional additional labels to add to the Webhook Pods
		podLabels: {}

		// Optional additional labels to add to the Webhook Service
		serviceLabels: {}

		image: {
			// The container registry to pull the webhook image from
			// +docs:property
			// registry: quay.io
			// The container image for the cert-manager webhook
			// +docs:property
			repository: "quay.io/jetstack/cert-manager-webhook"

			// Override the image tag to deploy by setting this variable.
			// If no value is set, the chart's appVersion will be used.
			// +docs:property
			// tag: vX.Y.Z
			// Setting a digest will override any tag
			// +docs:property
			// digest: sha256:0e072dddd1f7f8fc8909a2ca6f65e76c5f0d2fcfb8be47935ae3457e8bbceb20
			// Kubernetes imagePullPolicy on Deployment.
			pullPolicy: "IfNotPresent"
		}

		serviceAccount: {
			// Specifies whether a service account should be created
			create: true

			// The name of the service account to use.
			// If not set and create is true, a name is generated using the fullname template
			// +docs:property
			// name: ""
			// Optional additional annotations to add to the controller's ServiceAccount
			// +docs:property
			// annotations: {}
			// Optional additional labels to add to the webhook's ServiceAccount
			// +docs:property
			// labels: {}
			// Automount API credentials for a Service Account.
			automountServiceAccountToken: true
		}

		// Automounting API credentials for a particular pod
		// +docs:property
		// automountServiceAccountToken: true
		// The port that the webhook should listen on for requests.
		// In GKE private clusters, by default kubernetes apiservers are allowed to
		// talk to the cluster nodes only on 443 and 10250. so configuring
		// securePort: 10250, will work out of the box without needing to add firewall
		// rules or requiring NET_BIND_SERVICE capabilities to bind port numbers <1000
		securePort: 10250

		// Specifies if the webhook should be started in hostNetwork mode.
		//
		// Required for use in some managed kubernetes clusters (such as AWS EKS) with custom
		// CNI (such as calico), because control-plane managed by AWS cannot communicate
		// with pods' IP CIDR and admission webhooks are not working
		//
		// Since the default port for the webhook conflicts with kubelet on the host
		// network, `webhook.securePort` should be changed to an available port if
		// running in hostNetwork mode.
		hostNetwork: false

		// Specifies how the service should be handled. Useful if you want to expose the
		// webhook to outside of the cluster. In some cases, the control plane cannot
		// reach internal services.
		serviceType: "ClusterIP"

		// Specify the load balancer IP for the created service
		// +docs:property
		// loadBalancerIP: "10.10.10.10"
		// Overrides the mutating webhook and validating webhook so they reach the webhook
		// service using the `url` field instead of a service.
		url: {}
		// host:
		// Enables default network policies for webhooks.
		networkPolicy: {
			// Create network policies for the webhooks
			enabled: false

			// Ingress rule for the webhook network policy, by default will allow all
			// inbound traffic
			// +docs:property
			ingress: [{
				from: [{
					ipBlock: cidr: "0.0.0.0/0"
				}]
			}]

			// Egress rule for the webhook network policy, by default will allow all
			// outbound traffic traffic to ports 80 and 443, as well as DNS ports
			// +docs:property
			egress: [{
				ports: [{
					port:     80
					protocol: "TCP"
				}, {
					port:     443
					protocol: "TCP"
				}, {
					port:     53
					protocol: "TCP"
				}, {
					port:     53
					protocol: "UDP"
				}, {
					// On OpenShift and OKD, the Kubernetes API server listens on
					// port 6443.
					port:     6443
					protocol: "TCP"
				}]
				to: [{
					ipBlock: cidr: "0.0.0.0/0"
				}]
			}]
		}

		// Additional volumes to add to the cert-manager controller pod.
		volumes: []

		// Additional volume mounts to add to the cert-manager controller container.
		volumeMounts: []

		// enableServiceLinks indicates whether information about services should be
		// injected into pod's environment variables, matching the syntax of Docker
		// links.
		enableServiceLinks: false
	}

	// +docs:section=CA Injector

	cainjector: {
		// Create the CA Injector deployment
		enabled: true

		// Number of replicas of the cert-manager cainjector to run.
		//
		// The default is 1, but in production you should set this to 2 or 3 to provide high
		// availability.
		//
		// If `replicas > 1` you should also consider setting `cainjector.podDisruptionBudget.enabled=true`.
		//
		// Note: cert-manager uses leader election to ensure that there can
		// only be a single instance active at a time.
		replicaCount: 1

		// Used to configure options for the cainjector pod.
		// This allows setting options that'd usually be provided via flags.
		// An APIVersion and Kind must be specified in your values.yaml file.
		// Flags will override options that are set here.
		//
		// For example:
		//  apiVersion: cainjector.config.cert-manager.io/v1alpha1
		//  kind: CAInjectorConfiguration
		//  logging:
		//   verbosity: 2
		//   format: text
		//  leaderElectionConfig:
		//   namespace: kube-system
		config: {}

		// Deployment update strategy for the cert-manager cainjector deployment.
		// See https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy
		//
		// For example:
		//  strategy:
		//    type: RollingUpdate
		//    rollingUpdate:
		//      maxSurge: 0
		//      maxUnavailable: 1
		strategy: {}

		// Pod Security Context to be set on the cainjector component Pod
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
		// +docs:property
		securityContext: {
			runAsNonRoot: true
			seccompProfile: type: "RuntimeDefault"
		}

		// Container Security Context to be set on the cainjector component container
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
		// +docs:property
		containerSecurityContext: {
			allowPrivilegeEscalation: false
			capabilities: drop: ["ALL"]
			readOnlyRootFilesystem: true
		}

		podDisruptionBudget: {
			// Enable or disable the PodDisruptionBudget resource
			//
			// This prevents downtime during voluntary disruptions such as during a Node upgrade.
			// For example, the PodDisruptionBudget will block `kubectl drain`
			// if it is used on the Node where the only remaining cert-manager
			// Pod is currently running.
			enabled: false
		}

		// Configures the minimum available pods for disruptions. Can either be set to
		// an integer (e.g. 1) or a percentage value (e.g. 25%).
		// Cannot be used if `maxUnavailable` is set.
		// +docs:property
		// minAvailable: 1
		// Configures the maximum unavailable pods for disruptions. Can either be set to
		// an integer (e.g. 1) or a percentage value (e.g. 25%).
		// Cannot be used if `minAvailable` is set.
		// +docs:property
		// maxUnavailable: 1
		// Optional additional annotations to add to the cainjector Deployment
		// +docs:property
		// deploymentAnnotations: {}
		// Optional additional annotations to add to the cainjector Pods
		// +docs:property
		// podAnnotations: {}
		// Additional command line flags to pass to cert-manager cainjector binary.
		// To see all available flags run docker run quay.io/jetstack/cert-manager-cainjector:<version> --help
		extraArgs: []
		// Enable profiling for cainjector
		// - --enable-profiling=true
		// Comma separated list of feature gates that should be enabled on the
		// cainjector pod.
		featureGates: ""

		// Resources to provide to the cert-manager cainjector pod
		//
		// For example:
		//  requests:
		//    cpu: 10m
		//    memory: 32Mi
		//
		// ref: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
		resources: #PodResources

		// The nodeSelector on Pods tells Kubernetes to schedule Pods on the nodes with
		// matching labels.
		// See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/
		//
		// This default ensures that Pods are only scheduled to Linux nodes.
		// It prevents Pods being scheduled to Windows nodes in a mixed OS cluster.
		// +docs:property
		nodeSelector: #NodeSelector

		// A Kubernetes Affinity, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#affinity-v1-core
		//
		// For example:
		//   affinity:
		//     nodeAffinity:
		//      requiredDuringSchedulingIgnoredDuringExecution:
		//        nodeSelectorTerms:
		//        - matchExpressions:
		//          - key: foo.bar.com/role
		//            operator: In
		//            values:
		//            - master
		affinity: {}

		// A list of Kubernetes Tolerations, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#toleration-v1-core
		//
		// For example:
		//   tolerations:
		//   - key: foo.bar.com/role
		//     operator: Equal
		//     value: master
		//     effect: NoSchedule
		tolerations: []

		// A list of Kubernetes TopologySpreadConstraints, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#topologyspreadconstraint-v1-core
		//
		// For example:
		//   topologySpreadConstraints:
		//   - maxSkew: 2
		//     topologyKey: topology.kubernetes.io/zone
		//     whenUnsatisfiable: ScheduleAnyway
		//     labelSelector:
		//       matchLabels:
		//         app.kubernetes.io/instance: cert-manager
		//         app.kubernetes.io/component: controller
		topologySpreadConstraints: []

		// Optional additional labels to add to the CA Injector Pods
		podLabels: {}

		image: {
			// The container registry to pull the cainjector image from
			// +docs:property
			// registry: quay.io
			// The container image for the cert-manager cainjector
			// +docs:property
			repository: "quay.io/jetstack/cert-manager-cainjector"

			// Override the image tag to deploy by setting this variable.
			// If no value is set, the chart's appVersion will be used.
			// +docs:property
			// tag: vX.Y.Z
			// Setting a digest will override any tag
			// +docs:property
			// digest: sha256:0e072dddd1f7f8fc8909a2ca6f65e76c5f0d2fcfb8be47935ae3457e8bbceb20
			// Kubernetes imagePullPolicy on Deployment.
			pullPolicy: "IfNotPresent"
		}

		serviceAccount: {
			// Specifies whether a service account should be created
			create: true

			// The name of the service account to use.
			// If not set and create is true, a name is generated using the fullname template
			// +docs:property
			// name: ""
			// Optional additional annotations to add to the controller's ServiceAccount
			// +docs:property
			// annotations: {}
			// Optional additional labels to add to the cainjector's ServiceAccount
			// +docs:property
			// labels: {}
			// Automount API credentials for a Service Account.
			automountServiceAccountToken: true
		}

		// Automounting API credentials for a particular pod
		// +docs:property
		// automountServiceAccountToken: true
		// Additional volumes to add to the cert-manager controller pod.
		volumes: []

		// Additional volume mounts to add to the cert-manager controller container.
		volumeMounts: []

		// enableServiceLinks indicates whether information about services should be
		// injected into pod's environment variables, matching the syntax of Docker
		// links.
		enableServiceLinks: false
	}

	// +docs:section=ACME Solver

	acmesolver: {
		image: {
			// The container registry to pull the acmesolver image from
			// +docs:property
			// registry: quay.io
			// The container image for the cert-manager acmesolver
			// +docs:property
			repository: "quay.io/jetstack/cert-manager-acmesolver"

			// Override the image tag to deploy by setting this variable.
			// If no value is set, the chart's appVersion will be used.
			// +docs:property
			// tag: vX.Y.Z
			// Setting a digest will override any tag
			// +docs:property
			// digest: sha256:0e072dddd1f7f8fc8909a2ca6f65e76c5f0d2fcfb8be47935ae3457e8bbceb20
			// Kubernetes imagePullPolicy on Deployment.
			pullPolicy: "IfNotPresent"
		}
	}

	// +docs:section=Startup API Check
	// This startupapicheck is a Helm post-install hook that waits for the webhook
	// endpoints to become available.
	// The check is implemented using a Kubernetes Job - if you are injecting mesh
	// sidecar proxies into cert-manager pods, you probably want to ensure that they
	// are not injected into this Job's pod. Otherwise the installation may time out
	// due to the Job never being completed because the sidecar proxy does not exit.
	// See https://github.com/cert-manager/cert-manager/pull/4414 for context.

	startupapicheck: {
		// Enables the startup api check
		enabled: *true | false

		// Pod Security Context to be set on the startupapicheck component Pod
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
		// +docs:property
		securityContext: {
			runAsNonRoot: true
			seccompProfile: type: "RuntimeDefault"
		}

		// Container Security Context to be set on the controller component container
		// ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
		// +docs:property
		containerSecurityContext: {
			allowPrivilegeEscalation: false
			capabilities: drop: ["ALL"]
			readOnlyRootFilesystem: true
		}

		// Timeout for 'kubectl check api' command
		timeout: "1m"

		// Job backoffLimit
		backoffLimit: 4

		// Optional additional annotations to add to the startupapicheck Job
		// +docs:property
		jobAnnotations: {
			"helm.sh/hook":               "post-install"
			"helm.sh/hook-weight":        "1"
			"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded"
		}

		// Optional additional annotations to add to the startupapicheck Pods
		// +docs:property
		// podAnnotations: {}
		// Additional command line flags to pass to startupapicheck binary.
		// To see all available flags run docker run quay.io/jetstack/cert-manager-ctl:<version> --help
		//
		// We enable verbose logging by default so that if startupapicheck fails, users
		// can know what exactly caused the failure. Verbose logs include details of
		// the webhook URL, IP address and TCP connect errors for example.
		// +docs:property
		extraArgs: ["-v"]

		// Resources to provide to the cert-manager controller pod
		//
		// For example:
		//  requests:
		//    cpu: 10m
		//    memory: 32Mi
		//
		// ref: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
		resources: #PodResources

		// The nodeSelector on Pods tells Kubernetes to schedule Pods on the nodes with
		// matching labels.
		// See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/
		//
		// This default ensures that Pods are only scheduled to Linux nodes.
		// It prevents Pods being scheduled to Windows nodes in a mixed OS cluster.
		// +docs:property
		nodeSelector: #NodeSelector

		// A Kubernetes Affinity, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#affinity-v1-core
		//
		// For example:
		//   affinity:
		//     nodeAffinity:
		//      requiredDuringSchedulingIgnoredDuringExecution:
		//        nodeSelectorTerms:
		//        - matchExpressions:
		//          - key: foo.bar.com/role
		//            operator: In
		//            values:
		//            - master
		affinity: {}

		// A list of Kubernetes Tolerations, if required; see https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#toleration-v1-core
		//
		// For example:
		//   tolerations:
		//   - key: foo.bar.com/role
		//     operator: Equal
		//     value: master
		//     effect: NoSchedule
		tolerations: []

		// Optional additional labels to add to the startupapicheck Pods
		podLabels: {}

		image: {
			// The container registry to pull the startupapicheck image from
			// +docs:property
			// registry: quay.io
			// The container image for the cert-manager startupapicheck
			// +docs:property
			repository: "quay.io/jetstack/cert-manager-startupapicheck"

			// Override the image tag to deploy by setting this variable.
			// If no value is set, the chart's appVersion will be used.
			// +docs:property
			// tag: vX.Y.Z
			// Setting a digest will override any tag
			// +docs:property
			// digest: sha256:0e072dddd1f7f8fc8909a2ca6f65e76c5f0d2fcfb8be47935ae3457e8bbceb20
			// Kubernetes imagePullPolicy on Deployment.
			pullPolicy: "IfNotPresent"
		}

		rbac: {
			// annotations for the startup API Check job RBAC and PSP resources
			// +docs:property
			annotations: {
				"helm.sh/hook":               "post-install"
				"helm.sh/hook-weight":        "-5"
				"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded"
			}
		}

		// Automounting API credentials for a particular pod
		// +docs:property
		// automountServiceAccountToken: true

		serviceAccount: {
			// Specifies whether a service account should be created
			create: true

			// The name of the service account to use.
			// If not set and create is true, a name is generated using the fullname template
			// +docs:property
			// name: ""
			// Optional additional annotations to add to the Job's ServiceAccount
			// +docs:property
			annotations: {
				"helm.sh/hook":               "post-install"
				"helm.sh/hook-weight":        "-5"
				"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded"
			}

			// Automount API credentials for a Service Account.
			// +docs:property
			automountServiceAccountToken: true
		}

		// Optional additional labels to add to the startupapicheck's ServiceAccount
		// +docs:property
		// labels: {}
		// Additional volumes to add to the cert-manager controller pod.
		volumes: []

		// Additional volume mounts to add to the cert-manager controller container.
		volumeMounts: []

		// enableServiceLinks indicates whether information about services should be
		// injected into pod's environment variables, matching the syntax of Docker
		// links.
		enableServiceLinks: false
	}
}