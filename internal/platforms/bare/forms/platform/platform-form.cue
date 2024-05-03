package forms

import formsv1 "github.com/holos-run/forms/v1alpha1"

let Platform = formsv1.#Platform & {
	name:        "bare"
	displayName: "Bare Platform"

	sections: org: {
		displayName: "Organization"
		description: "Organization config values are used to derive more specific configuration values throughout the platform."

		fieldConfigs: {
			// platform.spec.config.user.sections.org.fields.name
			name: props: {
				label:       "Name"
				placeholder: "example"
				description: "DNS label, e.g. 'example'"
			}
			// platform.spec.config.user.sections.org.fields.domain
			domain: props: {
				label:       "Domain"
				placeholder: "example.com"
				description: "DNS domain, e.g. 'example.com'"
			}
			// platform.spec.config.user.sections.org.fields.displayName
			displayName: props: {
				label:       "Display Name"
				placeholder: "Example Organization"
				description: "Display name, e.g. 'Example Organization'"
			}
			// platform.spec.config.user.sections.org.fields.contactEmail
			contactEmail: {
				props: {
					label:       "Contact Email"
					placeholder: "platform-team@example.com"
					description: "Technical contact email address"
				}
			}
		}
	}

	sections: privacy: {
		displayName: "Data Privacy"
		description: "Configure data privacy aspects of the platform."

		fieldConfigs: {
			country: {
				// https://formly.dev/docs/api/ui/material/select/
				type: "select"
				props: {
					label:       "Select Planet"
					description: "Juridiction of applicable data privacy laws."
					options: [
						{value: "mercury", label: "Mercury"},
						{value: "venus", label:   "Venus"},
						{value: "earth", label:   "Earth"},
						{value: "mars", label:    "Mars"},
						{value: "jupiter", label: "Jupiter"},
						{value: "saturn", label:  "Saturn"},
						{value: "uranus", label:  "Uranus"},
						{value: "neptune", label: "Neptune"},
					]
				}
			}
			regions: {
				// https://formly.dev/docs/api/ui/material/select/
				type: "select"
				props: {
					label:           "Select Regions"
					description:     "Select the regions this platform operates in."
					multiple:        true
					selectAllOption: "Select All"
					options: [
						{value: "us-east-2", label: "Ohio"},
						{value: "us-west-2", label: "Oregon"},
						{value: "eu-west-1", label: "Ireland"},
						{value: "eu-west-2", label: "London", disabled: true},
					]
				}
			}
		}
	}

	// https://v5.formly.dev/ui/material
	sections: terms: {
		displayName: "Terms and Conditions"
		description: "Example of a boolean checkbox."

		fieldConfigs: {
			// platform.spec.config.user.sections.terms.fields.didAgree
			didAgree: {
				type: "checkbox"
				props: {
					label:       "Accept terms"
					description: "In order to proceed, please accept terms"
					pattern:     "true"
					required:    true
				}
				validation: {
					messages: {
						pattern: "Please accept the terms"
					}
				}
			}
		}
	}
}

// Provide the output form fields
Platform.Form
