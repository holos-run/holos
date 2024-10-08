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
					required:  pattern
				}
			}

			displayName: {
				type: "input"
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
				type: "input"
				props: {
					label:       "Domain Name"
					placeholder: "example.com"
					description: "Domain name, e.g. 'example.com'"
					pattern:     "^([a-z0-9]([a-z0-9]|-){0,61}[a-z0-9]{0,1}\\.?)+[a-z]{2,6}$"
					maxLength:   100
					required:    true
				}
				validation: messages: {
					required: "Enter a valid domain name."
					pattern:  required
				}
			}
		}
	}

	Sections: github: {
		displayName: "GitHub"
		description: "Configure the platform GitHub integration.  These values are used by the Backstage component."

		fieldConfigs: {
			primaryOrg: {
				type: "input"
				props: {
					label:       "Organization"
					description: "Primary GitHub orgranization where code repositories reside. \(validation.messages.required)"
					pattern:     "^[a-z]([a-z0-9]|-){\(minLength-2),\(maxLength-2)}[a-z]$"
					minLength:   2
					maxLength:   39
					required:    true
				}
				validation: messages: {
					pattern:  "It must be \(props.minLength) to \(props.maxLength) lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited. \(required)"
					required: "GitHub organization name"
				}
			}
		}
	}

	Sections: eso: {
		displayName: "Secret Store"
		description: "Configure the platform secret store.  These values are used by the external-secrets-creds component.  Note: this information is not sufficient to read secrets.  To read secrets, the credential refresher job requires the workload clusters to be configured as workload identity providers."

		fieldConfigs: {
			gcpProjectID: {
				type: "input"
				props: {
					label:       "GCP Project ID"
					description: "GCP Project ID of the management cluster. \(validation.messages.required)"
					pattern:     "^[a-z]([a-z0-9]|-){4,28}[a-z]$"
					minLength:   6
					maxLength:   30
					required:    true
				}
				validation: messages: {
					pattern:  "It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited. \(required)"
					required: "gcloud projects list"
				}

			}

			gcpProjectNumber: {
				type: "input"
				props: {
					label:       "GCP Project Number"
					description: "GCP Project Number of the management cluster.  \(validation.messages.required)"
					pattern:     "^[0-9]+$"
					required:    true
				}
				validation: messages: {
					pattern:  "Must be a positive integer.  \(required)"
					required: "gcloud projects list"
				}
			}

			gcpServiceAccount: {
				type: "input"
				props: {
					label:       "ESO Credential Refresher Service Account"
					placeholder: "eso-creds-refresher@my-project-id.iam.gserviceaccount.com"
					description: "GCP Service Account email for the external secrets creds refresher job.  \(validation.messages.required)"
					minLength:   3
					maxLength:   250
					required:    true
				}
				validation: messages: {
					required: "gcloud iam service-accounts list"
				}
			}

			gkeClusterName: {
				type: "input"
				props: {
					label:       "GKE Cluster Name"
					placeholder: "management"
					description: "GKE Cluster Name of the management cluster.  \(validation.messages.required)"
					pattern:     "^[a-z]$|^[a-z][0-9a-z]$|^[a-z]([0-9a-z]|-){0,38}[0-9a-z]$"
					minLength:   1
					maxLength:   40
					required:    true
				}
				validation: messages: {
					pattern:  "Lowercase letters, numbers, and hyphens only.  Must start with a letter.  Must end with a number or letter.  \(required)"
					required: "gcloud container clusters list"
				}
			}

			gkeRegion: {
				type: "select"
				props: {
					label:       "Select Region"
					description: "Select the GCP region of the management cluster."
					multiple:    false
					options:     GCPRegions
				}
			}

			gkeCABundle: {
				type: "input"
				props: {
					label:       "Management Cluster CA Bundle"
					description: "Enter the management cluster ca bundle.  \(validation.messages.required)"
					pattern:     "^[0-9a-zA-Z]+=*$"
					minLength:   1
					required:    true
				}
				validation: messages: {
					required: "kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.certificate-authority-data}'"
					pattern:  "Must be a base64 encoded pem encoded certificate bundle.  \(required)"
				}
			}

			gkeClusterURL: {
				type: "input"
				props: {
					label:       "Management Cluster URL"
					description: "Enter the URL of the management cluster API endpoint.  \(validation.messages.required)"
					pattern:     "^https://.*$"
					required:    true
				}
				validation: messages: {
					pattern:  "Must be a https:// URL.  \(required)"
					required: "kubectl config view --minify --flatten -ojsonpath='{.clusters[0].cluster.server}'"
				}
			}
		}
	}

	Sections: letsencrypt: {
		displayName: "Let's Encrypt"
		description: "Configure Let's Encrypt account settings to provision TLS certificates."

		fieldConfigs: {
			// platform.spec.config.user.sections.org.fields.name
			email: {
				type: "input"
				props: {
					label: "Account Email"
					// placeholder: "example" placeholder cannot be used with validation?
					description: "Let's Encrypt Account Email, e.g. 'platform@example.com'"
					minLength:   3
					maxLength:   100
					required:    true
				}
				validation: messages: {
					minLength: "Must be at least \(props.minLength) characters"
					maxLength: "Must be at most \(props.maxLength) characters"
				}
			}
		}
	}

	Sections: zitadel: {
		displayName: "ZITADEL"
		description: "Configure platform level ZITADEL settings.  Requires AWS."

		fieldConfigs: {
			backupBucketName: {
				type: "input"
				props: {
					label:       "ZITADEL Backup Bucket Name"
					description: "S3 bucket name, e.g. 'myorg-zitadel-backups'"
					pattern:     "^[a-z]([a-z0-9.-]){1,61}[0-9a-z]$"
					minLength:   3
					maxLength:   63
					required:    true
				}
				validation: messages: {
					pattern:   "It must be \(props.minLength) to \(props.maxLength) lowercase letters, numbers, dots, or hyphens. It must start with a letter. It must end with a letter or number."
					minLength: "Must be at least \(props.minLength) characters"
					maxLength: "Must be at most \(props.maxLength) characters"
					required:  pattern
				}
			}
			backupBucketRegion: {
				type: "select"
				props: {
					label:       "ZITADEL Backup Bucket Region"
					description: "Select the S3 bucket region for database backups."
					multiple:    false
					options:     AWSRegions
					required:    true
				}
			}
			backupRestore: {
				type: "checkbox"
				props: {
					label:       "Provision from Backup"
					description: "If checked, the PostgresCluster will provision from backups.  Uncheck this box for initial provisioning.  Check it after a successful backup.  NOTE: Moving from false to true will trigger an in-place restore from the most recent backup.  Data since the most recent backup will be lost."
					required:    true
				}
			}
		}
	}

	Sections: authproxy: {
		displayName: "Auth Proxy"
		description: "Enter configuration values from the Auth Proxy application in the Holos Platform project in the platform identity provider (ZITADEL)."

		fieldConfigs: clientID: (#StandardFields & {DisplayName: displayName}).clientID
		fieldConfigs: {
			projectID: {
				type: "input"
				props: {
					label:       "Project ID"
					description: "Project Resource ID of the Holos Platform project in ZITADEL."
					pattern:     "^[0-9]+$"
					minLength:   1
					maxLength:   100
					required:    true
				}
				validation: messages: {
					minLength: "Must be at least \(props.minLength) characters"
					maxLength: "Must be at most \(props.maxLength) characters"
					pattern:   "Must be a positive whole number."
					required:  "Required: " + props.description
				}
			}
		}
	}

	Sections: argocd: {
		displayName: "ArgoCD"
		description: "Configure ArgoCD platform settings."

		fieldConfigs: {
			clientID: (#StandardFields & {DisplayName: displayName}).clientID

			repoURL: {
				type: "input"
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
					required:    false
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

	Sections: aws: {
		displayName: "AWS"
		description: "Configure AWS settings."

		fieldConfigs: {
			accountNumber: {
				type: "input"
				props: {
					label:       "AWS Account Number"
					description: "AWS Account Number.  \(validation.messages.required)"
					pattern:     "^[0-9]+$"
					required:    true
				}
				validation: messages: {
					pattern:  "Must be a positive integer.  \(required)"
					required: "aws sts get-caller-identity"
				}
			}
			primaryRegion: {
				type: "select"
				props: {
					label:       "Primary Region"
					description: "Select the primary region for AWS resources."
					multiple:    false
					options:     AWSRegions
					required:    true
				}
			}
		}
	}
}

#StandardFields: {
	DisplayName: string
	clientID: {
		type: "input"
		props: {
			label:       "Client ID"
			description: "Client ID of the \(DisplayName) application in the Holos Platform ZITADEL project."
			pattern:     "^[0-9]+@[a-z0-9_]+$"
			minLength:   1
			maxLength:   512
			required:    true
		}
		validation: messages: {
			minLength: "Must be at least \(props.minLength) characters"
			maxLength: "Must be at most \(props.maxLength) characters"
			pattern:   "Must be like 999999999999999999@holos_platform"
			required:  "Required: " + props.description
		}
	}
}

let GCPRegions = [
	{value: "africa-south1", label: "africa-south1"},
	{value: "asia-east1", label: "asia-east1"},
	{value: "asia-east2", label: "asia-east2"},
	{value: "asia-northeast1", label: "asia-northeast1"},
	{value: "asia-northeast2", label: "asia-northeast2"},
	{value: "asia-northeast3", label: "asia-northeast3"},
	{value: "asia-south1", label: "asia-south1"},
	{value: "asia-south2", label: "asia-south2"},
	{value: "asia-southeast1", label: "asia-southeast1"},
	{value: "asia-southeast2", label: "asia-southeast2"},
	{value: "australia-southeast1", label: "australia-southeast1"},
	{value: "australia-southeast2", label: "australia-southeast2"},
	{value: "europe-central2", label: "europe-central2"},
	{value: "europe-north1", label: "europe-north1"},
	{value: "europe-southwest1", label: "europe-southwest1"},
	{value: "europe-west1", label: "europe-west1"},
	{value: "europe-west10", label: "europe-west10"},
	{value: "europe-west12", label: "europe-west12"},
	{value: "europe-west2", label: "europe-west2"},
	{value: "europe-west3", label: "europe-west3"},
	{value: "europe-west4", label: "europe-west4"},
	{value: "europe-west6", label: "europe-west6"},
	{value: "europe-west8", label: "europe-west8"},
	{value: "europe-west9", label: "europe-west9"},
	{value: "me-central1", label: "me-central1"},
	{value: "me-central2", label: "me-central2"},
	{value: "me-west1", label: "me-west1"},
	{value: "northamerica-northeast1", label: "northamerica-northeast1"},
	{value: "northamerica-northeast2", label: "northamerica-northeast2"},
	{value: "southamerica-east1", label: "southamerica-east1"},
	{value: "southamerica-west1", label: "southamerica-west1"},
	{value: "us-central1", label: "us-central1"},
	{value: "us-east1", label: "us-east1"},
	{value: "us-east4", label: "us-east4"},
	{value: "us-east5", label: "us-east5"},
	{value: "us-south1", label: "us-south1"},
	{value: "us-west1", label: "us-west1"},
	{value: "us-west2", label: "us-west2"},
	{value: "us-west3", label: "us-west3"},
	{value: "us-west4", label: "us-west4"},
]

let AWSRegions = [
	{value: "us-east-1", label: "N. Virginia (us-east-1)"},
	{value: "us-east-2", label: "Ohio (us-east-2)"},
	{value: "us-west-1", label: "N. California (us-west-1)"},
	{value: "us-west-2", label: "Oregon (us-west-2)"},
	{value: "us-gov-west1", label: "US GovCloud West (us-gov-west1)"},
	{value: "us-gov-east1", label: "US GovCloud East (us-gov-east1)"},
	{value: "ca-central-1", label: "Canada (ca-central-1)"},
	{value: "eu-north-1", label: "Stockholm (eu-north-1)"},
	{value: "eu-west-1", label: "Ireland (eu-west-1)"},
	{value: "eu-west-2", label: "London (eu-west-2)"},
	{value: "eu-west-3", label: "Paris (eu-west-3)"},
	{value: "eu-central-1", label: "Frankfurt (eu-central-1)"},
	{value: "eu-south-1", label: "Milan (eu-south-1)"},
	{value: "af-south-1", label: "Cape Town (af-south-1)"},
	{value: "ap-northeast-1", label: "Tokyo (ap-northeast-1)"},
	{value: "ap-northeast-2", label: "Seoul (ap-northeast-2)"},
	{value: "ap-northeast-3", label: "Osaka (ap-northeast-3)"},
	{value: "ap-southeast-1", label: "Singapore (ap-southeast-1)"},
	{value: "ap-southeast-2", label: "Sydney (ap-southeast-2)"},
	{value: "ap-east-1", label: "Hong Kong (ap-east-1)"},
	{value: "ap-south-1", label: "Mumbai (ap-south-1)"},
	{value: "me-south-1", label: "Bahrain (me-south-1)"},
	{value: "sa-east-1", label: "São Paulo (sa-east-1)"},
	{value: "cn-north-1", label: "Bejing (cn-north-1)"},
	{value: "cn-northwest-1", label: "Ningxia (cn-northwest-1)"},
	{value: "ap-southeast-3", label: "Jakarta (ap-southeast-3)"},
]
