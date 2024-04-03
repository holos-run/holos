package holos

import "list"

// The primary istio Gateway, named default
let Name = "gateway"
#InputKeys: component: Name
#TargetNamespace: "istio-ingress"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-stores":  _
		_dependsOn: "prod-mesh-istio-base": _
		_dependsOn: "prod-mesh-ingress":    _

		metadata: name: "\(#InstancePrefix)-\(Name)"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

// GatewayServers represents all hosts for all VirtualServices in the cluster attached to Gateway/default
// NOTE: This is a critical structure because the default Gateway should be used in most cases.
let GatewayServers = {
	for Project in _Projects {
		for server in (#ProjectTemplate & {project: Project}).ClusterGatewayServers {
			(server.port.name): server
		}
	}

	for k, svc in #OptionalServices {
		if svc.enabled && list.Contains(svc.clusterNames, #ClusterName) {
			for server in svc.servers {
				(server.port.name): server
			}
		}
	}

	if #PlatformServers[#ClusterName] != _|_ {
		for server in #PlatformServers[#ClusterName] {
			(server.port.name): server
		}
	}
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		Gateway: default: #Gateway & {
			metadata: name:      "default"
			metadata: namespace: #TargetNamespace

			spec: selector: istio: "ingressgateway"
			spec: servers: [for x in GatewayServers {x}]
		}

		for k, svc in #OptionalServices {
			if svc.enabled && list.Contains(svc.clusterNames, #ClusterName) {
				for k, s in svc.servers {
					ExternalSecret: "\(s.tls.credentialName)": _
				}
			}
		}
	}
}
