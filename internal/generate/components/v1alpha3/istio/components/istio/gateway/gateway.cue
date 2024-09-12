package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "istio-gateway"
	Version:   #Istio.Version
	Namespace: #Istio.Gateway.Namespace

	Chart: chart: name: "gateway"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #Istio.Values
}
