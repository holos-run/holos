package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "istio-base"
	Version:   #IstioVersion
	Namespace: "istio-system"

	Chart: chart: name: "base"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #IstioValues
}
