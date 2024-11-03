package holos

import author "github.com/holos-run/holos/api/author/v1alpha5"

// Define the default organization name.
_Organization: author.#OrganizationStrict & {
	DisplayName: string | *"Bank of Holos"
	Name:        string | *"bank-of-holos"
	Domain:      string | *"holos.localhost"
}

// Projects represents a way to organize components into projects with owners.
// https://holos.run/docs/api/author/v1alpha5/#Projects
_Projects: author.#Projects

// ArgoConfig represents the configuration of ArgoCD Application resources for
// each component.
// https://holos.run/docs/api/author/v1alpha5/#ArgoConfig
_ArgoConfig: author.#ArgoConfig

#ComponentConfig: author.#ComponentConfig & {
	Name:    _Tags.component.name
	Path:    _Tags.component.path
	Cluster: _Tags.component.cluster
	ArgoConfig: _ArgoConfig & {
		if _Tags.project != "no-project" {
			AppProject: _Tags.project
		}
	}
	Resources: #Resources

	// Mix in project labels if the project is defined by the platform.
	if _Tags.project != "no-project" {
		CommonLabels: _Projects[_Tags.project].CommonLabels
	}
}

// https://holos.run/docs/api/author/v1alpha5/#Kubernetes
#Kubernetes: close({
	#ComponentConfig
	author.#Kubernetes
})

// https://holos.run/docs/api/author/v1alpha5/#Kustomize
#Kustomize: close({
	#ComponentConfig
	author.#Kustomize
})

// https://holos.run/docs/api/author/v1alpha5/#Helm
#Helm: close({
	#ComponentConfig
	author.#Helm
})
