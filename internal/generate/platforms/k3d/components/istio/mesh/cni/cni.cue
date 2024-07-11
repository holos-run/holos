package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "istio-cni"
	Version:   #IstioVersion
	Namespace: "istio-system"

	Chart: chart: name: "cni"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #CNIValues
}
