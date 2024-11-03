package holos

import "github.com/holos-run/holos/api/core/v1alpha5:core"

// Note: tags should have a reasonable default value for cue export.
_Tags: {
	// Standardized parameters
	component: core.#Component & {
		name: string | *"no-name" @tag(holos_component_name, type=string)
		path: string | *"no-path" @tag(holos_component_path, type=string)
	}
}
