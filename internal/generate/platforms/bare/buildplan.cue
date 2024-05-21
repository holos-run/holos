package holos

import "encoding/yaml"

import v1 "github.com/holos-run/holos/api/v1alpha1"

// #Helm represents a holos build plan composed of one or more helm charts.
#Helm: {
	Name:      string
	Version:   string
	Namespace: string

	Repo: {
		name: string
		url:  string
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
