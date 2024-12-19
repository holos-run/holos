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
			"cert-manager":     _
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
				instances: [{extractYAML: path: "./data/cert-manager/cert-manager.yaml"}]
			}
			"local-ca": {
				name: "local-ca"
				path: "projects/security/components/local-ca"
			}
		}
	}
}
