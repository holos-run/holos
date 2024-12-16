package holos

// Unify the components for each project into the platform.
// This works using the name label idiom, for example
//  things: [NAME=string]: #Thing & {name: NAME}
for PROJECT in Projects {
	for KEY, COMPONENT in PROJECT.components {
		Platform: Components: (KEY): COMPONENT & {
			parameters: ProjectName: PROJECT.name
		}
	}
}
