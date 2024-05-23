package holos

import corev1 "k8s.io/api/core/v1"

let Objects = {
	Name:      "namespaces"
	Namespace: "default"

	Resources: {
		for ns in #Namespaces {
			Namespace: "\(ns.name)": corev1.#Namespace
		}
	}
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output
