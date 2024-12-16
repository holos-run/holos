package holos

import "github.com/holos-run/holos/api/core/v1alpha5:core"

// #Project represents a collection of related components.
#Project: {
	name:       string
	components: #Components
	namespaces: #Namespaces
	httpRoutes: #HTTPRoutes
}

#Components: [string]: core.#Component

#Projects: [NAME=string]: #Project & {name: NAME}

#ProjectBuilder: {
	Name:       string | *"default"
	Components: #Components
	Stages:     #Stages

	Project: #Project & {
		name:       Name
		components: Components
	}
}
