package v1alpha4

import (
	ks "sigs.k8s.io/kustomize/api/types"
	app "argoproj.io/application/v1alpha1"
	core "github.com/holos-run/holos/api/core/v1alpha4"
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
	_Transformer: core.#Transformer & {
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

// https://holos.run/docs/api/author/v1alpha4/#Kustomize
#Kustomize: {
	Name:         _
	Component:    _
	Cluster:      _
	Resources:    _
	ArgoConfig:   _
	CommonLabels: _
	Namespace?:   _

	Kustomization: ks.#Kustomization & {
		apiVersion: "kustomize.config.k8s.io/v1beta1"
		kind:       "Kustomization"
	}

	// Transformer used as a generator.
	_Kustomization: core.#Transformer & {
		kind: "Kustomize"
		kustomize: kustomization: Kustomization
	}

	// Kustomize to add custom labels and manage the namespace.  More advanced
	// functionality than this should use the Core API directly and propose
	// extending the Author API if the need is common.
	_Transformer: core.#Transformer & {
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
			let Intermediate = "intermediate.gen.yaml"
			generators: [{
				kind:      "Resources"
				output:    Output
				resources: Resources
			}]
			transformers: [
				_Kustomization & {
					inputs: []
					output: Intermediate
				},
				_Transformer & {
					inputs: [Output, Intermediate]
					output: artifact
					kustomize: kustomization: resources: inputs
					if Namespace != _|_ {
						kustomize: kustomization: namespace: Namespace
					}
				},
			]
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

// https://holos.run/docs/api/author/v1alpha4/#Helm
#Helm: {
	Name:         _
	Component:    _
	Cluster:      _
	Resources:    _
	ArgoConfig:   _
	CommonLabels: _
	Namespace?:   _

	Chart: name: string | *Name
	Values:      _
	EnableHooks: true | *false

	Kustomization: ks.#Kustomization & {
		apiVersion: "kustomize.config.k8s.io/v1beta1"
		kind:       "Kustomization"
	}

	// Kustomize to add custom labels and manage the namespace.  More advanced
	// functionality than this should use the Core API directly and propose
	// extending the Author API if the need is common.
	_Transformer: core.#Transformer & {
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
			let HelmOutput = "helm.gen.yaml"
			let ResourcesOutput = "resources.gen.yaml"
			let IntermediateOutput = "combined.gen.yaml"
			generators: [
				{
					kind:   "Helm"
					output: HelmOutput
					helm: core.#Helm & {
						chart:       Chart
						values:      Values
						enableHooks: EnableHooks
						if Namespace != _|_ {
							namespace: Namespace
						}
					}
				},
				{
					kind:      "Resources"
					output:    ResourcesOutput
					resources: Resources
				},
			]
			transformers: [
				core.#Transformer & {
					kind: "Kustomize"
					inputs: [HelmOutput, ResourcesOutput]
					output: IntermediateOutput
					kustomize: kustomization: Kustomization & {
						resources: inputs
					}
				},
				_Transformer & {
					inputs: [IntermediateOutput]
					output: artifact
					kustomize: kustomization: resources: inputs
					if Namespace != _|_ {
						kustomize: kustomization: namespace: Namespace
					}
				},
			]
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
