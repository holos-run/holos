package holos

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-namespaces": _

		metadata: name: "prod-mesh-istio-base"
		namespace: "istio-system"
		chart: {
			name:    "base"
			version: "1.20.3"
			repository: {
				name: "istio"
				url:  "https://istio-release.storage.googleapis.com/charts"
			}
		}
		_values: #IstioValues
	},
]
