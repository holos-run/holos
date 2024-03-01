package holos

#InputKeys: component: "istio-base"
#TargetNamespace: "istio-system"

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name:    "base"
		version: "1.20.3"
		repository: {
			name: "istio"
			url:  "https://istio-release.storage.googleapis.com/charts"
		}
	}
	values: #IstioValues
}
