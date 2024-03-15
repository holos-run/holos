package holos

let CoreDomain = "core.\(#Platform.org.domain)"
let TargetNamespace = "prod-core-vault"

#OptionalServices: {
	vault: {
		enabled: true
		clusters: core1: _
		clusters: core2: _
		managedNamespaces: "prod-core-vault": {
			namespace: metadata: labels: "istio-injection": "enabled"
		}
		certs: "vault-core": #Certificate & {
			metadata: name:      "vault-core"
			metadata: namespace: "istio-ingress"
			spec: {
				commonName: "vault.\(CoreDomain)"
				dnsNames: [commonName]
				secretName: metadata.name
				issuerRef: kind: "ClusterIssuer"
				issuerRef: name: string | *"letsencrypt"
			}
		}
		servers: "https-vault-core": {
			hosts: ["\(TargetNamespace)/vault.\(CoreDomain)"]
			tls: credentialName: certs."vault-core".spec.secretName
		}
		for k, v in clusters {
			let obj = (Cert & {Name: "vault-core", Cluster: v.name}).APIObject
			certs: "\(obj.metadata.name)": obj
			servers: "https-\(obj.metadata.name)": {
				hosts: [for host in obj.spec.dnsNames {"\(TargetNamespace)/\(host)"}]
				tls: credentialName: obj.spec.secretName
			}
		}
	}
}

// Cert provisions a cluster specific certificate.
let Cert = {
	Name:    string
	Cluster: string

	APIObject: #Certificate & {
		metadata: name:      "\(Cluster)-\(Name)"
		metadata: namespace: string | *"istio-ingress"
		spec: {
			commonName: string | *"vault.\(Cluster).\(CoreDomain)"
			dnsNames: [commonName]
			secretName: metadata.name
			issuerRef: kind: "ClusterIssuer"
			issuerRef: name: string | *"letsencrypt"
		}
	}
}
