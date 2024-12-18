package holos

import holos "example.com/platform/schemas/holos/v1alpha5"

// Namespaces represents all managed namespaces across the platform.  Mix
// namespaces into this structure to manage them automatically from the
// namespaces component.
Namespaces: holos.#Namespaces

for PROJECT in Projects {
	Namespaces: PROJECT.namespaces
}
