package holos

import (
	h "github.com/holos-run/holos/api/v1alpha1"
	"encoding/yaml"
)

// CUE provides a #BuildPlan to the holos cli.  Holos requires the output of CUE
// to conform to the #BuildPlan schema.
{} & h.#BuildPlan

// #HolosComponent defines struct fields common to all holos component types.
#HolosComponent: {
	h.#HolosComponent
	metadata: name: string
	_NameLengthConstraint: len(metadata.name) & >=1
	...
}

#KubernetesObjects: #HolosComponent & h.#KubernetesObjects

// #HolosTypeMeta is similar to kubernetes api TypeMeta, but for holos api
// objects such as the Platform config resource.
#HolosTypeMeta: {
	kind:       string @go(Kind)
	apiVersion: string @go(APIVersion)
}

// #HolosObjectMeta is similar to kubernetes api ObjectMeta, but for holos api
// objects.
#HolosObjectMeta: {
	name: string @go(Name)
	labels: {[string]: string} @go(Labels,map[string]string)
	annotations: {[string]: string} @go(Annotations,map[string]string)
}

// #APIObjects defines the output format for kubernetes api objects.  The holos
// cli expects the yaml representation of each api object in the apiObjectMap
// field.
#APIObjects: {
	// apiObjects represents the un-marshalled form of each kubernetes api object
	// managed by a holos component.
	apiObjects: {
		[Kind=_]: {
			[string]: {
				kind: Kind
				...
			}
		}
		ConfigMap?: [Name=_]: #ConfigMap & {metadata: name: Name}
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
