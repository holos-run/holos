package holos

let Objects = {
	Name: "httproutes"
	Resources: HTTPRoute: #HTTPRoutes
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan
