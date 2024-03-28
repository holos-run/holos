package holos

for Project in _Projects {
	spec: components: resources: (#ProjectTemplate & {project: Project}).workload.resources
}
