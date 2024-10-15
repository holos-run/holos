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

let ComponentConfig = {
	Name:       _Tags.name
	Component:  _Tags.component
	Cluster:    _Tags.cluster
	ArgoConfig: #ArgoConfig
	Resources:  #Resources
}

// https://holos.run/docs/api/author/v1alpha4/#Kubernetes
#Kubernetes: api.#Kubernetes & ComponentConfig

// https://holos.run/docs/api/author/v1alpha4/#Kustomize
#Kustomize: api.#Kustomize & ComponentConfig

// https://holos.run/docs/api/author/v1alpha4/#Helm
#Helm: api.#Helm & ComponentConfig
