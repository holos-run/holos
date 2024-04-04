package holos

spec: components: KustomizeBuildList: [
	#KustomizeBuild & {
		_dependsOn: "prod-secrets-stores": _

		metadata: name: "prod-iam-obs"
	},
]
