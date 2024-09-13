package holos

let Objects = {
	Name: "namespaces"
	Resources: Namespace: #Namespaces
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan
