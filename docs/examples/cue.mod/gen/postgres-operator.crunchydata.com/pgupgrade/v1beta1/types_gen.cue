// Code generated by timoni. DO NOT EDIT.

//timoni:generate timoni vendor crd -f /home/jeff/workspace/holos-run/holos-infra/deploy/clusters/core2/components/prod-pgo-crds/prod-pgo-crds.gen.yaml

package v1beta1

import "strings"

// PGUpgrade is the Schema for the pgupgrades API
#PGUpgrade: {
	// APIVersion defines the versioned schema of this representation
	// of an object. Servers should convert recognized schemas to the
	// latest internal value, and may reject unrecognized values.
	// More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	apiVersion: "postgres-operator.crunchydata.com/v1beta1"

	// Kind is a string value representing the REST resource this
	// object represents. Servers may infer this from the endpoint
	// the client submits requests to. Cannot be updated. In
	// CamelCase. More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	kind: "PGUpgrade"
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

	// PGUpgradeSpec defines the desired state of PGUpgrade
	spec!: #PGUpgradeSpec
}

// PGUpgradeSpec defines the desired state of PGUpgrade
#PGUpgradeSpec: {
	// Scheduling constraints of the PGUpgrade pod. More info:
	// https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node
	affinity?: {
		// Describes node affinity scheduling rules for the pod.
		nodeAffinity?: {
			// The scheduler will prefer to schedule pods to nodes that
			// satisfy the affinity expressions specified by this field, but
			// it may choose a node that violates one or more of the
			// expressions. The node that is most preferred is the one with
			// the greatest sum of weights, i.e. for each node that meets all
			// of the scheduling requirements (resource request,
			// requiredDuringScheduling affinity expressions, etc.), compute
			// a sum by iterating through the elements of this field and
			// adding "weight" to the sum if the node matches the
			// corresponding matchExpressions; the node(s) with the highest
			// sum are the most preferred.
			preferredDuringSchedulingIgnoredDuringExecution?: [...{
				// A node selector term, associated with the corresponding weight.
				preference: {
					// A list of node selector requirements by node's labels.
					matchExpressions?: [...{
						// The label key that the selector applies to.
						key: string

						// Represents a key's relationship to a set of values. Valid
						// operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
						operator: string

						// An array of string values. If the operator is In or NotIn, the
						// values array must be non-empty. If the operator is Exists or
						// DoesNotExist, the values array must be empty. If the operator
						// is Gt or Lt, the values array must have a single element,
						// which will be interpreted as an integer. This array is
						// replaced during a strategic merge patch.
						values?: [...string]
					}]

					// A list of node selector requirements by node's fields.
					matchFields?: [...{
						// The label key that the selector applies to.
						key: string

						// Represents a key's relationship to a set of values. Valid
						// operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
						operator: string

						// An array of string values. If the operator is In or NotIn, the
						// values array must be non-empty. If the operator is Exists or
						// DoesNotExist, the values array must be empty. If the operator
						// is Gt or Lt, the values array must have a single element,
						// which will be interpreted as an integer. This array is
						// replaced during a strategic merge patch.
						values?: [...string]
					}]
				}

				// Weight associated with matching the corresponding
				// nodeSelectorTerm, in the range 1-100.
				weight: int
			}]
			requiredDuringSchedulingIgnoredDuringExecution?: {
				// Required. A list of node selector terms. The terms are ORed.
				nodeSelectorTerms: [...{
					// A list of node selector requirements by node's labels.
					matchExpressions?: [...{
						// The label key that the selector applies to.
						key: string

						// Represents a key's relationship to a set of values. Valid
						// operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
						operator: string

						// An array of string values. If the operator is In or NotIn, the
						// values array must be non-empty. If the operator is Exists or
						// DoesNotExist, the values array must be empty. If the operator
						// is Gt or Lt, the values array must have a single element,
						// which will be interpreted as an integer. This array is
						// replaced during a strategic merge patch.
						values?: [...string]
					}]

					// A list of node selector requirements by node's fields.
					matchFields?: [...{
						// The label key that the selector applies to.
						key: string

						// Represents a key's relationship to a set of values. Valid
						// operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
						operator: string

						// An array of string values. If the operator is In or NotIn, the
						// values array must be non-empty. If the operator is Exists or
						// DoesNotExist, the values array must be empty. If the operator
						// is Gt or Lt, the values array must have a single element,
						// which will be interpreted as an integer. This array is
						// replaced during a strategic merge patch.
						values?: [...string]
					}]
				}]
			}
		}

		// Describes pod affinity scheduling rules (e.g. co-locate this
		// pod in the same node, zone, etc. as some other pod(s)).
		podAffinity?: {
			// The scheduler will prefer to schedule pods to nodes that
			// satisfy the affinity expressions specified by this field, but
			// it may choose a node that violates one or more of the
			// expressions. The node that is most preferred is the one with
			// the greatest sum of weights, i.e. for each node that meets all
			// of the scheduling requirements (resource request,
			// requiredDuringScheduling affinity expressions, etc.), compute
			// a sum by iterating through the elements of this field and
			// adding "weight" to the sum if the node has pods which matches
			// the corresponding podAffinityTerm; the node(s) with the
			// highest sum are the most preferred.
			preferredDuringSchedulingIgnoredDuringExecution?: [...{
				// Required. A pod affinity term, associated with the
				// corresponding weight.
				podAffinityTerm: {
					// A label query over a set of resources, in this case pods.
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
							// NotIn, the values array must be non-empty. If the operator is
							// Exists or DoesNotExist, the values array must be empty. This
							// array is replaced during a strategic merge patch.
							values?: [...string]
						}]

						// matchLabels is a map of {key,value} pairs. A single {key,value}
						// in the matchLabels map is equivalent to an element of
						// matchExpressions, whose key field is "key", the operator is
						// "In", and the values array contains only "value". The
						// requirements are ANDed.
						matchLabels?: {
							[string]: string
						}
					}

					// A label query over the set of namespaces that the term applies
					// to. The term is applied to the union of the namespaces
					// selected by this field and the ones listed in the namespaces
					// field. null selector and null or empty namespaces list means
					// "this pod's namespace". An empty selector ({}) matches all
					// namespaces.
					namespaceSelector?: {
						// matchExpressions is a list of label selector requirements. The
						// requirements are ANDed.
						matchExpressions?: [...{
							// key is the label key that the selector applies to.
							key: string

							// operator represents a key's relationship to a set of values.
							// Valid operators are In, NotIn, Exists and DoesNotExist.
							operator: string

							// values is an array of string values. If the operator is In or
							// NotIn, the values array must be non-empty. If the operator is
							// Exists or DoesNotExist, the values array must be empty. This
							// array is replaced during a strategic merge patch.
							values?: [...string]
						}]

						// matchLabels is a map of {key,value} pairs. A single {key,value}
						// in the matchLabels map is equivalent to an element of
						// matchExpressions, whose key field is "key", the operator is
						// "In", and the values array contains only "value". The
						// requirements are ANDed.
						matchLabels?: {
							[string]: string
						}
					}

					// namespaces specifies a static list of namespace names that the
					// term applies to. The term is applied to the union of the
					// namespaces listed in this field and the ones selected by
					// namespaceSelector. null or empty namespaces list and null
					// namespaceSelector means "this pod's namespace".
					namespaces?: [...string]

					// This pod should be co-located (affinity) or not co-located
					// (anti-affinity) with the pods matching the labelSelector in
					// the specified namespaces, where co-located is defined as
					// running on a node whose value of the label with key
					// topologyKey matches that of any node on which any of the
					// selected pods is running. Empty topologyKey is not allowed.
					topologyKey: string
				}

				// weight associated with matching the corresponding
				// podAffinityTerm, in the range 1-100.
				weight: int
			}]

			// If the affinity requirements specified by this field are not
			// met at scheduling time, the pod will not be scheduled onto the
			// node. If the affinity requirements specified by this field
			// cease to be met at some point during pod execution (e.g. due
			// to a pod label update), the system may or may not try to
			// eventually evict the pod from its node. When there are
			// multiple elements, the lists of nodes corresponding to each
			// podAffinityTerm are intersected, i.e. all terms must be
			// satisfied.
			requiredDuringSchedulingIgnoredDuringExecution?: [...{
				// A label query over a set of resources, in this case pods.
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
						// NotIn, the values array must be non-empty. If the operator is
						// Exists or DoesNotExist, the values array must be empty. This
						// array is replaced during a strategic merge patch.
						values?: [...string]
					}]

					// matchLabels is a map of {key,value} pairs. A single {key,value}
					// in the matchLabels map is equivalent to an element of
					// matchExpressions, whose key field is "key", the operator is
					// "In", and the values array contains only "value". The
					// requirements are ANDed.
					matchLabels?: {
						[string]: string
					}
				}

				// A label query over the set of namespaces that the term applies
				// to. The term is applied to the union of the namespaces
				// selected by this field and the ones listed in the namespaces
				// field. null selector and null or empty namespaces list means
				// "this pod's namespace". An empty selector ({}) matches all
				// namespaces.
				namespaceSelector?: {
					// matchExpressions is a list of label selector requirements. The
					// requirements are ANDed.
					matchExpressions?: [...{
						// key is the label key that the selector applies to.
						key: string

						// operator represents a key's relationship to a set of values.
						// Valid operators are In, NotIn, Exists and DoesNotExist.
						operator: string

						// values is an array of string values. If the operator is In or
						// NotIn, the values array must be non-empty. If the operator is
						// Exists or DoesNotExist, the values array must be empty. This
						// array is replaced during a strategic merge patch.
						values?: [...string]
					}]

					// matchLabels is a map of {key,value} pairs. A single {key,value}
					// in the matchLabels map is equivalent to an element of
					// matchExpressions, whose key field is "key", the operator is
					// "In", and the values array contains only "value". The
					// requirements are ANDed.
					matchLabels?: {
						[string]: string
					}
				}

				// namespaces specifies a static list of namespace names that the
				// term applies to. The term is applied to the union of the
				// namespaces listed in this field and the ones selected by
				// namespaceSelector. null or empty namespaces list and null
				// namespaceSelector means "this pod's namespace".
				namespaces?: [...string]

				// This pod should be co-located (affinity) or not co-located
				// (anti-affinity) with the pods matching the labelSelector in
				// the specified namespaces, where co-located is defined as
				// running on a node whose value of the label with key
				// topologyKey matches that of any node on which any of the
				// selected pods is running. Empty topologyKey is not allowed.
				topologyKey: string
			}]
		}

		// Describes pod anti-affinity scheduling rules (e.g. avoid
		// putting this pod in the same node, zone, etc. as some other
		// pod(s)).
		podAntiAffinity?: {
			// The scheduler will prefer to schedule pods to nodes that
			// satisfy the anti-affinity expressions specified by this field,
			// but it may choose a node that violates one or more of the
			// expressions. The node that is most preferred is the one with
			// the greatest sum of weights, i.e. for each node that meets all
			// of the scheduling requirements (resource request,
			// requiredDuringScheduling anti-affinity expressions, etc.),
			// compute a sum by iterating through the elements of this field
			// and adding "weight" to the sum if the node has pods which
			// matches the corresponding podAffinityTerm; the node(s) with
			// the highest sum are the most preferred.
			preferredDuringSchedulingIgnoredDuringExecution?: [...{
				// Required. A pod affinity term, associated with the
				// corresponding weight.
				podAffinityTerm: {
					// A label query over a set of resources, in this case pods.
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
							// NotIn, the values array must be non-empty. If the operator is
							// Exists or DoesNotExist, the values array must be empty. This
							// array is replaced during a strategic merge patch.
							values?: [...string]
						}]

						// matchLabels is a map of {key,value} pairs. A single {key,value}
						// in the matchLabels map is equivalent to an element of
						// matchExpressions, whose key field is "key", the operator is
						// "In", and the values array contains only "value". The
						// requirements are ANDed.
						matchLabels?: {
							[string]: string
						}
					}

					// A label query over the set of namespaces that the term applies
					// to. The term is applied to the union of the namespaces
					// selected by this field and the ones listed in the namespaces
					// field. null selector and null or empty namespaces list means
					// "this pod's namespace". An empty selector ({}) matches all
					// namespaces.
					namespaceSelector?: {
						// matchExpressions is a list of label selector requirements. The
						// requirements are ANDed.
						matchExpressions?: [...{
							// key is the label key that the selector applies to.
							key: string

							// operator represents a key's relationship to a set of values.
							// Valid operators are In, NotIn, Exists and DoesNotExist.
							operator: string

							// values is an array of string values. If the operator is In or
							// NotIn, the values array must be non-empty. If the operator is
							// Exists or DoesNotExist, the values array must be empty. This
							// array is replaced during a strategic merge patch.
							values?: [...string]
						}]

						// matchLabels is a map of {key,value} pairs. A single {key,value}
						// in the matchLabels map is equivalent to an element of
						// matchExpressions, whose key field is "key", the operator is
						// "In", and the values array contains only "value". The
						// requirements are ANDed.
						matchLabels?: {
							[string]: string
						}
					}

					// namespaces specifies a static list of namespace names that the
					// term applies to. The term is applied to the union of the
					// namespaces listed in this field and the ones selected by
					// namespaceSelector. null or empty namespaces list and null
					// namespaceSelector means "this pod's namespace".
					namespaces?: [...string]

					// This pod should be co-located (affinity) or not co-located
					// (anti-affinity) with the pods matching the labelSelector in
					// the specified namespaces, where co-located is defined as
					// running on a node whose value of the label with key
					// topologyKey matches that of any node on which any of the
					// selected pods is running. Empty topologyKey is not allowed.
					topologyKey: string
				}

				// weight associated with matching the corresponding
				// podAffinityTerm, in the range 1-100.
				weight: int
			}]

			// If the anti-affinity requirements specified by this field are
			// not met at scheduling time, the pod will not be scheduled onto
			// the node. If the anti-affinity requirements specified by this
			// field cease to be met at some point during pod execution (e.g.
			// due to a pod label update), the system may or may not try to
			// eventually evict the pod from its node. When there are
			// multiple elements, the lists of nodes corresponding to each
			// podAffinityTerm are intersected, i.e. all terms must be
			// satisfied.
			requiredDuringSchedulingIgnoredDuringExecution?: [...{
				// A label query over a set of resources, in this case pods.
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
						// NotIn, the values array must be non-empty. If the operator is
						// Exists or DoesNotExist, the values array must be empty. This
						// array is replaced during a strategic merge patch.
						values?: [...string]
					}]

					// matchLabels is a map of {key,value} pairs. A single {key,value}
					// in the matchLabels map is equivalent to an element of
					// matchExpressions, whose key field is "key", the operator is
					// "In", and the values array contains only "value". The
					// requirements are ANDed.
					matchLabels?: {
						[string]: string
					}
				}

				// A label query over the set of namespaces that the term applies
				// to. The term is applied to the union of the namespaces
				// selected by this field and the ones listed in the namespaces
				// field. null selector and null or empty namespaces list means
				// "this pod's namespace". An empty selector ({}) matches all
				// namespaces.
				namespaceSelector?: {
					// matchExpressions is a list of label selector requirements. The
					// requirements are ANDed.
					matchExpressions?: [...{
						// key is the label key that the selector applies to.
						key: string

						// operator represents a key's relationship to a set of values.
						// Valid operators are In, NotIn, Exists and DoesNotExist.
						operator: string

						// values is an array of string values. If the operator is In or
						// NotIn, the values array must be non-empty. If the operator is
						// Exists or DoesNotExist, the values array must be empty. This
						// array is replaced during a strategic merge patch.
						values?: [...string]
					}]

					// matchLabels is a map of {key,value} pairs. A single {key,value}
					// in the matchLabels map is equivalent to an element of
					// matchExpressions, whose key field is "key", the operator is
					// "In", and the values array contains only "value". The
					// requirements are ANDed.
					matchLabels?: {
						[string]: string
					}
				}

				// namespaces specifies a static list of namespace names that the
				// term applies to. The term is applied to the union of the
				// namespaces listed in this field and the ones selected by
				// namespaceSelector. null or empty namespaces list and null
				// namespaceSelector means "this pod's namespace".
				namespaces?: [...string]

				// This pod should be co-located (affinity) or not co-located
				// (anti-affinity) with the pods matching the labelSelector in
				// the specified namespaces, where co-located is defined as
				// running on a node whose value of the label with key
				// topologyKey matches that of any node on which any of the
				// selected pods is running. Empty topologyKey is not allowed.
				topologyKey: string
			}]
		}
	}

	// The major version of PostgreSQL before the upgrade.
	fromPostgresVersion: uint & >=10 & <=16

	// The image name to use for major PostgreSQL upgrades.
	image?: string

	// ImagePullPolicy is used to determine when Kubernetes will
	// attempt to pull (download) container images. More info:
	// https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy
	imagePullPolicy?: "Always" | "Never" | "IfNotPresent"

	// The image pull secrets used to pull from a private registry.
	// Changing this value causes all running PGUpgrade pods to
	// restart.
	// https://k8s.io/docs/tasks/configure-pod-container/pull-image-private-registry/
	imagePullSecrets?: [...{
		// Name of the referent. More info:
		// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
		name?: string
	}]

	// Metadata contains metadata for custom resources
	metadata?: {
		annotations?: {
			[string]: string
		}
		labels?: {
			[string]: string
		}
	}

	// The name of the cluster to be updated
	postgresClusterName: strings.MinRunes(1)

	// Priority class name for the PGUpgrade pod. Changing this value
	// causes PGUpgrade pod to restart. More info:
	// https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/
	priorityClassName?: string

	// Resource requirements for the PGUpgrade container.
	resources?: {
		// Limits describes the maximum amount of compute resources
		// allowed. More info:
		// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
		limits?: {
			[string]: (int | string) & =~"^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$"
		}

		// Requests describes the minimum amount of compute resources
		// required. If Requests is omitted for a container, it defaults
		// to Limits if that is explicitly specified, otherwise to an
		// implementation-defined value. More info:
		// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
		requests?: {
			[string]: (int | string) & =~"^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$"
		}
	}

	// The image name to use for PostgreSQL containers after upgrade.
	// When omitted, the value comes from an operator environment
	// variable.
	toPostgresImage?: string

	// The major version of PostgreSQL to be upgraded to.
	toPostgresVersion: uint & >=10 & <=16

	// Tolerations of the PGUpgrade pod. More info:
	// https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration
	tolerations?: [...{
		// Effect indicates the taint effect to match. Empty means match
		// all taint effects. When specified, allowed values are
		// NoSchedule, PreferNoSchedule and NoExecute.
		effect?: string

		// Key is the taint key that the toleration applies to. Empty
		// means match all taint keys. If the key is empty, operator must
		// be Exists; this combination means to match all values and all
		// keys.
		key?: string

		// Operator represents a key's relationship to the value. Valid
		// operators are Exists and Equal. Defaults to Equal. Exists is
		// equivalent to wildcard for value, so that a pod can tolerate
		// all taints of a particular category.
		operator?: string

		// TolerationSeconds represents the period of time the toleration
		// (which must be of effect NoExecute, otherwise this field is
		// ignored) tolerates the taint. By default, it is not set, which
		// means tolerate the taint forever (do not evict). Zero and
		// negative values will be treated as 0 (evict immediately) by
		// the system.
		tolerationSeconds?: int

		// Value is the taint value the toleration matches to. If the
		// operator is Exists, the value should be empty, otherwise just
		// a regular string.
		value?: string
	}]
}