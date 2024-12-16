package holos

// Produce a kubernetes objects build plan.
holos: Component.BuildPlan

Component: #Kubernetes & {
	Name: "httproutes"
	Resources: HTTPRoute: Projects[ProjectName].httpRoutes
}
