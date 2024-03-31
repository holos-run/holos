package holos

#HelmChart: {
	chart: {
		version: #IstioVersion
		repository: {
			name: "istio"
			url:  "https://istio-release.storage.googleapis.com/charts"
		}
	}
}
