package holos

import api "github.com/holos-run/holos/api/author/v1alpha4"

// Manage a workload cluster named workload for use with the guides.
#Fleets: api.#StandardFleets & {
	workload: clusters: workload: _
}

// Define the default organization name.
#Organization: api.#OrganizationStrict & {
	DisplayName: string | *"Bank of Holos"
	Name:        string | *"bank-of-holos"
	Domain:      string | *"holos.localhost"
}

// https://holos.run/docs/api/author/v1alpha4/#ArgoConfig
#ArgoConfig: api.#ArgoConfig

#ComponentConfig: api.#ComponentConfig & {
	Name:       _Tags.name
	Component:  _Tags.component
	Cluster:    _Tags.cluster
	ArgoConfig: #ArgoConfig
	Resources:  #Resources
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
