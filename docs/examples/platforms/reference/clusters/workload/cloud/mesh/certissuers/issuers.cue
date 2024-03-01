package holos

// Lets Encrypt certificate issuers for public tls certs
#InputKeys: component: "certissuers"
#TargetNamespace: "cert-manager"

#KubernetesObjects & {
	apiObjects: {
		ClusterIssuer: {
			letsencrypt: {

			}
		}
	}
}
