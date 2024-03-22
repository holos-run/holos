package holos

// Lets Encrypt certificate issuers for public tls certs
#TargetNamespace: "cert-manager"

let Name = "letsencrypt"

// The cloudflare api token is platform scoped, not cluster scoped.
#SecretName: "cloudflare-api-token-secret"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "\(#InstancePrefix)-letsencrypt"

		_dependsOn: "prod-secrets-namespaces":        _
		_dependsOn: "\(#InstancePrefix)-certmanager": _
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		ClusterIssuer: {
			letsencrypt: #ClusterIssuer & {
				metadata: name: Name
				spec: {
					acme: {
						email:  #Platform.org.contact.email
						server: "https://acme-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name
						solvers: [{
							dns01: cloudflare: {
								email: #Platform.org.cloudflare.email
								apiTokenSecretRef: name: #SecretName
								apiTokenSecretRef: key:  "api_token"
							}}]
					}
				}
			}
			letsencryptStaging: #ClusterIssuer & {
				metadata: name: Name + "-staging"
				spec: {
					acme: {
						email:  #Platform.org.contact.email
						server: "https://acme-staging-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name + "-staging"
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
	}
}

// _HTTPSolvers are disabled in the provisioner cluster, dns is the method supported by holos.
_HTTPSolvers: {
	letsencryptHTTP: #ClusterIssuer & {
		metadata: name: Name + "-http"
		spec: {
			acme: {
				email:  #Platform.org.contact.email
				server: "https://acme-v02.api.letsencrypt.org/directory"
				privateKeySecretRef: name: Name
				solvers: [{http01: ingress: class: "istio"}]
			}
		}
	}
	letsencryptHTTPStaging: #ClusterIssuer & {
		metadata: name: Name + "-http-staging"
		spec: {
			acme: {
				email:  #Platform.org.contact.email
				server: "https://acme-staging-v02.api.letsencrypt.org/directory"
				privateKeySecretRef: name: Name + "-staging"
				solvers: [{http01: ingress: class: "istio"}]
			}
		}
	}
}
