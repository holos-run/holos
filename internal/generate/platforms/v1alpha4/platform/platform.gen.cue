package holos

import author "github.com/holos-run/holos/api/author/v1alpha4"

_Platform: author.#Platform & {Name: "guide"}

// Render a Platform for holos to process
_Platform.Resource
