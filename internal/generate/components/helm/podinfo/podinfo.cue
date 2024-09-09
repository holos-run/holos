package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "{{ .Name }}"
	Version:   "{{ .Version }}"
	Namespace: "{{ .Namespace }}"

	Repo: name: "{{ .RepoName }}"
	Repo: url:  "{{ .RepoURL }}"

	Values: {}
}
