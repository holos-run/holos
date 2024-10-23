package holos

import api "github.com/holos-run/holos/api/author/v1alpha4"

_Platform: api.#Platform & {
	Name: "default"
}

// Render a Platform resource for holos to process
_Platform.Resource
