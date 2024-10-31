package holos

import api "github.com/holos-run/holos/api/author/v1alpha4"

// Define the default organization name.
_Organization: api.#OrganizationStrict & {
	DisplayName: string | *"Bank of Holos"
	Name:        string | *"bank-of-holos"
	Domain:      string | *"holos.localhost"
}

// Projects represents a way to organize components into projects with owners.
// https://holos.run/docs/api/author/v1alpha4/#Projects
_Projects: api.#Projects

// ArgoConfig represents the configuration of ArgoCD Application resources for
// each component.
// https://holos.run/docs/api/author/v1alpha4/#ArgoConfig
#ArgoConfig: api.#ArgoConfig

#ComponentConfig: api.#ComponentConfig & {
	Name:      _Tags.name
	Component: _Tags.component
	Cluster:   _Tags.cluster
	ArgoConfig: #ArgoConfig & {
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

// https://holos.run/docs/api/author/v1alpha4/#Kubernetes
#Kubernetes: close({
	#ComponentConfig
	api.#Kubernetes
})

// https://holos.run/docs/api/author/v1alpha4/#Kustomize
#Kustomize: close({
	#ComponentConfig
	api.#Kustomize
})

// https://holos.run/docs/api/author/v1alpha4/#Helm
#Helm: close({
	#ComponentConfig
	api.#Helm
})
