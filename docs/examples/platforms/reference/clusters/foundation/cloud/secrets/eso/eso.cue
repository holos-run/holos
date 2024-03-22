package holos

// Manages the External Secrets Operator from the official upstream Helm chart.

#TargetNamespace: "external-secrets"
#Kustomization: spec: targetNamespace: #TargetNamespace

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-namespaces": _

		metadata: name: "prod-secrets-eso"
		namespace: #TargetNamespace
		chart: {
			name:    "external-secrets"
			version: "0.9.12"
			repository: {
				name: "external-secrets"
				url:  "https://charts.external-secrets.io"
			}
		}
		_values: installCrds: true
	},
]
