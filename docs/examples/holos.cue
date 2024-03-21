package holos

import (
	"encoding/yaml"
	h "github.com/holos-run/holos/api/v1alpha1"
	kc "sigs.k8s.io/kustomize/api/types"
	ksv1 "kustomize.toolkit.fluxcd.io/kustomization/v1"
)

// The overall structure of the data is:
// 1 CUE Instance => 1 BuildPlan => 0..N HolosComponents

// Holos requires CUE to evaluate and provide a valid BuildPlan.
// Constrain each CUE instance to output a BuildPlan.
{} & h.#BuildPlan

let DependsOn = {[Name=_]: name: string & Name}

// #HolosComponent defines struct fields common to all holos component types.
#HolosComponent: {
	h.#HolosComponent
	_dependsOn: DependsOn
	let DEPENDS_ON = _dependsOn
	metadata: name: string
	#namelen: len(metadata.name) & >=1
	let Name = metadata.name
	ksContent: yaml.Marshal(#Kustomization & {
		_dependsOn: DEPENDS_ON
		metadata: name: Name
	})
	// Leave the HolosComponent open for components with additional fields like HelmChart.
	// Refer to https://cuelang.org/docs/tour/types/closed/
	...
}

//#KustomizeFiles represents resources for holos to write into files for kustomize post-processing.
#KustomizeFiles: {
	// Objects collects files for Holos to write for kustomize post-processing.
	Objects: "kustomization.yaml": #Kustomize
	// Files holds the marshaled output of Objects holos writes to the filesystem before calling the kustomize post-processor.
	Files: {
		for filename, obj in Objects {
			"\(filename)": yaml.Marshal(obj)
		}
	}
}

// Holos component types.
#HelmChart: #HolosComponent & h.#HelmChart & {
	_values:         _
	_kustomizeFiles: #KustomizeFiles

	// Render the values to yaml for holos to provide to helm.
	valuesContent: yaml.Marshal(_values)
	// Kustomize post-processor
	// resources is the intermediate file name for api objects.
	resourcesFile: h.#ResourcesFile
	// kustomizeFiles represents the files in a kustomize directory tree.
	kustomizeFiles: _kustomizeFiles.Files

	chart: h.#Chart & {
		name:    string
		release: string | *name
	}
}
#KubernetesObjects: #HolosComponent & h.#KubernetesObjects
#KustomizeBuild:    #HolosComponent & h.#KustomizeBuild

// #ClusterName is the cluster name for cluster scoped resources.
#ClusterName: #InputKeys.cluster

// Flux Kustomization CRDs
#Kustomization: #NamespaceObject & ksv1.#Kustomization & {
	_dependsOn: DependsOn

	metadata: {
		name:      string
		namespace: string | *"flux-system"
	}
	spec: ksv1.#KustomizationSpec & {
		interval:      string | *"30m0s"
		path:          string | *"deploy/clusters/\(#InputKeys.cluster)/components/\(metadata.name)"
		prune:         bool | *true
		retryInterval: string | *"2m0s"
		sourceRef: {
			kind: string | *"GitRepository"
			name: string | *"flux-system"
		}
		suspend?:         bool
		targetNamespace?: string
		timeout:          string | *"3m0s"
		// wait performs health checks for all reconciled resources. If set to true, .spec.healthChecks is ignored.
		// Setting this to true for all components generates considerable load on the api server from watches.
		// Operations are additionally more complicated when all resources are watched.  Consider setting wait true for
		// relatively simple components, otherwise target specific resources with spec.healthChecks.
		wait: true | *false
		dependsOn: [for k, v in _dependsOn {v}, ...]
	}
}

// #Kustomize represents the kustomize post processor.
#Kustomize: kc.#Kustomization & {
	_patches: {[_]: kc.#Patch}
	apiVersion: "kustomize.config.k8s.io/v1beta1"
	kind:       "Kustomization"
	// resources are file names holos will use to store intermediate component output for kustomize to post-process (i.e. helm template | kubectl kustomize)
	// See the related resourcesFile field of the holos component.
	resources: [h.#ResourcesFile]
	if len(_patches) > 0 {
		patches: [for v in _patches {v}]
	}
}
