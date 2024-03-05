package holos

#PlatformCerts: {
	// Login service for IAM.
	login: #PlatformCert & {
		_name:        "login"
		_wildcard:    true
		_description: "Cert for Zitadel the platform oidc identity provider for iam"
	}
	"k2-httpbin": #ClusterCert & {
		_name:        "httpbin"
		_cluster:     "k2"
		_description: "Test endpoint to verify the service mesh ingress gateway"
	}
}

// #PlatformCert provisions a cert in the provisioner cluster.  Workload clusters use ExternalSecret resources to fetch the Secret tls key and cert from the provisioner cluster.
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
	metadata: name:      string | *"\(_cluster)-\(_name)"
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
