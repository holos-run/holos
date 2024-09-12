package holos

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "istio-ztunnel"
	Version:   #Istio.Version
	Namespace: #Istio.System.Namespace

	Chart: chart: name: "ztunnel"

	Repo: name: "istio"
	Repo: url:  "https://istio-release.storage.googleapis.com/charts"

	Values: #Istio.Values
}
