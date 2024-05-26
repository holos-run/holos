package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "cert-manager"
	Version:   "1.14.5"
	Namespace: "cert-manager"

	Repo: name: "jetstack"
	Repo: url:  "https://charts.jetstack.io"

	Values: {
		installCRDs: true
		startupapicheck: enabled: false
		// Must not use kube-system on gke autopilot.  GKE Warden blocks access.
		global: leaderElection: namespace: Namespace

		// https://cloud.google.com/kubernetes-engine/docs/concepts/autopilot-resource-requests#min-max-requests
		resources: requests: {
			cpu:                 "250m"
			memory:              "512Mi"
			"ephemeral-storage": "100Mi"
		}
		webhook: resources:        Values.resources
		cainjector: resources:     Values.resources
		startupapicheck: resource: Values.resources

		// https://cloud.google.com/kubernetes-engine/docs/how-to/autopilot-spot-pods
		nodeSelector: {
			"kubernetes.io/os": "linux"
			if _ClusterName == "management" {
				"cloud.google.com/gke-spot": "true"
			}
		}
		webhook: nodeSelector:         Values.nodeSelector
		cainjector: nodeSelector:      Values.nodeSelector
		startupapicheck: nodeSelector: Values.nodeSelector
	}
}
