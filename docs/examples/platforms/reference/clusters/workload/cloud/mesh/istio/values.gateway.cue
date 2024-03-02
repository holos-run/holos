package holos

// Gateway default values.yaml imported from the gateway chart.

#GatewayValues: {

	// Name allows overriding the release name. Generally this should not be set
	name: "istio-ingressgateway"
	// revision declares which revision this gateway is a part of
	revision: ""

	// Controls the spec.replicas setting for the Gateway deployment if set.
	// Otherwise defaults to Kubernetes Deployment default (1).
	replicaCount: null

	kind: "Deployment"

	rbac: {
		// If enabled, roles will be created to enable accessing certificates from Gateways. This is not needed
		// when using http://gateway-api.org/.
		enabled: true
	}

	serviceAccount: {
		// If set, a service account will be created. Otherwise, the default is used
		create: true
		// Annotations to add to the service account
		annotations: {}
		// The name of the service account to use.
		// If not set, the release name is used
		name: ""
	}

	podAnnotations: {
		"prometheus.io/port":        "15020"
		"prometheus.io/scrape":      "true"
		"prometheus.io/path":        "/stats/prometheus"
		"inject.istio.io/templates": "gateway"
		"sidecar.istio.io/inject":   "true"
		...
	}

	// Define the security context for the pod.
	// If unset, this will be automatically set to the minimum privileges required to bind to port 80 and 443.
	// On Kubernetes 1.22+, this only requires the `net.ipv4.ip_unprivileged_port_start` sysctl.
	securityContext:          null
	containerSecurityContext: null

	service: {
		// Type of service. Set to "None" to disable the service entirely
		type: string | *"LoadBalancer"
		ports: [...] | *[{
			name:       "status-port"
			port:       15021
			protocol:   "TCP"
			targetPort: 15021
		}, {
			name:       "http2"
			port:       80
			protocol:   "TCP"
			targetPort: 80
		}, {
			name:       "https"
			port:       443
			protocol:   "TCP"
			targetPort: 443
		}]
		annotations: {...}
		loadBalancerIP: ""
		loadBalancerSourceRanges: []
		externalTrafficPolicy: string | *""
		externalIPs: []
		ipFamilyPolicy: ""
		ipFamilies: []
	}

	resources: {
		requests: {
			cpu:    "100m"
			memory: "128Mi"
		}
		limits: {
			cpu:    "2000m"
			memory: "1024Mi"
		}
	}

	autoscaling: {
		enabled:                        true
		minReplicas:                    1
		maxReplicas:                    5
		targetCPUUtilizationPercentage: 80
		autoscaleBehavior: {}
	}

	// Pod environment variables
	env: {}

	// Labels to apply to all resources
	labels: {}

	// Annotations to apply to all resources
	annotations: {}

	nodeSelector: {}

	tolerations: []

	topologySpreadConstraints: []

	affinity: {}

	// If specified, the gateway will act as a network gateway for the given network.
	networkGateway: ""

	// Specify image pull policy if default behavior isn't desired.
	// Default behavior: latest images will be Always else IfNotPresent
	imagePullPolicy: ""

	imagePullSecrets: []

	// This value is used to configure a Kubernetes PodDisruptionBudget for the gateway.
	//
	// By default, the `podDisruptionBudget` is disabled (set to `{}`),
	// which means that no PodDisruptionBudget resource will be created.
	//
	// To enable the PodDisruptionBudget, configure it by specifying the
	// `minAvailable` or `maxUnavailable`. For example, to set the
	// minimum number of available replicas to 1, you can update this value as follows:
	//
	// podDisruptionBudget:
	//   minAvailable: 1
	//
	// Or, to allow a maximum of 1 unavailable replica, you can set:
	//
	// podDisruptionBudget:
	//   maxUnavailable: 1
	//
	// You can also specify the `unhealthyPodEvictionPolicy` field, and the valid values are `IfHealthyBudget` and `AlwaysAllow`.
	// For example, to set the `unhealthyPodEvictionPolicy` to `AlwaysAllow`, you can update this value as follows:
	//
	// podDisruptionBudget:
	//   minAvailable: 1
	//   unhealthyPodEvictionPolicy: AlwaysAllow
	//
	// To disable the PodDisruptionBudget, you can leave it as an empty object `{}`:
	//
	// podDisruptionBudget: {}
	//
	podDisruptionBudget: {}

	terminationGracePeriodSeconds: 30

	// A list of `Volumes` added into the Gateway Pods. See
	// https://kubernetes.io/docs/concepts/storage/volumes/.
	volumes: []

	// A list of `VolumeMounts` added into the Gateway Pods. See
	// https://kubernetes.io/docs/concepts/storage/volumes/.
	volumeMounts: []

	// Configure this to a higher priority class in order to make sure your Istio gateway pods
	// will not be killed because of low priority class.
	// Refer to https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/#priorityclass
	// for more detail.
	priorityClassName: ""
}
