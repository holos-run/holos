package holos

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan

let Chart = {
	Name:      "istio-base"
	Version:   #Istio.Version
	Namespace: #Istio.System.Namespace

	Chart: chart: name: "base"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #Istio.Values
}
