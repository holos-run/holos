package holos

import "encoding/yaml"

let Objects = {
	Name:      "{{ .Name }}"
	Namespace: "{{ .Namespace }}"

	Resources: {
		ConfigMap: {
			example: {
				metadata: namespace: "{{ .Namespace }}"
				// _Platform.Model represents the web form model
				data: platform: yaml.Marshal({model: _Platform.Model})
			}
		}
	}
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output
