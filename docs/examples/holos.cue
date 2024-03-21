package holos

import (
	"encoding/yaml"
	h "github.com/holos-run/holos/api/v1alpha1"
	ksv1 "kustomize.toolkit.fluxcd.io/kustomization/v1"
)

// The overall structure of the data is:
// 1 CUE Instance => 1 BuildPlan => 0..N HolosComponents

// Holos requires CUE to evaluate and provide a valid BuildPlan.
// Constrain each CUE instance to output a BuildPlan.
{} & h.#BuildPlan

// #HolosComponent defines struct fields common to all holos component types.
#HolosComponent: {
	h.#HolosComponent
	metadata: name: string
	#namelen: len(metadata.name) & >=1
	let Name = metadata.name
	ksContent: yaml.Marshal(#Kustomization & {
		metadata: name: Name
	})
	...
}

// Holos component types.
#HelmChart:         #HolosComponent & h.#HelmChart
#KubernetesObjects: #HolosComponent & h.#KubernetesObjects
#KustomizeBuild:    #HolosComponent & h.#KustomizeBuild

// #ClusterName is the cluster name for cluster scoped resources.
#ClusterName: #InputKeys.cluster

// Flux Kustomization CRDs
#Kustomization: #NamespaceObject & ksv1.#Kustomization & {
	_dependsOn: [Name=_]: name: string & Name

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
