package holos

import ci "cert-manager.io/clusterissuer/v1"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "local-ca"
	Namespace: #CertManager.Namespace

	Resources: ClusterIssuer: LocalCA: ci.#ClusterIssuer & {
		metadata: name:      "local-ca"
		metadata: namespace: #CertManager.Namespace

		// The secret name must align with the local cluster guide at
		// https://holos.run/docs/guides/local-cluster/
		spec: ca: secretName: "local-ca"
	}
}
