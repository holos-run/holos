package holos

// Unify the components for each project into the platform.
// This works using the name label idiom, for example
//  things: [NAME=string]: #Thing & {name: NAME}
for PROJECT in Projects {
	Platform: Components: PROJECT.components
}
