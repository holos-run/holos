package holos

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan

let Chart = {
	Name:      "{{ .Name }}"
	Version:   #CertManager.Version
	Namespace: #CertManager.Namespace

	Repo: name: "{{ .RepoName }}"
	Repo: url:  "{{ .RepoURL }}"

	Values: installCRDs: true
	Values: startupapicheck: enabled: false
}
