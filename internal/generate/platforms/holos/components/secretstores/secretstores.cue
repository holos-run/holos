package holos

import ss "external-secrets.io/secretstore/v1beta1"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "secretstores"
	Namespace: "default"

	for Namespace in _Namespaces {
		Resources: SecretStore: "\(Namespace.metadata.name)": #SecretStore & {
			metadata: namespace: Namespace.metadata.name
		}
	}
}

#SecretStore: ss.#SecretStore & {
	metadata: {
		name:      "default"
		namespace: string
	}
	spec: provider: {
		kubernetes: {
			remoteNamespace: metadata.namespace
			auth: token: bearerToken: {
				name: "eso-reader"
				key:  "token"
			}
			server: {
				caBundle: _Platform.Model.eso.gkeCABundle
				url:      _Platform.Model.eso.gkeClusterURL
			}
		}
	}
}
