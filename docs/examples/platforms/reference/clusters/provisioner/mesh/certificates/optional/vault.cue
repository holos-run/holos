package holos

let Vault = #OptionalServices.vault

if Vault.enabled {
	spec: components: KubernetesObjectsList: [
		#KubernetesObjects & {
			metadata: name: "prod-core-vault"

			apiObjectMap: OBJECTS.apiObjectMap
		},
	]
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		for k, obj in Vault.certs {
			"\(obj.kind)": "\(obj.metadata.name)": obj
		}
	}
}
