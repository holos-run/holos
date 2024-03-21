package holos

#TargetNamespace: #ARCSystemNamespace
#InputKeys: component: "arc-system"

spec: components: HelmCharts: [
	#HelmChart & {
		metadata: name: "prod-github-arc-system"

		_dependsOn: "prod-secrets-namespaces": _
		_values:   #Values & #DefaultSecurityContext
		namespace: #TargetNamespace
		chart: {
			// Match the gha-base-name in the chart _helpers.tpl to avoid long full names.
			release: "gha-rs-controller"
			name:    "oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller"
			version: "0.8.3"
		}
	},
]
