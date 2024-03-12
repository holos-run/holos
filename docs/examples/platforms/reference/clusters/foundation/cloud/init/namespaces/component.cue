package holos

#TargetNamespace: "default"

#InputKeys: {
	project:   "secrets"
	component: "namespaces"
}

#KubernetesObjects & {
	apiObjects: {
		// #ManagedNamespaces is the set of all namespaces across all clusters in the platform.
		for k, ns in #ManagedNamespaces {
			Namespace: "\(ns.name)": #Namespace & {metadata: ns}
		}

		// #PlatformNamespaces is deprecated in favor of #ManagedNamespaces.
		for ns in #PlatformNamespaces {
			Namespace: "\(ns.name)": #Namespace & {metadata: ns}
		}
	}
}
