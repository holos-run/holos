package holos

import (
	"encoding/json"

	"github.com/holos-run/holos/api/core/v1alpha6:core"
	"github.com/holos-run/holos/api/author/v1alpha6:author"
)

#ComponentConfig: author.#ComponentConfig & {
	Name:      _Tags.component.name
	Path:      _Tags.component.path
	Resources: #Resources

	// Inject the output base directory from platform component parameters.
	OutputBaseDir: string | *"" @tag(outputBaseDir, type=string)

	// labels is an optional field, guard references to it.
	if _Tags.component.labels != _|_ {
		Labels: _Tags.component.labels
	}

	// annotations is an optional field, guard references to it.
	if _Tags.component.annotations != _|_ {
		Annotations: _Tags.component.annotations
	}
}

// https://holos.run/docs/api/author/v1alpha6/#Kubernetes
#Kubernetes: close({
	#ComponentConfig
	author.#Kubernetes
})

// https://holos.run/docs/api/author/v1alpha6/#Kustomize
#Kustomize: close({
	#ComponentConfig
	author.#Kustomize
})

// https://holos.run/docs/api/author/v1alpha6/#Helm
#Helm: close({
	#ComponentConfig
	author.#Helm
})

// Note: tags should have a reasonable default value for cue export.
_Tags: {
	// Standardized parameters
	component: core.#Component & {
		name: string | *"no-name" @tag(holos_component_name, type=string)
		path: string | *"no-path" @tag(holos_component_path, type=string)

		_labels_json: string | *"" @tag(holos_component_labels, type=string)
		_labels: {}
		if _labels_json != "" {
			_labels: json.Unmarshal(_labels_json)
		}
		for k, v in _labels {
			labels: (k): v
		}

		_annotations_json: string | *"" @tag(holos_component_annotations, type=string)
		_annotations: {}
		if _annotations_json != "" {
			_annotations: json.Unmarshal(_annotations_json)
		}
		for k, v in _annotations {
			annotations: (k): v
		}
	}
}
