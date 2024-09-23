package holos

import api "github.com/holos-run/holos/api/author/v1alpha3"

// Define the default organization name
#Organization: DisplayName: string | *"Bank of Holos"
#Organization: Name:        string | *"bank-of-holos"

#Organization: api.#OrganizationStrict
#Platform:     api.#Platform
#Fleets:       api.#StandardFleets

_ComponentConfig: {
	Resources:  #Resources
	ArgoConfig: #ArgoConfig
}

#Helm:       api.#Helm & _ComponentConfig
#Kustomize:  api.#Kustomize & _ComponentConfig
#Kubernetes: api.#Kubernetes & _ComponentConfig

#ArgoConfig: api.#ArgoConfig & {
	ClusterName: _ClusterName
}
