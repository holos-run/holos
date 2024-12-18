@if(!NoKargo)
package holos

import stage "kargo.akuity.io/stage/v1alpha1"

// #StageSpecBuilder builds a Kargo Stage spec field from a provided component
// name.  Useful to define how the #Stage.spec field is produced for each
// component in a Kargo Project.
//
// This is meant to be a simple builder useful as a starting example, you don't
// have to use it and can instead define your own to produce the desired kargo
// stage spec.
#StageSpecBuilder: {
	// Component name without the stage prefix, e.g. "podinfo"
	Component: string
	// Prior stage to promote from.  If "direct", freight is promoted directly
	// from the Warehouse.
	Prior: string
	// The warehouse to pull from, defaults to the component name.
	Warehouse: {
		kind: "Warehouse"
		name: string | *Component
	}

	spec: stage.#StageSpec & {
		requestedFreight: [{
			origin: Warehouse
			if Prior == "direct" {sources: direct: true}
			if Prior != "direct" {sources: stages: ["\(Prior)-\(Component)"]}
		}]
	}
}

// Mix Kargo Stage resource fields into the holos stage definition.  The
// #KargoProjectBuilder definition uses these fields to construct the promotion
// pipeline for each promotable component in the project.
#Stage: {
	prior: string
}
