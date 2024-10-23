package holos

import api "github.com/holos-run/holos/api/author/v1alpha4"

// Manage a workload cluster named local for use with the guides.
_Fleets: api.#StandardFleets & {
	workload: clusters: local: _
}
