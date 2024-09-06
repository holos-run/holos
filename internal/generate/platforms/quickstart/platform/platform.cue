package holos

import core "github.com/holos-run/holos/api/core/v1alpha3"

core.#Platform & {
	metadata: name: "quickstart"

	spec: model: {}
	spec: components: []
}
