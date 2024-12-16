package holos

import core "k8s.io/api/core/v1"

// #Namespace defines one namespace.
#Namespace: core.#Namespace & {
	metadata: name: string
	metadata: labels: "kubernetes.io/metadata.name": metadata.name
}

// #Namespaces defines the structure holding managed namespaces.
#Namespaces: {
	[NAME=string]: #Namespace & {
		metadata: name: NAME
	}
}
