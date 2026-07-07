package holos

import "github.com/holos-run/holos/api/core/v1alpha6:core"

// An OCI chart pulls directly from the registry.  The repository omits the url
// field and configures registry authentication with the auth field.
holos: core.#BuildPlan & {
	metadata: {
		name: "helm-oci"
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
				repository: auth: {
					username: fromEnv: "HOLOS_TEST_REGISTRY_USERNAME"
					password: fromEnv: "HOLOS_TEST_REGISTRY_PASSWORD"
				}
			}
		}]
	}]
}
