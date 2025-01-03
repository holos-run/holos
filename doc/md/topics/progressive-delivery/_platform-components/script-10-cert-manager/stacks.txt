package platform

import "holos.example/config/certmanager"

stacks: security: (#StackBuilder & {
	(#PromoterBuilder & {parameters: {
		name: "cert-manager"
		config: {
			datafile: certmanager.config.datafile
			chart:    certmanager.config.chart
		}
	}}).promoter

	// Manage the external-secrets namespace.
	stack: namespaces: "external-secrets": _

	parameters: {
		name: "security"
		components: {
			namespaces: {
				path: "stacks/security/components/namespaces"
				annotations: description: "configures namespaces for all stacks"
			}
			"external-secrets-crds": {
				path: "stacks/security/components/external-secrets-crds"
				annotations: description: "external secrets custom resource definitions"
			}
			"external-secrets": {
				path: "stacks/security/components/external-secrets"
				annotations: description: "external secrets custom resource definitions"
			}
			"cert-manager": {
				path: "stacks/security/components/cert-manager"
				annotations: description: "cert-manager operator and custom resource definitions"
				parameters: {
					kargoProject: "cert-manager"
					kargoStage:   "main"
				}
			}
			"local-ca": {
				path: "stacks/security/components/local-ca"
				annotations: description: "localhost mkcert certificate authority"
			}
		}
	}
}).stack
