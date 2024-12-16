package holos

// Produce a kubernetes objects build plan.
holos: Component.BuildPlan

Component: #Kubernetes & {
	Resources: Namespace: Namespaces
}
