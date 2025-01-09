package holos

// Produce a helm chart build plan.
holos: Helm.BuildPlan

Helm: #Helm & {
	Chart: {
		name:    "prometheus"
		version: "25.27.0"
		repository: {
			name: "prometheus-community"
			url:  "https://prometheus-community.github.io/helm-charts"
		}
	}
}
