package v1alpha3

import "encoding/yaml"

// #APIObjects defines the output format for kubernetes api objects.  The holos
// cli expects the yaml representation of each api object in the apiObjectMap
// field.
#APIObjects: {
	apiObjects: {...}

	for kind, v in apiObjects {
		for name, obj in v {
			apiObjectMap: (kind): (name): yaml.Marshal(obj)
		}
	}
}
