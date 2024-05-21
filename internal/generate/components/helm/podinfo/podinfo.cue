package holos

let Chart = {
	Name:      "{{ .Name }}"
	Version:   "{{ .Version }}"
	Namespace: "{{ .Namespace }}"

	Repo: name: "{{ .RepoName }}"
	Repo: url:  "{{ .RepoURL }}"

	Values: {}
}

// Produce a helm chart build plan.
(#Helm & Chart).Output
