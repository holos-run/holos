package holos

import author "github.com/holos-run/holos/api/author/v1alpha5"

_Platform: author.#Platform & {
	Name: "default"
}

// Render a Platform resource for holos to process
_Platform.Resource
