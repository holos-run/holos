package holos

import certv1 "cert-manager.io/certificate/v1"

let Objects = {
	Name:      "certificates"
	Namespace: "istio-gateways"

	Resources: Certificate: [NAME=string]: certv1.#Certificate & {
		metadata: name:      NAME
		metadata: namespace: Namespace
		spec: {
			commonName: NAME
			secretName: NAME
			dnsNames: [NAME]
			issuerRef: {
				kind: "ClusterIssuer"
				name: "local-ca"
			}
		}
	}

	Resources: Certificate: "httpbin.\(_Platform.Model.org.domain)":   _
	Resources: Certificate: "argocd.\(_Platform.Model.org.domain)":    _
	Resources: Certificate: "app.\(_Platform.Model.org.domain)":       _
	Resources: Certificate: "backstage.\(_Platform.Model.org.domain)": _
}

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output
