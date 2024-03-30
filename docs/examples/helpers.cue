package holos

import "encoding/yaml"

// #APIObjects is the output type for api objects produced by cue.
#APIObjects: {
	// apiObjects holds each the api objects produced by cue.
	apiObjects: {
		[Kind=_]: {
			[string]: {
				kind: Kind
				...
			}
		}
		Namespace?: [Name=_]: #Namespace & {metadata: name: Name}
		ExternalSecret?: [Name=_]: #ExternalSecret & {_name: Name}
		VirtualService?: [Name=_]: #VirtualService & {metadata: name: Name}
		Issuer?: [Name=_]: #Issuer & {metadata: name: Name}
		Gateway?: [Name=_]: #Gateway & {metadata: name: Name}
		ConfigMap?: [Name=_]: #ConfigMap & {metadata: name: Name}
		Deployment?: [_]: #Deployment
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
		...
	}
}
