package holos

// Lets Encrypt certificate issuers for public tls certs
#InputKeys: component: "certissuers"
#TargetNamespace: "cert-manager"

let Name = "letsencrypt"

// The cloudflare api token is platform scoped, not cluster scoped.
#SecretName: "cloudflare-api-token-secret"

// Depends on cert manager
#DependsOn: _CertManager

#KubernetesObjects & {
	apiObjects: {
		ClusterIssuer: {
			letsencrypt: #ClusterIssuer & {
				metadata: name: Name
				spec: {
					acme: {
						email:  #Platform.org.contact.email
						server: "https://acme-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name + "-istio"
						solvers: [{http01: ingress: class: "istio"}]
					}
				}
			}
			letsencryptStaging: #ClusterIssuer & {
				metadata: name: Name + "-staging"
				spec: {
					acme: {
						email:  #Platform.org.contact.email
						server: "https://acme-staging-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name + "-staging-istio"
						solvers: [{http01: ingress: class: "istio"}]
					}
				}
			}
			letsencryptDns: #ClusterIssuer & {
				metadata: name: Name + "-dns"
				spec: {
					acme: {
						email:  #Platform.org.contact.email
						server: "https://acme-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name + "-istio"
						solvers: [{
							dns01: cloudflare: {
								email: #Platform.org.cloudflare.email
								apiTokenSecretRef: name: #SecretName
								apiTokenSecretRef: key:  "api_token"
							}}]
					}
				}
			}
		}
		ExternalSecret: "\(#SecretName)": #ExternalSecret & {
			_name: #SecretName
		}
	}
}
