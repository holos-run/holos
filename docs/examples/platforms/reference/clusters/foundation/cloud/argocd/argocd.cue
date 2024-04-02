package holos

// This is a project managed namespace
#TargetNamespace: "prod-platform"

spec: components: KustomizeBuildList: [
	#KustomizeBuild & {
		_dependsOn: "prod-secrets-stores": _

		metadata: name: "argocd"
	},
]
