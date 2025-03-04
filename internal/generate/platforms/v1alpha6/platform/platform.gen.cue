package holos

import "github.com/holos-run/holos/api/author/v1alpha6:author"

Platform: author.#Platform & {
	Name: "default"
}

// Render a Platform resource for holos to process
holos: Platform.Resource
