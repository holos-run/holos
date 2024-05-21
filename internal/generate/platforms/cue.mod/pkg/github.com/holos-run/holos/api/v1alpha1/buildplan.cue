package v1alpha1

#BuildPlan: {
	apiVersion: #APIVersion
	kind:       #BuildPlanKind
}

#HelmChart: {
	apiVersion: #APIVersion
	kind:       "HelmChart"

	metadata: name: string
	chart: name:    string | *metadata.name
	chart: release: string | *metadata.name
}
