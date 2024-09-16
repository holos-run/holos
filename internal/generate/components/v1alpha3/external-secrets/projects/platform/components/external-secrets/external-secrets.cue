package holos

let Chart = {
	Name:      "external-secrets"
	Version:   "0.10.3"
	Namespace: "external-secrets"

	Repo: name: "external-secrets"
	Repo: url:  "https://charts.external-secrets.io"

	Values: installCRDs: false
}

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan
