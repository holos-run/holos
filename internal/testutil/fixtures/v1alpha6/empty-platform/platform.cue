package holos

import "github.com/holos-run/holos/api/author/v1alpha6:author"

// holos represents the field holos render platform evaluates, the resource
// field of the author.#Platform definition constructed from a components
// struct.
holos: platform.resource

platform: author.#Platform & {
	components: {
		// Empty platform with no components
	}
}