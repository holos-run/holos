package forms

import v1 "github.com/holos-run/holos/api/v1alpha1"

// Provides a concrete v1.#Form
FormBuilder.Output

let FormBuilder = v1.#FormBuilder & {
	Sections: org: {
		displayName: "Organization"
		description: "Organization config values are used to derive more specific configuration values throughout the platform."

		fieldConfigs: {
			// platform.spec.config.user.sections.org.fields.name
			name: {
				type:         "input"
				defaultValue: "demo"
				props: {
					label: "Name"
					// placeholder: "example" placeholder cannot be used with validation?
					description: "DNS label, e.g. 'example'"
					pattern:     "^[a-z]([0-9a-z]|-){1,28}[0-9a-z]$"
					minLength:   3
					maxLength:   30
					required:    true
				}
				validation: messages: {
					pattern:   "It must be \(props.minLength) to \(props.maxLength) lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited."
					minLength: "Must be at least \(props.minLength) characters"
					maxLength: "Must be at most \(props.maxLength) characters"
					required:  pattern
				}
			}

			displayName: {
				type:         "input"
				defaultValue: "Demo Organization"
				props: {
					label:       "Display Name"
					placeholder: "Example Organization"
					description: "Display name, e.g. 'Example Organization'"
					maxLength:   100
					required:    true
				}
				validation: messages: {
					required: "Enter a display name."
				}
			}

			domain: {
				type:         "input"
				defaultValue: "holos.local"
				props: {
					label:       "Domain Name"
					placeholder: "holos.local"
					description: "Domain name, e.g. 'holos.local'"
					pattern:     "^holos\\.local$"
					maxLength:   100
					required:    true
				}
				validation: messages: {
					required: "Must be '\(defaultValue)' for local evaluation."
					pattern:  required
				}
			}
		}
	}

	Sections: argocd: {
		displayName: "ArgoCD"
		description: "Configure ArgoCD platform settings."

		fieldConfigs: {
			clientID: {
				type:         "input"
				defaultValue: "275571279971577269@holos_quickstart"
				props: {
					label:       "Client ID"
					description: "Client ID of the ArgoCD application in the ZITADEL project."
					pattern:     "^275571279971577269@holos_quickstart$"
					minLength:   1
					maxLength:   512
					required:    true
				}
				validation: messages: {
					minLength: "Must be at least \(props.minLength) characters"
					maxLength: "Must be at most \(props.maxLength) characters"
					required:  "Must be '\(defaultValue)' for local evaluation."
					pattern:   required
				}
			}

			repoURL: {
				type:         "input"
				defaultValue: "https://github.com/holos-run/holos-k3d"
				props: {
					label:       "Git repository URL"
					description: "Git repository URL, e.g. 'git@github.com:myorg/holos-infra.git'"
					minLength:   3
					maxLength:   128
					required:    true
				}
				validation: messages: {
					minLength: "Must be at least \(props.minLength) characters"
					maxLength: "Must be at most \(props.maxLength) characters"
				}
			}

			deployRoot: {
				type:         "input"
				defaultValue: ""
				props: {
					label:       "Deploy Root"
					description: "Path to the parent directory of the holos deploy directory."
					pattern:     "^$"
					required:    false
				}
				validation: messages: {
					pattern: "Must be empty for local evaluation."
				}
			}

			targetRevision: {
				type:         "input"
				defaultValue: "HEAD"
				props: {
					label:       "Target Revision"
					description: "Git reference to deploy."
					required:    true
				}
			}
		}
	}
}
