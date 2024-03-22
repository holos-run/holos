package holos

// Refer to https://github.com/CrunchyData/postgres-operator-examples/tree/main/kustomize/install/crd
spec: components: KustomizeBuildList: [
	#KustomizeBuild & {
		metadata: name: "prod-pgo-crds"
	},
]
