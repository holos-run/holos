package holos

let Objects = {
	Name: "app-projects"
	Resources: AppProject: #AppProjects
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan
