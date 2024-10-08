package v1alpha3

// HelmChart represents a holos component which wraps around an upstream helm
// chart.  Holos orchestrates helm by providing values obtained from CUE,
// renders the output using `helm template`, then post-processes the helm output
// yaml using the general functionality provided by [Component], for
// example [Kustomize] post-rendering and mixing in additional kubernetes api
// objects.
type HelmChart struct {
	Component `json:",inline"`
	Kind      string `json:"kind" cue:"\"HelmChart\""`

	// Chart represents a helm chart to manage.
	Chart Chart `json:"chart"`
	// ValuesContent represents the values.yaml file holos passes to the `helm
	// template` command.
	ValuesContent string `json:"valuesContent"`
	// EnableHooks enables helm hooks when executing the `helm template` command.
	EnableHooks bool `json:"enableHooks" cue:"bool | *false"`
}

// Chart represents a helm chart.
type Chart struct {
	// Name represents the chart name.
	Name string `json:"name"`
	// Version represents the chart version.
	Version string `json:"version"`
	// Release represents the chart release when executing helm template.
	Release string `json:"release"`
	// Repository represents the repository to fetch the chart from.
	Repository Repository `json:"repository,omitempty"`
}

// Repository represents a helm chart repository.
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
