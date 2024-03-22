package holos

// Certificates usually go into the istio-system namespace, but they may go anywhere.
#TargetNamespace: "default"

#Kustomization: spec: wait: true

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "\(#InstancePrefix)-certificates"

		_dependsOn: "prod-secrets-namespaces": _
		_dependsOn: "prod-mesh-letsencrypt":   _

		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let Vault = #OptionalServices.vault

let OBJECTS = #APIObjects & {
	apiObjects: {
		for k, obj in #PlatformCerts {
			"\(obj.kind)": {
				"\(obj.metadata.namespace)/\(obj.metadata.name)": obj
			}
		}

		if Vault.enabled {
			for k, obj in Vault.certs {
				"\(obj.kind)": "\(obj.metadata.name)": obj
			}
		}
	}
}
