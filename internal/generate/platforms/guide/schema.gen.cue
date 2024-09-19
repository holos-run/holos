package holos

import api "github.com/holos-run/holos/api/author/v1alpha3"

#Platform: api.#Platform
#Fleets:   api.#StandardFleets

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
