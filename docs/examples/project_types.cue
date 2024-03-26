package holos

import h "github.com/holos-run/holos/api/v1alpha1"

// #Projects is a map of all the projects in the platform.
#Projects: [Name=_]: #Project & {name: Name}

#Project: {
	name: string
	let ProjectName = name
	description: string
	environments: [Name=string]: #Environment & {
		name:    Name
		project: ProjectName
	}
	features: [Name=string]: #Feature & {name: Name}
}

#Environment: {
	name:    string
	project: string
	stage:   string | "dev" | "prod"
	slug:    "\(name)-\(project)"
}

#Feature: {
	name:        string
	description: string
	enabled:     *true | false
}

#ProjectTemplate: {
	project: #Project

	resources: [Name=_]: h.#KubernetesObjects & {
		metadata: name: Name
	}
}
