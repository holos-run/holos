package holos

// https://cert-manager.io/docs/

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "cert-manager"
	Version:   "1.14.3"
	Namespace: "cert-manager"

	Repo: name: "jetstack"
	Repo: url:  "https://charts.jetstack.io"

	Values: {
		installCRDs: true
		startupapicheck: enabled: false
		// Must not use kube-system on gke autopilot.  GKE Warden blocks access.
		global: leaderElection: namespace: Namespace
	}
}
