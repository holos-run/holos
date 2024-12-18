@if(!NoKargo)
package holos

import (
	project "kargo.akuity.io/project/v1alpha1"
	warehouse "kargo.akuity.io/warehouse/v1alpha1"
	stage "kargo.akuity.io/stage/v1alpha1"
)

// Custom resource definitions.
#Resources: {
	Project?: [NAME=string]: project.#Project & {
		metadata: name: string | *NAME
		spec: {...}
	}
	Warehouse?: [_]: warehouse.#Warehouse
	Stage?: [_]:     stage.#Stage
}

// Kargo config definition.
#Kargo: {
	Namespace: string
	Version: string | *"1.1.1"
	Values: {...}
}
