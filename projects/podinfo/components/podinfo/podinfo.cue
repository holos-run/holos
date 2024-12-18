package holos

holos: Component.BuildPlan

_NamespaceName: string @tag(NamespaceName)
// image is the oci image repository to configure.
_image:   string            @tag(image)
_version: string | *"6.7.0" @tag(version)

Component: #Helm & {
	Chart: {
		name:    "oci://ghcr.io/stefanprodan/charts/podinfo"
		release: "podinfo"
		version: _version
	}
	Values: #ComponentValues & {
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
	Resources: ReferenceGrant: (#ReferenceGrantBuilder & {Namespace: _NamespaceName}).ReferenceGrant
}
