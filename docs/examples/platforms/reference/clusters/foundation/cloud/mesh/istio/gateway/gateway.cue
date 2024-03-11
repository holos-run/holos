package holos

// The primary istio Gateway, named default

let Name = "gateway"

#InputKeys: component: Name

#TargetNamespace: "istio-ingress"
#DependsOn:       _IngressGateway

let LoginCert = #PlatformCerts.login

#KubernetesObjects & {
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
	}
}
