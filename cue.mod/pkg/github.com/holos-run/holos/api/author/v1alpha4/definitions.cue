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

	KustomizeConfig: {
		Files:     _
		Resources: _
		Kustomization: ks.#Kustomization & {
			apiVersion: "kustomize.config.k8s.io/v1beta1"
			kind:       "Kustomization"
		}
	}

	// Kustomize to add custom labels and manage the namespace.  More advanced
	// functionality than this should use the Core API directly and propose
	// extending the Author API if the need is common.
	_TransformerArgo: core.#Transformer & {
		kind: "Kustomize"
		kustomize: kustomization: ks.#Kustomization & {
			commonLabels: "holos.run/component.name": BuildPlan.metadata.name
			commonLabels: CommonLabels
		}
	}

	// Add the argocd.argoproj.io/instance label to resources, but not to the
	// argocd Application config.
	_Transformer: _TransformerArgo & {
		kustomize: kustomization: commonLabels: {
			"argocd.argoproj.io/instance": Name
		}
	}

	_Artifacts: {
		component: {
			_path:    "clusters/\(Cluster)/components/\(Name)"
			artifact: "\(_path)/\(Name).gen.yaml"
			let ResourcesOutput = "resources.gen.yaml"
			let IntermediateOutput = "combined.gen.yaml"
			generators: [
				{
					kind:      "Resources"
					output:    ResourcesOutput
					resources: Resources
				},
				for x in KustomizeConfig.Files {
					kind:   "File"
					output: x.Source
					file: source: x.Source
				},
				for x in KustomizeConfig.Resources {
					kind:   "File"
					output: x.Source
					file: source: x.Source
				},
			]
			transformers: [
				core.#Transformer & {
					kind: "Kustomize"
					inputs: [for x in generators {x.output}]
					output: IntermediateOutput
					kustomize: kustomization: KustomizeConfig.Kustomization & {
						resources: [
							ResourcesOutput,
							for x in KustomizeConfig.Resources {x.Source},
						]
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

		// Mix in the ArgoCD Application gitops artifact.
		(#ArgoArtifact & {
			name:        Name
			cluster:     Cluster
			config:      ArgoConfig
			transformer: _TransformerArgo
			component:   _Artifacts.component._path
		}).Artifact
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

	KustomizeConfig: {
		Files:     _
		Resources: _
		Kustomization: ks.#Kustomization & {
			apiVersion: "kustomize.config.k8s.io/v1beta1"
			kind:       "Kustomization"
		}
	}

	// Kustomize to add custom labels and manage the namespace.  More advanced
	// functionality than this should use the Core API directly and propose
	// extending the Author API if the need is common.
	_TransformerArgo: core.#Transformer & {
		kind: "Kustomize"
		kustomize: kustomization: ks.#Kustomization & {
			commonLabels: "holos.run/component.name": BuildPlan.metadata.name
			commonLabels: CommonLabels
		}
	}

	// Add the argocd.argoproj.io/instance label to resources, but not to the
	// argocd Application config.
	_Transformer: _TransformerArgo & {
		kustomize: kustomization: commonLabels: {
			"argocd.argoproj.io/instance": Name
		}
	}

	_Artifacts: {
		component: {
			_path:    "clusters/\(Cluster)/components/\(Name)"
			artifact: "\(_path)/\(Name).gen.yaml"
			generators: [
				{
					kind:      "Resources"
					output:    "resources.gen.yaml"
					resources: Resources
				},
				for x in KustomizeConfig.Files {
					{
						kind: "File"
						file: source: x.Source
						output: file.source
					}
				},
			]
			let Intermediate = "intermediate.gen.yaml"
			transformers: [
				core.#Transformer & {
					kind: "Kustomize"
					inputs: [for x in generators {x.output}]
					output: Intermediate
					kustomize: kustomization: KustomizeConfig.Kustomization & {
						resources: [
							for x in inputs {x},
							for x in KustomizeConfig.Resources {x.Source},
						]
					}
				},
				_Transformer & {
					inputs: [Intermediate]
					output: artifact
					kustomize: kustomization: resources: inputs
					if Namespace != _|_ {
						kustomize: kustomization: namespace: Namespace
					}
				},
			]
		}

		// Mix in the ArgoCD Application gitops artifact.
		(#ArgoArtifact & {
			name:        Name
			cluster:     Cluster
			config:      ArgoConfig
			transformer: _TransformerArgo
			component:   _Artifacts.component._path
		}).Artifact
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

	Chart: {
		name:    string | *Name
		release: string | *name
	}
	Values:      _
	EnableHooks: true | *false

	KustomizeConfig: {
		Files:     _
		Resources: _
		Kustomization: ks.#Kustomization & {
			apiVersion: "kustomize.config.k8s.io/v1beta1"
			kind:       "Kustomization"
		}
	}

	// Kustomize to add custom labels and manage the namespace.  More advanced
	// functionality than this should use the Core API directly and propose
	// extending the Author API if the need is common.
	_TransformerArgo: core.#Transformer & {
		kind: "Kustomize"
		kustomize: kustomization: ks.#Kustomization & {
			commonLabels: "holos.run/component.name": BuildPlan.metadata.name
			commonLabels: CommonLabels
		}
	}

	// Add the argocd.argoproj.io/instance label to resources, but not to the
	// argocd Application config.
	_Transformer: _TransformerArgo & {
		kustomize: kustomization: commonLabels: {
			"argocd.argoproj.io/instance": Name
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
					kustomize: kustomization: KustomizeConfig.Kustomization & {
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

		// Mix in the ArgoCD Application gitops artifact.
		(#ArgoArtifact & {
			name:        Name
			cluster:     Cluster
			config:      ArgoConfig
			transformer: _TransformerArgo
			component:   _Artifacts.component._path
		}).Artifact
	}

	BuildPlan: {
		metadata: name:  Name
		spec: component: Component
		spec: artifacts: [for x in _Artifacts {x}]
	}
}

#ArgoArtifact: {
	name: string
	let Name = name
	cluster: string
	let Cluster = cluster
	config: #ArgoConfig
	let ArgoConfig = config
	transformer: core.#Transformer
	component:   string

	Artifact: {}
	if ArgoConfig.Enabled {
		Artifact: {
			argocd: core.#Artifact & {
				artifact: "clusters/\(Cluster)/gitops/\(Name).gen.yaml"
				generators: [{
					kind:   "Resources"
					output: "application.gen.yaml"
					resources: Application: (Name): app.#Application & {
						metadata: name:      Name
						metadata: namespace: string | *"argocd"
						spec: {
							destination: server: string | *"https://kubernetes.default.svc"
							project: ArgoConfig.AppProject
							source: {
								repoURL:        ArgoConfig.RepoURL
								path:           "\(ArgoConfig.Root)/\(component)"
								targetRevision: ArgoConfig.TargetRevision
							}
						}
					}
				}]
				transformers: [transformer & {
					inputs: [for x in generators {x.output}]
					output: artifact
					kustomize: kustomization: resources: inputs
				}]
			}
		}
	}
}
