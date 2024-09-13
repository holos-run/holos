package holos

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan

let Chart = {
	Name:      "istio-cni"
	Version:   #Istio.Version
	Namespace: #Istio.System.Namespace

	Chart: chart: name: "cni"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #Istio.Values
}
