package holos

#DependsOn: _ESOCreds

#TargetNamespace: "default"

#InputKeys: {
	project:   "secrets"
	component: "stores"
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

#KubernetesObjects & {
	apiObjects: {
		for ns in #PlatformNamespaces {
			for obj in (#PlatformNamespaceObjects & {_ns: ns}).objects {
				let Kind = obj.kind
				let NS = ns.name
				let Name = obj.metadata.name
				"\(Kind)": "\(NS)/\(Name)": obj
			}
		}

		for k, ns in #ManagedNamespaces {
			let obj = #SecretStore & {_namespace: ns.name}
			SecretStore: "\(ns.name)/\(obj.metadata.name)": obj
		}
	}
}
