package holos

let Objects = {
	Name:      "namespaces"
	Namespace: "default"

	Resources: Namespace: _Namespaces
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output
