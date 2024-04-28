package forms

import formsv1 "github.com/holos-run/forms/v1alpha1"

let Platform = formsv1.#Platform & {
	name:        "bare"
	displayName: "Bare Platform"

	sections: org: {
		displayName: "Organization"
		description: "Organization config values are used to derive more specific configuration values throughout the platform."

		fieldConfigs: {
			// platform.org.name
			name: props: {
				label:       "Name"
				placeholder: "example"
				description: "DNS label, e.g. 'example'"
			}
			// platform.org.domain
			domain: props: {
				label:       "Domain"
				placeholder: "example.com"
				description: "DNS domain, e.g. 'example.com'"
			}
			// platform.org.displayName
			displayName: props: {
				label:       "Display Name"
				placeholder: "Example Organization"
				description: "Display name, e.g. 'Example Organization'"
			}
			// platform.org.contactEmail
			contactEmail: props: {
				label:       "Contact Email"
				placeholder: "platform-team@example.com"
				description: "Technical contact email address"
			}
		}
	}
}

// Provide the output form fields
Platform.Form
