// Code generated by timoni. DO NOT EDIT.

//timoni:generate timoni vendor crd -f /home/jeff/workspace/holos-run/holos-infra/deploy/clusters/k2/components/prod-secrets-eso/prod-secrets-eso.gen.yaml

package v1alpha1

import "strings"

#PushSecret: {
	// APIVersion defines the versioned schema of this representation
	// of an object.
	// Servers should convert recognized schemas to the latest
	// internal value, and
	// may reject unrecognized values.
	// More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	apiVersion: "external-secrets.io/v1alpha1"

	// Kind is a string value representing the REST resource this
	// object represents.
	// Servers may infer this from the endpoint the client submits
	// requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	kind: "PushSecret"
	metadata!: {
		name!: strings.MaxRunes(253) & strings.MinRunes(1) & {
			string
		}
		namespace!: strings.MaxRunes(63) & strings.MinRunes(1) & {
			string
		}
		labels?: {
			[string]: string
		}
		annotations?: {
			[string]: string
		}
	}

	// PushSecretSpec configures the behavior of the PushSecret.
	spec!: #PushSecretSpec
}

// PushSecretSpec configures the behavior of the PushSecret.
#PushSecretSpec: {
	// Secret Data that should be pushed to providers
	data?: [...{
		// Match a given Secret Key to be pushed to the provider.
		match: {
			// Remote Refs to push to providers.
			remoteRef: {
				// Name of the property in the resulting secret
				property?: string

				// Name of the resulting provider secret.
				remoteKey: string
			}

			// Secret Key to be pushed
			secretKey?: string
		}

		// Metadata is metadata attached to the secret.
		// The structure of metadata is provider specific, please look it
		// up in the provider documentation.
		metadata?: _
	}]

	// Deletion Policy to handle Secrets in the provider. Possible
	// Values: "Delete/None". Defaults to "None".
	deletionPolicy?: "Delete" | "None" | *"None"

	// The Interval to which External Secrets will try to push a
	// secret definition
	refreshInterval?: string
	secretStoreRefs: [...{
		// Kind of the SecretStore resource (SecretStore or
		// ClusterSecretStore)
		// Defaults to `SecretStore`
		kind?: string | *"SecretStore"

		// Optionally, sync to secret stores with label selector
		labelSelector?: {
			// matchExpressions is a list of label selector requirements. The
			// requirements are ANDed.
			matchExpressions?: [...{
				// key is the label key that the selector applies to.
				key: string

				// operator represents a key's relationship to a set of values.
				// Valid operators are In, NotIn, Exists and DoesNotExist.
				operator: string

				// values is an array of string values. If the operator is In or
				// NotIn,
				// the values array must be non-empty. If the operator is Exists
				// or DoesNotExist,
				// the values array must be empty. This array is replaced during a
				// strategic
				// merge patch.
				values?: [...string]
			}]

			// matchLabels is a map of {key,value} pairs. A single {key,value}
			// in the matchLabels
			// map is equivalent to an element of matchExpressions, whose key
			// field is "key", the
			// operator is "In", and the values array contains only "value".
			// The requirements are ANDed.
			matchLabels?: {
				[string]: string
			}
		}

		// Optionally, sync to the SecretStore of the given name
		name?: string
	}]
	selector: {
		secret: {
			// Name of the Secret. The Secret must exist in the same namespace
			// as the PushSecret manifest.
			name: string
		}
	}

	// Template defines a blueprint for the created Secret resource.
	template?: {
		data?: {
			[string]: string
		}

		// EngineVersion specifies the template engine version
		// that should be used to compile/execute the
		// template specified in .data and .templateFrom[].
		engineVersion?: "v1" | "v2" | *"v2"
		mergePolicy?:   "Replace" | "Merge" | *"Replace"

		// ExternalSecretTemplateMetadata defines metadata fields for the
		// Secret blueprint.
		metadata?: {
			annotations?: {
				[string]: string
			}
			labels?: {
				[string]: string
			}
		}
		templateFrom?: [...{
			configMap?: {
				items: [...{
					key:         string
					templateAs?: "Values" | "KeysAndValues" | *"Values"
				}]
				name: string
			}
			literal?: string
			secret?: {
				items: [...{
					key:         string
					templateAs?: "Values" | "KeysAndValues" | *"Values"
				}]
				name: string
			}
			target?: "Data" | "Annotations" | "Labels" | *"Data"
		}]
		type?: string
	}
}
