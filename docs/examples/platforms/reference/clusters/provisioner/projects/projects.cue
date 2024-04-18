package holos

for Project in _Projects {

	// Debugging variable to enable inspecting the project host data:
	// cue eval --out json -t cluster=provisioner ./platforms/reference/clusters/provisioner/projects/... -e _ProjectHosts.holos  > hosts.json
	let ProjectData = (#ProjectTemplate & {project: Project})
	_ProjectHosts: "\(Project.name)": ProjectData.ProjectHosts

	spec: components: resources: ProjectData.provisioner.resources
}
