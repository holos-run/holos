package v1alpha4

import (
	ks "sigs.k8s.io/kustomize/api/types"
	app "argoproj.io/application/v1alpha1"
)

#Platform: {
	Name:       string | *"no-platform-name"
	Components: _
	Resource: {
		metadata: name: Name
		spec: components: [for x in Components {x}]
	}
}

// https://holos.run/docs/api/author/v1alpha4/#Kubernetes
#Kubernetes: {
	Name:         _
	Component:    _
	Cluster:      _
	Resources:    _
	ArgoConfig:   _
	CommonLabels: _
	Namespace?:   _

	// Kustomize to add custom labels and manage the namespace.  More advanced
	// functionality than this should use the Core API directly and propose
	// extending the Author API if the need is common.
	_Transformer: {
		kind: "Kustomize"
		kustomize: kustomization: ks.#Kustomization & {
			commonLabels: "holos.run/component.name": BuildPlan.metadata.name
			commonLabels: CommonLabels
		}
	}

	_Artifacts: {
		component: {
			_path:    "clusters/\(Cluster)/components/\(Name)"
			artifact: "\(_path)/\(Name).gen.yaml"
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
				if Namespace != _|_ {
					kustomize: kustomization: namespace: Namespace
				}
			}]
		}

		if ArgoConfig.Enabled {
			argocd: {
				artifact: "clusters/\(Cluster)/gitops/\(Name).gen.yaml"
				let Output = "application.gen.yaml"
				generators: [{
					kind:   "Resources"
					output: Output
					resources: Application: (Name): app.#Application & {
						metadata: name:      Name
						metadata: namespace: string | *"argocd"
						spec: {
							destination: server: string | *"https://kubernetes.default.svc"
							project: string | *"default"
							source: {
								repoURL:        ArgoConfig.RepoURL
								path:           "\(ArgoConfig.Root)/\(component._path)"
								targetRevision: ArgoConfig.TargetRevision
							}
						}
					}
				}]
				transformers: [_Transformer & {
					inputs: [Output]
					output: artifact
					kustomize: kustomization: resources: inputs
				}]
			}
		}
	}

	BuildPlan: {
		metadata: name:  Name
		spec: component: Component
		spec: artifacts: [for x in _Artifacts {x}]
	}
}
