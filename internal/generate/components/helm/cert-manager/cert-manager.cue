package holos

let Chart = {
	Name:      "{{ .Name }}"
	Version:   "{{ .Version }}"
	Namespace: "{{ .Namespace }}"

	Repo: name: "{{ .RepoName }}"
	Repo: url:  "{{ .RepoURL }}"

	Values: {
		installCRDs: true
		startupapicheck: enabled: false
		// Must not use kube-system on gke autopilot.  GKE Warden blocks access.
		global: leaderElection: namespace: Namespace
	}
}

// Produce a helm chart build plan.
(#Helm & Chart).Output
