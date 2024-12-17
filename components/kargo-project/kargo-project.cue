package holos

// Produce a kubernetes objects build plan.
holos: Component.BuildPlan

_ProjectName: string @tag(ProjectName)

Component: #Kubernetes & {
	Resources: Project: (_ProjectName): {
		spec: promotionPolicies: KargoProjects[_ProjectName].promotionPolicies
	}
}
