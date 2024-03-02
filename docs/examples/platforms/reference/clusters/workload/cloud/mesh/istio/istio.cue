package holos

#DependsOn: _IstioBase

#HelmChart: {
	chart: {
		version: "1.20.3"
		repository: {
			name: "istio"
			url:  "https://istio-release.storage.googleapis.com/charts"
		}
	}
}
