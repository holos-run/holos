package holos

// https://cert-manager.io/docs/

#TargetNamespace: "cert-manager"

spec: components: HelmChartList: [
	#HelmChart & {
		metadata: name: "prod-mesh-certmanager"

		_dependsOn: "prod-secrets-namespaces": _
		namespace: #TargetNamespace
		_values: #Values & {
			installCRDs: true
			startupapicheck: enabled: false
			// Must not use kube-system on gke autopilot.  GKE Warden authz blocks access.
			global: leaderElection: namespace: #TargetNamespace
		}
		chart: {
			name:    "cert-manager"
			version: "1.14.3"
			repository: {
				name: "jetstack"
				url:  "https://charts.jetstack.io"
			}
		}
	},
]

// https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-resource-requests#min-max-requests
#PodResources: {
	requests: {
		cpu:                 string | *"250m"
		memory:              string | *"512Mi"
		"ephemeral-storage": string | *"100Mi"
	}
}

// https://cloud.google.com/kubernetes-engine/docs/how-to/autopilot-spot-pods
#NodeSelector: {
	// "kubernetes.io/os":          "linux"
	// "cloud.google.com/gke-spot": "true"
}
