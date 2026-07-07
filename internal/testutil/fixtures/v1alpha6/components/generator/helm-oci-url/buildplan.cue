package holos

import "github.com/holos-run/holos/api/core/v1alpha6:core"

// An OCI chart with a repository url must fail CUE validation.  OCI charts
// pull directly from the registry, so the url field must be omitted.  See the
// #Chart constraint in cue.mod/pkg.
holos: core.#BuildPlan & {
	metadata: {
		name: "helm-oci-url"
		labels: "holos.run/component.name":       name
		annotations: "app.holos.run/description": "\(name) generator"
	}
	spec: artifacts: [{
		artifact: "components/generator/\(metadata.name)/\(metadata.name).gen.yaml"
		generators: [{
			kind:   "Helm"
			output: artifact
			helm: chart: {
				name:    "oci://registry.example.com/charts/mychart"
				version: "0.1.0"
				release: "mychart"
				repository: url: "https://charts.example.com"
			}
		}]
	}]
}
