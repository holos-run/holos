package holos

// Produce a helm chart build plan.
holos: Helm.BuildPlan

Helm: #Helm & {
	Chart: {
		name:    "prometheus-blackbox-exporter"
		version: "9.0.1"
		repository: {
			name: "prometheus-community"
			url:  "https://prometheus-community.github.io/helm-charts"
		}
	}
}
