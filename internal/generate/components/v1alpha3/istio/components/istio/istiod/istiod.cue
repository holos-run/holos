package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "istiod"
	Version:   #Istio.Version
	Namespace: #Istio.System.Namespace

	Chart: chart: name: "istiod"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #Istio.Values
}
