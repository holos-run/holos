package holos

holos: Component.BuildPlan

Component: #Helm & {
	Chart: {
		name:    "oci://ghcr.io/stefanprodan/charts/podinfo"
		release: "podinfo"
		version: "6.6.1"
	}
	Values: {
		ui: {
			message: "Hello! Stage: \(StageName)"
		}
	}
}
