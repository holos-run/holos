package holos

let Chart = {
	Name:      "{{ .Name }}"
	Version:   "{{ .Version }}"
	Namespace: "{{ .Namespace }}"

	// OCI helm charts use the image url as the chart name
	Chart: chart: name: "{{ .Chart }}"

	Values: {}
}

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan
