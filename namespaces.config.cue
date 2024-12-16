package holos

// Namespaces represents all managed namespaces across the platform.  Mix
// namespaces into this structure to manage them automatically from the
// namespaces component.
Namespaces: #Namespaces

for PROJECT in Projects {
	Namespaces: PROJECT.namespaces
}
