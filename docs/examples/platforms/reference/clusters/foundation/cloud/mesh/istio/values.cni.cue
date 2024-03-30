package holos

// Default values.yaml imported from the cni chart

#CNIValues: {
	cni: {
		hub:        ""
		tag:        ""
		variant:    ""
		image:      "install-cni"
		pullPolicy: ""

		// Refer to https://istio.io/latest/docs/setup/additional-setup/cni/#installing-with-helm
		enabled: #IstioValues.istio_cni.enabled

		// Configuration log level of istio-cni binary
		// by default istio-cni send all logs to UDS server
		// if want to see them you need change global.logging.level with cni:debug
		logLevel: "debug"

		// Configuration file to insert istio-cni plugin configuration
		// by default this will be the first file found in the cni-conf-dir
		// Example
		// cniConfFileName: 10-calico.conflist
		// CNI bin and conf dir override settings
		// defaults:
		cniBinDir:       "" // Auto-detected based on version; defaults to /opt/cni/bin.
		cniConfDir:      "/etc/cni/net.d"
		cniConfFileName: ""
		// This directory must exist on the node, if it does not, consult your container runtime
		// documentation for the appropriate path.
		cniNetnsDir: null // Defaults to '/var/run/netns', in minikube/docker/others can be '/var/run/docker/netns'.

		excludeNamespaces: [
			"istio-system",
			"kube-system",
		]

		// Allows user to set custom affinity for the DaemonSet
		affinity: {}

		// Custom annotations on pod level, if you need them
		podAnnotations: {}

		// If this value is set a RoleBinding will be created
		// in the same namespace as the istio-cni DaemonSet is created.
		// This can be used to bind a preexisting ClusterRole to the istio/cni ServiceAccount
		// e.g. if you use PodSecurityPolicies
		psp_cluster_role: ""

		// Deploy the config files as plugin chain (value "true") or as standalone files in the conf dir (value "false")?
		// Some k8s flavors (e.g. OpenShift) do not support the chain approach, set to false if this is the case
		chained: #IstioValues.istio_cni.chained

		// Allow the istio-cni container to run in privileged mode, needed for some platforms (e.g. OpenShift) or features (repairPods)
		privileged: false

		// Custom configuration happens based on the CNI provider.
		// Possible values: "default", "multus"
		provider: "default"

		// Configure ambient settings
		ambient: {
			// If enabled, ambient redirection will be enabled
			enabled: false
			// Set ambient redirection mode: "iptables" or "ebpf"
			redirectMode: "iptables"
			// Set ambient config dir path: defaults to /etc/ambient-config
			configDir: ""
		}

		repair: {
			enabled: true
			hub:     ""
			tag:     ""

			// Repair controller has 3 modes. Pick which one meets your use cases. Note only one may be used.
			// This defines the action the controller will take when a pod is detected as broken.
			// labelPods will label all pods with <brokenPodLabelKey>=<brokenPodLabelValue>.
			// This is only capable of identifying broken pods; the user is responsible for fixing them (generally, by deleting them).
			labelPods: false
			// deletePods will delete any broken pod. These will then be rescheduled, hopefully onto a node that is fully ready.
			deletePods: true
			// repairPods will dynamically repair any broken pod by setting up the pod networking configuration even after it has started.
			// Note the pod will be crashlooping, so this may take a few minutes to become fully functional based on when the retry occurs.
			// This requires no RBAC privilege, but does require `securityContext.privileged`.
			repairPods: false

			initContainerName: "istio-validation"

			brokenPodLabelKey:   "cni.istio.io/uninitialized"
			brokenPodLabelValue: "true"
		}

		// Set to `type: RuntimeDefault` to use the default profile if available.
		seccompProfile: {}

		resources: requests: {
			cpu:    "100m"
			memory: "100Mi"
		}

		resourceQuotas: {
			enabled: false
			pods:    5000
		}

		// The number of pods that can be unavailable during rolling update (see
		// `updateStrategy.rollingUpdate.maxUnavailable` here:
		// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/daemon-set-v1/#DaemonSetSpec).
		// May be specified as a number of pods or as a percent of the total number
		// of pods at the start of the update.
		rollingMaxUnavailable: 1
	}

	// Revision is set as 'version' label and part of the resource names when installing multiple control planes.
	revision: ""

	// For Helm compatibility.
	ownerName: ""

	global: {
		// Default hub for Istio images.
		// Releases are published to docker hub under 'istio' project.
		// Dev builds from prow are on gcr.io
		hub: "docker.io/istio"

		// Default tag for Istio images.
		tag: "1.20.3"

		// Variant of the image to use.
		// Currently supported are: [debug, distroless]
		variant: ""

		// Specify image pull policy if default behavior isn't desired.
		// Default behavior: latest images will be Always else IfNotPresent.
		imagePullPolicy: ""

		// change cni scope level to control logging out of istio-cni-node DaemonSet
		logging: {
			level: "default:info,cni:info"
		}

		logAsJson: false

		// ImagePullSecrets for all ServiceAccount, list of secrets in the same namespace
		// to use for pulling any images in pods that reference this ServiceAccount.
		// For components that don't use ServiceAccounts (i.e. grafana, servicegraph, tracing)
		// ImagePullSecrets will be added to the corresponding Deployment(StatefulSet) objects.
		// Must be set for any cluster configured with private docker registry.
		imagePullSecrets: []
		// - private-registry-key
		// Default resources allocated
		defaultResources: {
			requests: {
				cpu:    "100m"
				memory: "100Mi"
			}
		}
	}
}