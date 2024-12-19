package holos

import holos "example.com/platform/schemas/holos/v1alpha5"

ProjectName: string | *"no-project" @tag(ProjectName)

Projects: holos.#Projects & {
	network: {
		components: "gateway-api": {
			name: "gateway-api"
			path: "projects/network/components/gateway-api"
		}
	}

	security: {
		namespaces: {
			"external-secrets": _
			(CertManager.namespace): metadata: labels: "kargo.akuity.io/project": "true"
		}
		components: {
			"namespaces": {
				name: "namespaces"
				path: "projects/security/components/namespaces"
			}
			"external-secrets-crds": {
				name: "external-secrets-crds"
				path: "projects/security/components/external-secrets-crds"
			}
			"external-secrets": {
				name: "external-secrets"
				path: "projects/security/components/external-secrets"
			}
			"cert-manager": {
				name: "cert-manager"
				path: "projects/security/components/cert-manager"
				instances: [{extractYAML: path: CertManager.datafile}]
				parameters: KargoProjectName: "cert-manager"
				parameters: KargoStageName:   "main"
			}
			// Defines the Kargo promotion pipeline for cert-manager.
			"cert-manager-kargo": {
				name: "cert-manager-kargo"
				path: "projects/security/components/cert-manager-kargo"
				parameters: KargoProjectName: "cert-manager"
				parameters: KargoStageName:   "main"
			}
			"local-ca": {
				name: "local-ca"
				path: "projects/security/components/local-ca"
			}
		}
	}
}
