package holos

let Name = "httpbin"
let Host = Name + "." + #ClusterDomain

#InputKeys: component: Name

#TargetNamespace: "istio-ingress"
#DependsOn:       _IngressGateway

#Metadata: namespace: #TargetNamespace
SecretName: #InputKeys.cluster + "-" + Name

#KubernetesObjects & {
	apiObjects: {
		Certificate: {
			httpbin: #Certificate & {
				metadata: {
					#Metadata
					name: SecretName
				}
				spec: {
					commonName: Host
					dnsNames: [Host]
					secretName: SecretName
					issuerRef: kind: "ClusterIssuer"
					issuerRef: name: "letsencrypt"
				}
			}
		}
	}
}
