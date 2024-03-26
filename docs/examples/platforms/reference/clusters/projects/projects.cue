package holos

for ProjectName, Project in _Projects {
	spec: components: resources: (#ProjectTemplate & {project: Project}).resources
}
