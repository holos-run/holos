package holos

import "encoding/yaml"

// Mix in a Kustomize patch to the configuration.
Kustomize: KustomizeConfig: Kustomization: _patches: {
	probe: {
		target: kind: "Service"
		target: name: "httpbin"
		patch: yaml.Marshal([{
			op:    "add"
			path:  "/metadata/annotations/prometheus.io~1probe"
			value: "true"
		}])
	}
}
