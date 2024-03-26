package holos

_Projects: #Projects & {
	example: environments: {
		dev: stage:  "dev"
		jeff: stage: "dev"
		gary: stage: "dev"
		nate: stage: "dev"
	}
	iam: _
}

// Platform level definition of a project.
#Project: {
	// All projects have at least a prod environment and stage.
	environments: prod: stage: "prod"
}

#ProjectTemplate: {
	project: #Project

	resources: {
		// System namespace
		let SystemName = "\(project.name)-system"
		(SystemName): #KubernetesObjects & {
			metadata: name: SystemName
			apiObjectMap: (#APIObjects & {
				apiObjects: Namespace: (SystemName): #Namespace & {metadata: name: SystemName}
			}).apiObjectMap
		}

		// Project Namespaces
		let NamespacesName = "\(project.name)-namespaces"
		(NamespacesName): #KubernetesObjects & {
			metadata: name: NamespacesName
			apiObjectMap: (#APIObjects & {
				for _, env in project.environments {
					apiObjects: Namespace: (env.slug): #Namespace & {metadata: name: env.slug}
				}
			}).apiObjectMap
		}
	}
}
