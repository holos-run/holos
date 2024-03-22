package holos

import "list"

#TargetNamespace: "default"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-namespaces":          _
		_dependsOn: "prod-secrets-eso-creds-refresher": _

		metadata: name: "prod-secrets-stores"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		for ns in #PlatformNamespaces {
			for obj in (#PlatformNamespaceObjects & {_ns: ns}).objects {
				let Kind = obj.kind
				let NS = ns.name
				let Name = obj.metadata.name
				"\(Kind)": "\(NS)/\(Name)": obj
			}
		}

		for nsName, ns in #ManagedNamespaces {
			if list.Contains(ns.clusterNames, #ClusterName) {
				let obj = #SecretStore & {_namespace: nsName}
				SecretStore: "\(nsName)/\(obj.metadata.name)": obj
			}
		}
	}
}

// #PlatformNamespaceObjects defines the api objects necessary for eso SecretStores in external clusters to access secrets in a given namespace in the provisioner cluster.
#PlatformNamespaceObjects: {
	_ns: #PlatformNamespace

	objects: [
		#SecretStore & {
			_namespace: _ns.name
		},
	]
}
