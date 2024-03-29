package holos

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "prod-secrets-namespaces"
		apiObjectMap: (#APIObjects & {
			apiObjects: {
				// #ManagedNamespaces is the set of all namespaces across all clusters in the platform.
				for k, ns in #ManagedNamespaces {
					if ns.clusters[#ClusterName] != _|_ {
						Namespace: "\(k)": #Namespace & ns.namespace
					}
				}

				// #PlatformNamespaces is deprecated in favor of #ManagedNamespaces.
				for ns in #PlatformNamespaces {
					Namespace: "\(ns.name)": #Namespace & {metadata: ns}
				}
			}
		}).apiObjectMap
	},
]
