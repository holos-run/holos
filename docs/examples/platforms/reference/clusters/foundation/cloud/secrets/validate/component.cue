package holos

// Validate ESO by syncing a secret with a SecretStore.

#TargetNamespace: "holos-system"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-stores": _

		metadata: name: "prod-secrets-validate"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: validate: #ExternalSecret & {
			_name: "validate"
		}
	}
}
