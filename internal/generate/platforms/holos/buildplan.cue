package holos

import "encoding/yaml"

import v1 "github.com/holos-run/holos/api/v1alpha1"

import corev1 "k8s.io/api/core/v1"

// #Helm represents a holos build plan composed of one helm chart.
#Helm: {
	// Name represents the holos component name
	Name:      string
	Version:   string
	Namespace: string

	Repo: {
		name: string | *""
		url:  string | *""
	}

	Values: {...}

	Chart: v1.#HelmChart & {
		metadata: name: string | *Name
		namespace: string | *Namespace
		chart: name:       string | *Name
		chart: version:    string | *Version
		chart: repository: Repo
		// Render the values to yaml for holos to provide to helm.
		valuesContent: yaml.Marshal(Values)
	}

	// output represents the build plan provided to the holos cli.
	Output: v1.#BuildPlan & {
		spec: components: helmChartList: [Chart]
	}
}

// #Kustomize represents a holos build plan composed of one kustomize build.
#Kustomize: {
	// Name represents the holos component name
	Name: string

	Kustomization: v1.#KustomizeBuild & {
		metadata: name: string | *Name
	}

	// output represents the build plan provided to the holos cli.
	Output: v1.#BuildPlan & {
		spec: components: kustomizeBuildList: [Kustomization]
	}
}

// #Kubernetes represents a holos build plan composed of inline kubernetes api
// objects.
#Kubernetes: {
	// Name represents the holos component name
	Name:      string
	Namespace: string

	Resources: [Kind=string]: [NAME=string]: {
		kind: Kind
		metadata: name: string | *NAME
	}

	Resources: Namespace: [string]: corev1.#Namespace

	// output represents the build plan provided to the holos cli.
	Output: v1.#BuildPlan & {
		// resources is a map unlike other build plans which use a list.
		spec: components: resources: "\(Name)": {
			metadata: name: Name
			apiObjectMap: (v1.#APIObjects & {apiObjects: Resources}).apiObjectMap
		}
	}
}
