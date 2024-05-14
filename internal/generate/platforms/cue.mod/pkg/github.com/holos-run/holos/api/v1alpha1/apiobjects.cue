package v1alpha1

import "encoding/yaml"

import core "k8s.io/api/core/v1"

// #APIObjects defines the output format for kubernetes api objects.  The holos
// cli expects the yaml representation of each api object in the apiObjectMap
// field.
#APIObjects: {
	// apiObjects represents the un-marshalled form of each kubernetes api object
	// managed by a holos component.
	apiObjects: {
		[Kind=string]: {
			[string]: {
				kind: Kind
				...
			}
		}
		ConfigMap: [string]: core.#ConfigMap & {apiVersion: "v1"}
	}

	// apiObjectMap holds the marshalled representation of apiObjects
	apiObjectMap: {
		for kind, v in apiObjects {
			"\(kind)": {
				for name, obj in v {
					"\(name)": yaml.Marshal(obj)
				}
			}
		}
	}
}
