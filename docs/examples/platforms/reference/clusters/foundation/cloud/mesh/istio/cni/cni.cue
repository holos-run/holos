package holos

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-namespaces": _
		_dependsOn: "prod-mesh-istio-base":    _

		_values: #IstioValues
		metadata: name: "\(#InstancePrefix)-\(chart.name)"
		namespace: "kube-system"
		chart: name: "cni"
	},
]
