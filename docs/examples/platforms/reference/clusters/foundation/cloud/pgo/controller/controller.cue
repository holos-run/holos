package holos

spec: components: KustomizeBuildList: [
	#KustomizeBuild & {
		_dependsOn: "prod-secrets-namespaces": _
		_dependsOn: "prod-pgo-crds":           _

		metadata: name: "prod-pgo-controller"
	},
]
