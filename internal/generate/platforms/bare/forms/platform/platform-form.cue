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
				type: "input"
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
				}
			}

			// platform.spec.config.user.sections.org.fields.displayName
			displayName: {
				type: "input"
				props: {
					label:       "Display Name"
					placeholder: "Example Organization"
					description: "Display name, e.g. 'Example Organization'"
					maxLength:   100
					required:    true
				}
			}
		}
	}

	Sections: cloud: {
		displayName: "Cloud Providers"
		description: "Select the services that provide resources for the platform."

		fieldConfigs: {
			providers: {
				// https://formly.dev/docs/api/ui/material/select/
				type: "select"
				props: {
					label:           "Select Providers"
					description:     "Select the cloud providers the platform builds upon."
					multiple:        true
					selectAllOption: "Select All"
					options: [
						{value: "aws", label:        "Amazon Web Services"},
						{value: "gcp", label:        "Google Cloud Platform"},
						{value: "azure", label:      "Microsoft Azure"},
						{value: "cloudflare", label: "Cloudflare"},
						{value: "github", label:     "GitHub"},
						{value: "ois", label:        "Open Infrastructure Services"},
						{value: "onprem", label:     "On Premises", disabled: true},
					]
				}
			}
		}
	}

	Sections: aws: {
		displayName: "Amazon Web Services"
		description: "Provide the information necessary for Holos to manage AWS resources to provide the platform."

		expressions: hide: "!\(AWSSelected)"

		fieldConfigs: {
			primaryRoleARN: {
				// https://formly.dev/docs/api/ui/material/input
				type: "input"
				props: {
					label:       "Holos Admin Role ARN"
					description: "Enter the AWS Role ARN Holos will use to bootstrap resources.  For example, arn:aws:iam::123456789012:role/HolosAdminAccess"
					pattern:     "^arn:.*"
					minLength:   4
					required:    true
				}
				validation: messages: {
					pattern: "Must be a valid ARN.  Refer to https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html"
				}
			}

			regions: {
				// https://formly.dev/docs/api/ui/material/select/
				type: "select"
				props: {
					label:           "Select Regions"
					description:     "Select the AWS regions this platform operates in."
					multiple:        true
					required:        true
					selectAllOption: "Select All"
					options:         AWSRegions
				}
			}
		}
	}

	Sections: gcp: {
		displayName: "Google Cloud Platform"
		description: "Use this form to configure platform level GCP settings."

		expressions: hide: "!\(GCPSelected)"

		fieldConfigs: {
			regions: {
				// https://formly.dev/docs/api/ui/material/select/
				type: "select"
				props: {
					label:           "Select Regions"
					description:     "Select the GCP regions this platform operates in."
					multiple:        true
					selectAllOption: "Select All"
					// gcloud compute regions list --format=json | jq '.[] | {value: .name, label: .description}' regions.json | jq -s | cue export --out cue
					options: GCPRegions
				}
			}

			gcpProjectID: {
				// https://formly.dev/docs/api/ui/material/input
				type: "input"
				props: {
					label:       "Project ID"
					description: "Enter the project id where the provisioner cluster resides."
					pattern:     "^[a-z]([0-9a-z]|-){1,28}[0-9a-z]$"
					minLength:   6
					maxLength:   30
					required:    true
				}
				validation: messages: {
					pattern:   "It must be \(props.minLength) to \(props.maxLength) lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited."
					minLength: "Must be at least \(props.minLength) characters."
					maxLength: "Must be at most \(props.maxLength) characters."
				}
			}

			gcpProjectNumber: {
				// https://formly.dev/docs/api/ui/material/input
				type: "input"
				props: {
					label: "Project Number"
					// note type number here
					type:        "number"
					description: "Enter the project number where the provisioner cluster resides."
					pattern:     "^[0-9]+$"
					required:    true
				}
				validation: messages: {
					pattern: "Must be a valid project number."
				}
			}

			provisionerCABundle: {
				type: "input"
				props: {
					label:       "Provisioner CA Bundle"
					description: "Enter the provisioner cluster ca bundle.  kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.certificate-authority-data}'"
					pattern:     "^[0-9a-zA-Z]+=*$"
					required:    true
				}
				validation: messages: {
					pattern: "Must be a base64 encoded pem encoded certificate bundle."
				}
			}

			provisionerURL: {
				type: "input"
				props: {
					label:       "Provisioner URL"
					description: "Enter the URL of the provisioner cluster API endpoint.  kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.server}'"
					pattern:     "^https://.*$"
					required:    true
				}
				validation: messages: {
					pattern: "Must be a https:// URL."
				}
			}
		}
	}

	Sections: cloudflare: {
		displayName: "Cloudflare"
		description: "Cloudflare is primarily used for DNS automation."

		expressions: hide: "!" + CloudflareSelected

		fieldConfigs: {
			email: {
				// https://formly.dev/docs/api/ui/material/input
				type: "input"
				props: {
					label:       "Account Email"
					description: "Enter the Cloudflare email address to manage DNS"
					minLength:   3
					required:    true
				}
			}
		}
	}

	Sections: github: {
		displayName: "GitHub"
		description: "GitHub is primarily used to host Git repositories and execute Actions workflows."

		expressions: hide: "!\(GitHubSelected)"

		fieldConfigs: {
			primaryOrg: {
				// https://formly.dev/docs/api/ui/material/input
				type: "input"
				props: {
					label:       "Organization"
					description: "Enter the primary GitHub organization associed with the platform."
					pattern:     "^(?!-)(?!.*--)([a-zA-Z0-9]|-){1,39}$"
					minLength:   1
					maxLength:   39
					required:    true
				}
				validation: messages: {
					pattern: "All characters must be either a hyphen or alphanumeric.  Cannot start with a hyphen.  Cannot include consecutive hyphens."
				}
			}
		}
	}
}

let GCPRegions = [
	{value: "africa-south1", label:           "africa-south1"},
	{value: "asia-east1", label:              "asia-east1"},
	{value: "asia-east2", label:              "asia-east2"},
	{value: "asia-northeast1", label:         "asia-northeast1"},
	{value: "asia-northeast2", label:         "asia-northeast2"},
	{value: "asia-northeast3", label:         "asia-northeast3"},
	{value: "asia-south1", label:             "asia-south1"},
	{value: "asia-south2", label:             "asia-south2"},
	{value: "asia-southeast1", label:         "asia-southeast1"},
	{value: "asia-southeast2", label:         "asia-southeast2"},
	{value: "australia-southeast1", label:    "australia-southeast1"},
	{value: "australia-southeast2", label:    "australia-southeast2"},
	{value: "europe-central2", label:         "europe-central2"},
	{value: "europe-north1", label:           "europe-north1"},
	{value: "europe-southwest1", label:       "europe-southwest1"},
	{value: "europe-west1", label:            "europe-west1"},
	{value: "europe-west10", label:           "europe-west10"},
	{value: "europe-west12", label:           "europe-west12"},
	{value: "europe-west2", label:            "europe-west2"},
	{value: "europe-west3", label:            "europe-west3"},
	{value: "europe-west4", label:            "europe-west4"},
	{value: "europe-west6", label:            "europe-west6"},
	{value: "europe-west8", label:            "europe-west8"},
	{value: "europe-west9", label:            "europe-west9"},
	{value: "me-central1", label:             "me-central1"},
	{value: "me-central2", label:             "me-central2"},
	{value: "me-west1", label:                "me-west1"},
	{value: "northamerica-northeast1", label: "northamerica-northeast1"},
	{value: "northamerica-northeast2", label: "northamerica-northeast2"},
	{value: "southamerica-east1", label:      "southamerica-east1"},
	{value: "southamerica-west1", label:      "southamerica-west1"},
	{value: "us-central1", label:             "us-central1"},
	{value: "us-east1", label:                "us-east1"},
	{value: "us-east4", label:                "us-east4"},
	{value: "us-east5", label:                "us-east5"},
	{value: "us-south1", label:               "us-south1"},
	{value: "us-west1", label:                "us-west1"},
	{value: "us-west2", label:                "us-west2"},
	{value: "us-west3", label:                "us-west3"},
	{value: "us-west4", label:                "us-west4"},
]

let AWSRegions = [
	{value: "us-east-1", label:      "N. Virginia (us-east-1)"},
	{value: "us-east-2", label:      "Ohio (us-east-2)"},
	{value: "us-west-1", label:      "N. California (us-west-1)"},
	{value: "us-west-2", label:      "Oregon (us-west-2)"},
	{value: "us-gov-west1", label:   "US GovCloud West (us-gov-west1)"},
	{value: "us-gov-east1", label:   "US GovCloud East (us-gov-east1)"},
	{value: "ca-central-1", label:   "Canada (ca-central-1)"},
	{value: "eu-north-1", label:     "Stockholm (eu-north-1)"},
	{value: "eu-west-1", label:      "Ireland (eu-west-1)"},
	{value: "eu-west-2", label:      "London (eu-west-2)"},
	{value: "eu-west-3", label:      "Paris (eu-west-3)"},
	{value: "eu-central-1", label:   "Frankfurt (eu-central-1)"},
	{value: "eu-south-1", label:     "Milan (eu-south-1)"},
	{value: "af-south-1", label:     "Cape Town (af-south-1)"},
	{value: "ap-northeast-1", label: "Tokyo (ap-northeast-1)"},
	{value: "ap-northeast-2", label: "Seoul (ap-northeast-2)"},
	{value: "ap-northeast-3", label: "Osaka (ap-northeast-3)"},
	{value: "ap-southeast-1", label: "Singapore (ap-southeast-1)"},
	{value: "ap-southeast-2", label: "Sydney (ap-southeast-2)"},
	{value: "ap-east-1", label:      "Hong Kong (ap-east-1)"},
	{value: "ap-south-1", label:     "Mumbai (ap-south-1)"},
	{value: "me-south-1", label:     "Bahrain (me-south-1)"},
	{value: "sa-east-1", label:      "São Paulo (sa-east-1)"},
	{value: "cn-north-1", label:     "Bejing (cn-north-1)"},
	{value: "cn-northwest-1", label: "Ningxia (cn-northwest-1)"},
	{value: "ap-southeast-3", label: "Jakarta (ap-southeast-3)"},
]

let AWSSelected = "formState.model.cloud?.providers?.includes(\"aws\")"
let GCPSelected = "formState.model.cloud?.providers?.includes(\"gcp\")"
let GitHubSelected = "formState.model.cloud?.providers?.includes(\"github\")"
let CloudflareSelected = "formState.model.cloud?.providers?.includes(\"cloudflare\")"
