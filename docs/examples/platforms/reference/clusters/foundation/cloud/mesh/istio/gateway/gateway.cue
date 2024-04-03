package holos

import "list"

// The primary istio Gateway, named default
let Name = "gateway"
#InputKeys: component: Name
#TargetNamespace: "istio-ingress"

let LoginCert = #PlatformCerts.login

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-stores":  _
		_dependsOn: "prod-mesh-istio-base": _
		_dependsOn: "prod-mesh-ingress":    _

		metadata: name: "\(#InstancePrefix)-\(Name)"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: login: #ExternalSecret & {
			_name: "login"
		}
		Gateway: default: #Gateway & {
			metadata: name:      "default"
			metadata: namespace: #TargetNamespace
			spec: selector: istio: "ingressgateway"
			spec: servers: [
				{
					hosts: [for dnsName in LoginCert.spec.dnsNames {"prod-iam-zitadel/\(dnsName)"}]
					port: name:          "https-prod-iam-login"
					port: number:        443
					port: protocol:      "HTTPS"
					tls: credentialName: LoginCert.spec.secretName
					tls: mode:           "SIMPLE"
				},
			]
		}

		for k, svc in #OptionalServices {
			if svc.enabled && list.Contains(svc.clusterNames, #ClusterName) {
				Gateway: "\(svc.name)": #Gateway & {
					metadata: name:      svc.name
					metadata: namespace: #TargetNamespace
					spec: selector: istio: "ingressgateway"
					spec: servers: [for s in svc.servers {s}]
				}
				for k, s in svc.servers {
					ExternalSecret: "\(s.tls.credentialName)": _
				}
			}
		}
	}
}
