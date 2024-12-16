package holos

// Produce a helm chart build plan.
holos: Component.BuildPlan

Component: #Kubernetes & {
	Resources: [_]: [_]: metadata: namespace: Kargo.Namespace

	Resources: {
		ExternalSecret: creds: {
			metadata: name: Kargo.Values.api.secret.name
			spec: {
				refreshInterval: "24h"
				target: {
					creationPolicy: "Owner"
					deletionPolicy: "Delete"
					template: {
						type:          "Opaque"
						mergePolicy:   "Merge"
						engineVersion: "v2"
						data: ADMIN_ACCOUNT_PASSWORD_HASH: "{{ .password | bcrypt }}"
					}
				}
				dataFrom: [
					{
						// Specify the top level key for the generated value.  This key is
						// used in the ExternalSecret.spec.target.template.data templates.
						rewrite: [{transform: template: "password"}]
						sourceRef: {
							generatorRef: {
								apiVersion: "generators.external-secrets.io/v1alpha1"
								kind:       "Password"
								name:       Password.admin.metadata.name
							}
						}
					},
					{
						// Specify the top level key for the generated value.  This key is
						// used in the ExternalSecret.spec.target.template.data templates.
						rewrite: [{transform: template: "ADMIN_ACCOUNT_TOKEN_SIGNING_KEY"}]
						sourceRef: {
							generatorRef: {
								apiVersion: "generators.external-secrets.io/v1alpha1"
								kind:       "Password"
								name:       Password.signingKey.metadata.name
							}
						}
					},
				]
			}
		}

		Password: admin: {
			metadata: name: "admin-password-generator"
			spec: SPEC
		}
		Password: signingKey: {
			metadata: name: "admin-signing-key-generator"
			spec: SPEC
		}

		let SPEC = {
			length:      32
			digits:      4
			symbols:     4
			allowRepeat: true
			noUpper:     false
		}
	}
}
