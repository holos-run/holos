package holos

_cluster: string @tag(cluster)
_project: string @tag(project)

holos: Component.BuildPlan

// Manage all app projects.
Component: #Kubernetes & {Resources: AppProject: AppProjects}
