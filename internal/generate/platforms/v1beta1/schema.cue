package holos

import "strings"

// The v1beta1 runtime is intentionally local to an initialized platform. It
// models only Holos task assembly; Kubernetes schemas remain an opt-in choice
// for component authors instead of being vendored into every platform.
// #KustomizeConfig describes the optional local transformation of generated
// resources. BasePath keeps existing Kustomize bases first-class.
#KustomizeConfig: {
	BasePath:        string | *""
	LoadRestrictor?: "LoadRestrictionsNone"
	Files: {[string]: _} | *{}
	Resources: {[string]: {Source: string}} | *{}
	CommonLabels: {[string]: string} | *{}
	Kustomization: {...} | *{}
}

_TaskKey: {
	IN: string
	if IN =~ "^[a-z0-9]([a-z0-9-]*[a-z0-9])?$" && len(IN) <= 63 {
		out: IN
	}
}

_TaskName: {
	IN: string
	out: (_TaskKey & {
		"IN": strings.ToLower(strings.Replace(strings.Replace(strings.Replace(IN, "/", "-", -1), ".", "-", -1), "_", "-", -1))
	}).out
}

// #Kubernetes produces raw Kubernetes resources and optionally transforms
// them with Kustomize.
#Kubernetes: close({
	Name:            _Tags.component.name
	Path:            _Tags.component.path
	Resources:       #Resources
	KustomizeConfig: #KustomizeConfig
	Tasks: {[string]: _} | *{}
	let componentName = Name
	let componentResources = Resources
	let componentTasks = Tasks
	let kustomizeConfig = KustomizeConfig
	TaskSet: {
		metadata: {
			name: componentName
			if _Tags.component.labels != _|_ {
				labels: _Tags.component.labels
			}
			if _Tags.component.annotations != _|_ {
				annotations: _Tags.component.annotations
			}
		}
		spec: tasks: {
			for taskName, task in componentTasks {
				((_TaskKey & {IN: taskName}).out): task
			}
			resources: {
				kind:        "Resources"
				output:      "resources.gen.yaml"
				"resources": componentResources
			}
			for source, _ in kustomizeConfig.Files {
				let fileSource = source
				((_TaskName & {IN: "file-\(fileSource)"}).out): {
					kind:   "File"
					output: source
					file: source: fileSource
				}
			}
			kustomize: {
				kind: "Kustomize"
				if kustomizeConfig.BasePath == "" {
					inputs: ["resources.gen.yaml", for source, _ in kustomizeConfig.Files {source}]
				}
				output: "\(componentName).gen.yaml"
				"kustomize": {
					if kustomizeConfig.BasePath != "" {
						basePath: kustomizeConfig.BasePath
						if kustomizeConfig.LoadRestrictor != _|_ {
							loadRestrictor: kustomizeConfig.LoadRestrictor
						}
					}
					if kustomizeConfig.BasePath == "" {
						kustomization: kustomizeConfig.Kustomization & {
							resources: ["resources.gen.yaml", for source, _ in kustomizeConfig.Files {source}, for _, resource in kustomizeConfig.Resources {resource.Source}]
							if len(kustomizeConfig.CommonLabels) > 0 {
								labels: [{includeSelectors: false, pairs: kustomizeConfig.CommonLabels}]
							}
						}
					}
				}
			}
			deploy: {
				kind: "Artifact"
				inputs: ["\(componentName).gen.yaml"]
				artifact: path: "components/\(componentName)/\(componentName).gen.yaml"
			}
		}
	}
})

// #Kustomize has the same task shape as #Kubernetes.
#Kustomize: close({
	#Kubernetes
})

// #Helm renders a chart and combines it with generated resources through
// Kustomize. No chart-specific wrapper is required.
#Helm: close({
	Name:            _Tags.component.name
	Path:            _Tags.component.path
	Resources:       #Resources
	KustomizeConfig: #KustomizeConfig
	Tasks: {[string]: _} | *{}
	let componentName = Name
	let componentResources = Resources
	let componentTasks = Tasks
	let componentKustomizeConfig = KustomizeConfig
	Chart: {
		name:    string | *componentName
		version: string
		release: string | *name
		repository?: {
			name?: string
			url?:  string
		}
	}
	Values: {...} | *{}
	ValueFiles?: [...{name: string, kind: "Values", values?: {...}}]
	EnableHooks: bool | *false
	Namespace?:  string
	APIVersions?: [...string]
	KubeVersion?: string
	TaskSet: {
		metadata: name: componentName
		spec: tasks: {
			for taskName, task in componentTasks {
				((_TaskKey & {IN: taskName}).out): task
			}
			helm: {
				kind:   "Helm"
				output: "helm.gen.yaml"
				"helm": {
					chart:        Chart
					values:       Values
					valueFiles?:  ValueFiles
					enableHooks:  EnableHooks
					namespace?:   Namespace
					apiVersions?: APIVersions
					kubeVersion?: KubeVersion
				}
			}
			resources: {
				kind:        "Resources"
				output:      "resources.gen.yaml"
				"resources": componentResources
			}
			kustomize: {
				kind: "Kustomize"
				inputs: ["helm.gen.yaml", "resources.gen.yaml"]
				output: "\(componentName).gen.yaml"
				"kustomize": kustomization: componentKustomizeConfig.Kustomization & {
					resources: ["helm.gen.yaml", "resources.gen.yaml"]
				}
			}
			deploy: {
				kind: "Artifact"
				inputs: ["\(componentName).gen.yaml"]
				artifact: path: "components/\(componentName)/\(componentName).gen.yaml"
			}
		}
	}
})
