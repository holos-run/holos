package holos

let Chart = {
	Name:      "external-secrets"
	Version:   "0.9.12"
	Namespace: "external-secrets"

	Repo: name: "external-secrets"
	Repo: url:  "https://charts.external-secrets.io"

	Values: installCrds: true
}

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan
