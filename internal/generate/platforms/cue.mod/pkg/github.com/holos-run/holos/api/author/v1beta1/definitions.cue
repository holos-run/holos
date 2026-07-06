// Hand-written assembly CUE for the author v1beta1 package.  The generated
// definitions in cue.mod/gen define types only; this file fills
// TaskSet.spec.tasks from the author wrapper fields per the design doc
// doc/design/v1beta1/schema.md.  Mirrors the v1alpha6 sibling, replacing the
// BuildPlan artifact/generator/transformer lists with struct-keyed tasks.
package author

import (
	"strings"

	ks "sigs.k8s.io/kustomize/api/types"
	core "github.com/holos-run/holos/api/core/v1beta1:core"
)

#Platform: {
	name:       _
	components: _
	resource: {
		metadata: "name": name
		spec: "components": [for x in components {x}]
	}
}

#KustomizeConfig: {
	CommonLabels: _
	Kustomization: ks.#Kustomization & {
		apiVersion: "kustomize.config.k8s.io/v1beta1"
		kind:       "Kustomization"
		_labels: {}
		if len(CommonLabels) > 0 {
			_labels: commonLabels: {
				includeSelectors: false
				pairs:            CommonLabels
			}
			labels: [for x in _labels {x}]
		}
	}
}

// _TaskKey validates a task name is an RFC 1123 label per
// doc/design/v1beta1/schema.md#d3-task-naming-and-namespacing: lowercase
// alphanumerics and interior hyphens, at most 63 characters.  The out field
// is defined only when IN is valid, so an invalid task name fails evaluation
// immediately (undefined field: out) instead of exporting.
_TaskKey: {
	IN: string
	if IN =~ "^[a-z0-9]([a-z0-9-]*[a-z0-9])?$" && len(IN) <= 63 {
		out: IN
	}
}

// _TaskName converts a component directory file path into a valid task name
// validated by _TaskKey.  A source path this conversion cannot make valid
// fails evaluation immediately.  Use the ComponentConfig Tasks field directly
// for such sources.
_TaskName: {
	IN: string
	out: (_TaskKey & {
		"IN": strings.ToLower(strings.Replace(strings.Replace(strings.Replace(IN, "/", "-", -1), ".", "-", -1), "_", "-", -1))
	}).out
}

// Kustomize and Kubernetes are identical.

// https://holos.run/docs/api/author/v1beta1/#Kustomize
#Kustomize: #Kubernetes

// https://holos.run/docs/api/author/v1beta1/#Kubernetes
#Kubernetes: {
	Name:            _
	Resources:       _
	KustomizeConfig: _

	TaskSet: spec: tasks: {
		let ResourcesOutput = "resources.gen.yaml"
		resources: {
			kind:        "Resources"
			output:      ResourcesOutput
			"resources": Resources
		}
		for x in KustomizeConfig.Files {
			((_TaskName & {IN: "file-\(x.Source)"}).out): {
				kind:   "File"
				output: x.Source
				file: source: x.Source
			}
		}
		kustomize: {
			kind: "Kustomize"
			inputs: [
				ResourcesOutput,
				for x in KustomizeConfig.Files {x.Source},
			]
			output: "\(Name).gen.yaml"
			"kustomize": kustomization: KustomizeConfig.Kustomization & {
				"resources": [
					ResourcesOutput,
					for x in KustomizeConfig.Files {x.Source},
					for x in KustomizeConfig.Resources {x.Source},
				]
			}
		}
	}
}

// https://holos.run/docs/api/author/v1beta1/#Helm
#Helm: {
	Name:            _
	Resources:       _
	KustomizeConfig: _

	Chart: {
		name:    string | *Name
		release: string | *name
	}
	Values:       _
	ValueFiles?:  _
	EnableHooks:  _
	Namespace?:   _
	APIVersions?: _
	KubeVersion?: _

	TaskSet: spec: tasks: {
		let HelmOutput = "helm.gen.yaml"
		let ResourcesOutput = "resources.gen.yaml"
		helm: {
			kind:   "Helm"
			output: HelmOutput
			"helm": core.#Helm & {
				chart:  Chart
				values: Values
				if ValueFiles != _|_ {
					valueFiles: ValueFiles
				}
				enableHooks: EnableHooks
				if Namespace != _|_ {
					namespace: Namespace
				}
				if APIVersions != _|_ {
					apiVersions: APIVersions
				}
				if KubeVersion != _|_ {
					kubeVersion: KubeVersion
				}
			}
		}
		resources: {
			kind:        "Resources"
			output:      ResourcesOutput
			"resources": Resources
		}
		for x in KustomizeConfig.Files {
			((_TaskName & {IN: "file-\(x.Source)"}).out): {
				kind:   "File"
				output: x.Source
				file: source: x.Source
			}
		}
		kustomize: {
			kind: "Kustomize"
			inputs: [
				HelmOutput,
				ResourcesOutput,
				for x in KustomizeConfig.Files {x.Source},
			]
			output: "\(Name).gen.yaml"
			"kustomize": kustomization: KustomizeConfig.Kustomization & {
				"resources": [
					HelmOutput,
					ResourcesOutput,
					for x in KustomizeConfig.Files {x.Source},
					for x in KustomizeConfig.Resources {x.Source},
				]
			}
		}
	}
}

#ComponentConfig: {
	Name:          _
	Labels:        _
	Annotations:   _
	OutputBaseDir: _
	Tasks:         _

	_outputPath: string
	if OutputBaseDir == "" {
		_outputPath: "components/\(Name)"
	}
	if OutputBaseDir != "" {
		_outputPath: "\(OutputBaseDir)/components/\(Name)"
	}

	// TaskSet represents the derived TaskSet produced for the holos render
	// component command.  Tasks mix in by unification; the deploy sink writes
	// the final artifact per doc/design/v1beta1/schema.md#d2-artifact-writing.
	TaskSet: core.#TaskSet & {
		metadata: "name": Name
		if len(Labels) != 0 {
			metadata: labels: Labels
		}
		if len(Annotations) != 0 {
			metadata: annotations: Annotations
		}

		// Mix in Tasks, validating each task name per schema.md D3.
		spec: "tasks": {
			for k, v in Tasks {
				((_TaskKey & {IN: k}).out): v
			}
		}
		spec: "tasks": deploy: {
			kind: "Artifact"
			inputs: ["\(Name).gen.yaml"]
			artifact: path: "\(_outputPath)/\(Name).gen.yaml"
		}
	}
}
