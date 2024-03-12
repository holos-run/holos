package holos

let Vault = #OptionalServices.vault

if Vault.enabled {
	#KubernetesObjects & {
		apiObjects: {
			for k, obj in Vault.certs {
				"\(obj.kind)": "\(obj.metadata.name)": obj
			}
		}
	}
}
