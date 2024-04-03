package holos

#PlatformServers: {
	for cluster in #Platform.clusters {
		(cluster.name): {
			"https-istio-ingress-httpbin": {
				let cert = #PlatformCerts[cluster.name+"-httpbin"]
				hosts: [for host in cert.spec.dnsNames {"istio-ingress/\(host)"}]
				port: name:          "https-istio-ingress-httpbin"
				port: number:        443
				port: protocol:      "HTTPS"
				tls: credentialName: cert.spec.secretName
				tls: mode:           "SIMPLE"
			}
		}
	}
}

#PlatformCerts: {
	// Globally scoped platform services are defined here.
	login: #PlatformCert & {
		_name:        "login"
		_wildcard:    true
		_description: "Cert for Zitadel oidc identity provider for iam services"
	}

	// Cluster scoped services are defined here.
	for cluster in #Platform.clusters {
		"\(cluster.name)-httpbin": #ClusterCert & {
			_name:        "httpbin"
			_cluster:     cluster.name
			_description: "Test endpoint to verify the service mesh ingress gateway"
		}
	}
}

// #PlatformCert provisions a cert in the provisioner cluster.
// Workload clusters use ExternalSecret resources to fetch the Secret tls key and cert from the provisioner cluster.
#PlatformCert: #Certificate & {
	_name:     string
	_wildcard: true | *false
	metadata: name:      string | *_name
	metadata: namespace: string | *"istio-ingress"
	spec: {
		commonName: string | *"\(_name).\(#Platform.org.domain)"
		if _wildcard {
			dnsNames: [commonName, "*.\(commonName)"]
		}
		if !_wildcard {
			dnsNames: [commonName]
		}
		secretName: metadata.name
		issuerRef: kind: "ClusterIssuer"
		issuerRef: name: string | *"letsencrypt"
	}
}

// #ClusterCert provisions a cluster specific certificate.
#ClusterCert: #Certificate & {
	_name:     string
	_cluster:  string
	_wildcard: true | *false
	// Enforce this value
	metadata: name:      "\(_cluster)-\(_name)"
	metadata: namespace: string | *"istio-ingress"
	spec: {
		commonName: string | *"\(_name).\(_cluster).\(#Platform.org.domain)"
		if _wildcard {
			dnsNames: [commonName, "*.\(commonName)"]
		}
		if !_wildcard {
			dnsNames: [commonName]
		}
		secretName: metadata.name
		issuerRef: kind: "ClusterIssuer"
		issuerRef: name: string | *"letsencrypt"
	}
}
