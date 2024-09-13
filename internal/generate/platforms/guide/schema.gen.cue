package holos

import schema "github.com/holos-run/holos/api/schema/v1alpha3"

#Helm: schema.#Helm & {
	ArgoConfig: #ArgoConfig
}

#Kustomize: schema.#Kustomize

#Kubernetes: schema.#Kubernetes

#ArgoConfig: schema.#ArgoConfig & {
	ClusterName: _ClusterName
}

#Fleets: schema.#StandardFleets

#Platform: schema.#Platform