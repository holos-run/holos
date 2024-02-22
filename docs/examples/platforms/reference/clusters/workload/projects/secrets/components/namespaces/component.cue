package holos

import "list"

#TargetNamespace: "default"

#InputKeys: {
	project:   "secrets"
	component: "namespaces"
}

// #PlatformNamespaceObjects defines the api objects necessary for eso SecretStores in external clusters to access secrets in a given namespace in the provisioner cluster.
#PlatformNamespaceObjects: {
	_ns: #PlatformNamespace

	objects: [
		#Namespace & {
			metadata: name: _ns.name
		},
	]
}

objects: list.FlattenN(_objects, 1)

_objects: [
	for ns in #PlatformNamespaces {
		(#PlatformNamespaceObjects & {_ns: ns}).objects
	},
]

{} & #KubernetesObjects
