@if(!NoIstio)
package holos

import "github.com/holos-run/holos/api/core/v1alpha5:core"

Projects: network: {
	namespaces: (Istio.System.Namespace):  _
	namespaces: (Istio.Gateway.Namespace): _
	namespaces: istio: metadata: labels: "kargo.akuity.io/project": "true"

	let DATAFILE = {instances: [{extractYAML: path: Istio.datafile}]}

	components: {
		"istio-base": {
			name: "istio-base"
			path: "projects/network/components/istio-base"
			DATAFILE
		}
		"istiod": {
			name: "istiod"
			path: "projects/network/components/istiod"
			DATAFILE
		}
		"istio-cni": {
			name: "istio-cni"
			path: "projects/network/components/istio-cni"
			DATAFILE
		}
		"istio-ztunnel": {
			name: "istio-ztunnel"
			path: "projects/network/components/istio-ztunnel"
			DATAFILE
		}
		"istio-gateway": {
			name: "istio-gateway"
			path: "projects/network/components/istio-gateway"
			DATAFILE
		}
		"istio-kargo": {
			name: "istio-promoter"
			path: "components/addon-promoter"
			parameters: {
				KargoProjectName: "istio"
				KargoStageName:   "main"
				KargoDataFile:    Istio.datafile
				KargoDataKey:     "Istio.chart.version"
				GitRepoURL:       Organization.RepoURL
				ChartName:        "base"
				ChartRepoURL:     Istio.chart.repository.url
			}
			DATAFILE
		}
		"httproutes": {
			name: "httproutes"
			path: "projects/network/components/httproutes"
		}
	}
}

// #Istio represents platform wide configuration
Istio: #Istio & {
	System: Namespace:  "istio-system"
	Gateway: Namespace: "istio-ingress"

	datafile: "./data/istio/istio.yaml"
	chart: {
		// defined in the datafile
		version: string
		repository: {
			name: "istio"
			url:  "https://istio-release.storage.googleapis.com/charts"
		}
	}

	// Constrain Helm values for safer, easier upgrades and consistency across
	// platform components.
	Values: global: istioNamespace: System.Namespace
	// Configure ambient mode
	Values: profile: "ambient"
}

#Istio: {
	Version: chart.version
	System: Namespace:  string
	Gateway: Namespace: string

	datafile: string
	chart: core.#Chart & {
		version: string
		repository: {
			name: string
			url:  string
		}
	}

	Values: {...}
}
