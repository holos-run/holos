package holos

import h "example.com/platform/schemas/holos/v1alpha5"

// Holos specific integration goes in this file.
Component: Resources: {
	// Grant the HTTPRoute access to route to this namespace.
	ReferenceGrant: (h.#ReferenceGrantBuilder & {Namespace: Parameters.namespace}).ReferenceGrant
}
