@extern(embed)
package certmanager

import "github.com/holos-run/holos/api/core/v1alpha5:core"

// Unify data from yaml for Kargo integration.
_data: _ @embed(file=cert-manager.yaml)

config: #Config & {
	namespace: "cert-manager"
	// datafile value must align to the embed file directive above for proper
	// configuration of Kargo promotion stages.
	datafile: "./config/certmanager/cert-manager.yaml"
	chart: {
		name:    "cert-manager"
		version: _data.chart.version
		repository: {
			name: "jetstack"
			url:  "https://charts.jetstack.io"
		}
	}
}

#Config: {
	namespace: string
	datafile:  string
	chart: core.#Chart & {
		version: =~"^v{0,1}[0-9]+\\.[0-9]+\\.[0-9]+$"
	}
}
