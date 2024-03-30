package holos

#MeshConfig: {
	projects:    _
	clusterName: _

	extensionProviderMap: {
		"cluster-trace": {
			zipkin: {
				maxTagLength: 56
				port:         9411
				service:      "zipkin.istio-system.svc"
			}
		}
	}

	config: {
		// For PROXY PROTOCOL at the ingress gateway.
		gatewayTopology: {
			numTrustedProxies: 2
		}
	}

	// Configure an ExtAuthzHttp provider for each stage's authproxy
	for Project in projects {
		if Project.clusters[clusterName] != _|_ {
			for Stage in Project.stages {
				let Name = "\(Stage.slug)-authproxy"
				extensionProviderMap: (Name): #ExtAuthzProxy & {
					envoyExtAuthzHttp: service: "authproxy.\(Stage.namespace).svc.cluster.local"
				}
			}
		}
	}
}
