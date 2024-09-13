package holos

import schema "github.com/holos-run/holos/api/schema/v1alpha3"

#Platform: schema.#Platform
#Fleets:   schema.#StandardFleets

_ComponentConfig: {
	Resources:  #Resources
	ArgoConfig: #ArgoConfig
}

#Helm:       schema.#Helm & _ComponentConfig
#Kustomize:  schema.#Kustomize & _ComponentConfig
#Kubernetes: schema.#Kubernetes & _ComponentConfig

#ArgoConfig: schema.#ArgoConfig & {
	ClusterName: _ClusterName
}
