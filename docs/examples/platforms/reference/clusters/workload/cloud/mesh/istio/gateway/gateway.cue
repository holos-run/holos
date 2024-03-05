package holos

// The primary istio Gateway, named default

let Name = "gateway"

#InputKeys: component: Name

#TargetNamespace: "istio-ingress"
#DependsOn:       _IngressGateway

// TODO: We need to generalize this for multiple services hanging off the default gateway.
let LoginCert = #Certificate & {
	metadata: {
		name:      "login"
		namespace: #TargetNamespace
	}
	spec: {
		commonName: "login.\(#Platform.org.domain)"
		dnsNames: [commonName]
		secretName: metadata.name
		issuerRef: kind: "ClusterIssuer"
		issuerRef: name: "letsencrypt"
	}
}

#KubernetesObjects & {
	apiObjects: {
		Certificate: login: LoginCert
		Gateway: default: #Gateway & {
			metadata: name:      "default"
			metadata: namespace: #TargetNamespace
			spec: selector: istio: "ingressgateway"
			spec: servers: [
				{
					hosts: ["prod-iam-zitadel/\(LoginCert.spec.commonName)"]
					port: name:          "https-prod-iam-zitadel"
					port: number:        443
					port: protocol:      "HTTPS"
					tls: credentialName: LoginCert.spec.secretName
					tls: mode:           "SIMPLE"
				},
			]
		}
	}
}
