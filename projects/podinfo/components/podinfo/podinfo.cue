package holos

// Parameters injected from the platform spec.
Parameters: {
	namespace:    string | *"podinfo-demo"                 @tag(NamespaceName)
	image:        string | *"ghcr.io/stefanprodan/podinfo" @tag(image)
	message:      string | *"Hello World"                  @tag(message)
	version:      string | *"6.7.0"                        @tag(version)
	replicaCount: int | *1                                 @tag(replicaCount, type=int)
}

// BuildPlan for holos to execute.
holos: Component.BuildPlan

// Configure the component from input parameters.
Component: #Helm & {
	Chart: {
		name:    "oci://ghcr.io/stefanprodan/charts/podinfo"
		release: "podinfo"
		version: Parameters.version
	}

	// Ensure all resources are located in the provided namespace
	KustomizeConfig: Kustomization: namespace: Parameters.namespace

	// The #Values definition is imported from the chart and defined in
	// values_schema.cue
	Values: #Values & {
		replicaCount: Parameters.replicaCount
		ui: message: Parameters.message
		image: {
			tag:        Chart.version
			repository: Parameters.image
		}
	}
}
