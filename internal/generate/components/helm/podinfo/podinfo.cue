package holos

import schema "github.com/holos-run/holos/api/schema/v1alpha3"

// Produce a helm chart build plan.
(schema.#Helm & Chart).Output

let Chart = {
	Name:      "{{ .Name }}"
	Version:   "{{ .Version }}"
	Namespace: "{{ .Namespace }}"

	Repo: name: "{{ .RepoName }}"
	Repo: url:  "{{ .RepoURL }}"

	Values: {}
}
