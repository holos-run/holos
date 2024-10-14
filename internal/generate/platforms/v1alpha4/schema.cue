package holos

import (
	api "github.com/holos-run/holos/api/author/v1alpha4"
	core "github.com/holos-run/holos/api/core/v1alpha4"
	ks "sigs.k8s.io/kustomize/api/types"
	app "argoproj.io/application/v1alpha1"
)

// Manage a workload cluster named workload for use with the guides.
_Fleets: api.#StandardFleets & {
	workload: clusters: workload: _
}

// Define the default organization name.
_Organization: api.#OrganizationStrict & {
	DisplayName: string | *"Bank of Holos"
	Name:        string | *"bank-of-holos"
}

_ComponentConfig: {
	Resources:  #Resources
	ArgoConfig: api.#ArgoConfig

	// Kustomize all generators of all build steps to add common labels.
	_Transformer: {
		kind: "Kustomize"
		kustomize: kustomization: ks.#Kustomization & {
			commonLabels: "holos.run/component.name": BuildPlan.metadata.name
		}
	}

	// Tags injected from holos render platform into holos render component.
	BuildPlan: core.#BuildPlan & {
		metadata: name:  _Tags.name
		spec: component: _Tags.component
		spec: artifacts: [
			{
				artifact: "clusters/\(_Tags.cluster)/components/\(metadata.name)/\(metadata.name).gen.yaml"
				let Output = "resources.gen.yaml"
				generators: [{
					kind:      "Resources"
					output:    Output
					resources: Resources
				}]
				transformers: [_Transformer & {
					inputs: [Output]
					output: artifact
					kustomize: kustomization: resources: inputs
				}]
			},
			{
				artifact: "clusters/\(_Tags.cluster)/gitops/\(metadata.name).gen.yaml"
				let Output = "application.gen.yaml"
				generators: [{
					kind:   "Resources"
					output: Output
					resources: Application: argocd: app.#Application & {
						metadata: name:      BuildPlan.metadata.name
						metadata: namespace: "argocd"
						spec: {
							destination: server: "https://kubernetes.default.svc"
							project: "default"
							source: {
								path:           "deploy/clusters/\(_Tags.cluster)/components/\(metadata.name)"
								repoURL:        "https://github.com/holos-run/bank-of-holos"
								targetRevision: "main"
							}
						}
					}
				}]
				transformers: [_Transformer & {
					inputs: [Output]
					output: artifact
					kustomize: kustomization: resources: inputs
				}]
			},
		]
	}
}

#Kubernetes: api.#Kubernetes & _ComponentConfig
// #Helm:       api.#Helm & _ComponentConfig
// #Kustomize:  api.#Kustomize & _ComponentConfig

// #ArgoConfig: api.#ArgoConfig & {
// 	ClusterName: _ClusterName
// }
