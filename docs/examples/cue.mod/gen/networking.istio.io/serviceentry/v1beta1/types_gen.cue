// Code generated by timoni. DO NOT EDIT.

//timoni:generate timoni vendor crd -f /home/jeff/workspace/holos-run/holos-infra/deploy/clusters/k2/components/prod-mesh-istio-base/prod-mesh-istio-base.gen.yaml

package v1beta1

import "strings"

#ServiceEntry: {
	// Configuration affecting service registry. See more details at:
	// https://istio.io/docs/reference/config/networking/service-entry.html
	spec!:      #ServiceEntrySpec
	apiVersion: "networking.istio.io/v1beta1"
	kind:       "ServiceEntry"
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
}

// Configuration affecting service registry. See more details at:
// https://istio.io/docs/reference/config/networking/service-entry.html
#ServiceEntrySpec: {
	// The virtual IP addresses associated with the service.
	addresses?: [...string]

	// One or more endpoints associated with the service.
	endpoints?: [...{
		// Address associated with the network endpoint without the port.
		address?: string

		// One or more labels associated with the endpoint.
		labels?: {
			[string]: string
		}

		// The locality associated with the endpoint.
		locality?: string

		// Network enables Istio to group endpoints resident in the same
		// L3 domain/network.
		network?: string

		// Set of ports associated with the endpoint.
		ports?: {
			[string]: int
		}

		// The service account associated with the workload if a sidecar
		// is present in the workload.
		serviceAccount?: string

		// The load balancing weight associated with the endpoint.
		weight?: int
	}]

	// A list of namespaces to which this service is exported.
	exportTo?: [...string]

	// The hosts associated with the ServiceEntry.
	hosts: [...string]

	// Specify whether the service should be considered external to
	// the mesh or part of the mesh.
	location?: "MESH_EXTERNAL" | "MESH_INTERNAL"

	// The ports associated with the external service.
	ports?: [...{
		// Label assigned to the port.
		name: string

		// A valid non-negative integer port number.
		number: int

		// The protocol exposed on the port.
		protocol?: string

		// The port number on the endpoint where the traffic will be
		// received.
		targetPort?: int
	}]

	// Service resolution mode for the hosts.
	resolution?: "NONE" | "STATIC" | "DNS" | "DNS_ROUND_ROBIN"

	// If specified, the proxy will verify that the server
	// certificate's subject alternate name matches one of the
	// specified values.
	subjectAltNames?: [...string]
	workloadSelector?: {
		// One or more labels that indicate a specific set of pods/VMs on
		// which the configuration should be applied.
		labels?: {
			[string]: string
		}
	}
}