package holos

// Produce a kubernetes objects build plan.
holos: Component.BuildPlan

Component: #Kubernetes & {
	Resources: Project: (ProjectName): {
		spec: promotionPolicies: KargoProjects[ProjectName].promotionPolicies
	}
}
