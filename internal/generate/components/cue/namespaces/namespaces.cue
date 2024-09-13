package holos

let Objects = {
	Name:      "{{ .Name }}"
	Namespace: "{{ .Namespace }}"

	Resources: Namespace: _Namespaces
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan
