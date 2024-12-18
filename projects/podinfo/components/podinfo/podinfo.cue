package holos

import h "example.com/platform/schemas/holos/v1alpha5"

holos: Component.BuildPlan

_NamespaceName: string @tag(NamespaceName)
// image is the oci image repository to configure.
_image:   string            @tag(image)
_version: string | *"6.7.0" @tag(version)

TierName: Stages[StageName].tier

ValuesByTier: {
	prod: replicaCount:    2
	nonprod: replicaCount: 1
}

Component: #Helm & {
	Chart: {
		name:    "oci://ghcr.io/stefanprodan/charts/podinfo"
		release: "podinfo"
		version: _version
	}
	Values: #ComponentValues & {
		// Embed the values for this tier.
		ValuesByTier[Stages[StageName].tier]

		ui: {
			message: "Hello! Stage: \(StageName)"
		}
		image: {
			tag:        Chart.version
			repository: _image
		}
	}
	// Ensure all resources are located in this namespace.
	KustomizeConfig: Kustomization: namespace: _NamespaceName
	// Grant the HTTPRoute access to route to this namespace.
	Resources: ReferenceGrant: (h.#ReferenceGrantBuilder & {Namespace: _NamespaceName}).ReferenceGrant
}
