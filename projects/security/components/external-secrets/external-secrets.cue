package holos

// Produce a helm chart build plan.
holos: Component.BuildPlan

Component: #Helm & {
	Name:      "external-secrets"
	Namespace: "external-secrets"

	Chart: {
		version: "0.10.3"
		repository: {
			name: "external-secrets"
			url:  "https://charts.external-secrets.io"
		}
	}

	Values: installCRDs: false
}
