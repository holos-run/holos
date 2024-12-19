package holos

import "github.com/holos-run/holos/api/core/v1alpha5:core"

CertManager: {
	namespace: "cert-manager"
	datafile:  "./data/cert-manager/cert-manager.yaml"

	chart: core.#Chart & {
		name: "cert-manager"
		// See ./data/cert-manager/cert-manager.yaml
		version: string
		// repository config is used by multiple components.
		repository: {
			name: "jetstack"
			url:  "https://charts.jetstack.io"
		}
	}
}
