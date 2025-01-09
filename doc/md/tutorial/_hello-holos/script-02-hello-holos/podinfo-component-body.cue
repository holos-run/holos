package holos

// Produce a helm chart build plan.
holos: HelmChart.BuildPlan

HelmChart: #Helm & {
	Name: "podinfo"
	Chart: {
		version: "6.6.2"
		repository: {
			name: "podinfo"
			url:  "https://stefanprodan.github.io/podinfo"
		}
	}
	// Holos marshals Values into values.yaml for Helm.
	Values: {
		// message is a string with a default value.  @tag indicates a value may
		// be injected from the platform spec component parameters.
		ui: {
			message: string | *"Hello World" @tag(greeting, type=string)
		}
	}
}
