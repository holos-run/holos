package holos

import ci "cert-manager.io/clusterissuer/v1"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

// The cloudflare api token is platform scoped, not cluster scoped.
#SecretName: "cloudflare-api-token-secret"

let Objects = {
	Name:      "letsencrypt"
	Namespace: "cert-manager"

	Resources: {
		ClusterIssuer: {
			letsencrypt: ci.#ClusterIssuer & {
				metadata: name: Name
				spec: {
					acme: {
						email:  _Platform.Model.letsencrypt.email
						server: "https://acme-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name
						solvers: [{
							dns01: cloudflare: {
								email: spec.acme.email
								apiTokenSecretRef: name: #SecretName
								apiTokenSecretRef: key:  "api_token"
							}}]
					}
				}
			}
			letsencryptStaging: ci.#ClusterIssuer & {
				metadata: name: Name + "-staging"
				spec: {
					acme: {
						email:  _Platform.Model.letsencrypt.email
						server: "https://acme-staging-v02.api.letsencrypt.org/directory"
						privateKeySecretRef: name: Name + "-staging"
						solvers: [{
							dns01: cloudflare: {
								email: spec.acme.email
								apiTokenSecretRef: name: #SecretName
								apiTokenSecretRef: key:  "api_token"
							}}]
					}
				}
			}
		}
	}
}
