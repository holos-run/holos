package holos

import ci "cert-manager.io/clusterissuer/v1"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "local-ca"
	Namespace: "cert-manager"

	Resources: {
		ClusterIssuer: {
			"local-ca": ci.#ClusterIssuer & {
				metadata: name: Name
				spec: ca: secretName: Name
			}
		}
	}
}
