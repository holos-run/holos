package holos

holos: Component.BuildPlan

_NamespaceName: string @tag(NamespaceName)

Component: #Helm & {
	Chart: {
		name:    "oci://ghcr.io/stefanprodan/charts/podinfo"
		release: "podinfo"
		version: "6.6.1"
	}
	Values: {
		ui: {
			message: "Hello! Stage: \(StageName)"
		}
	}
	// Ensure all resources are located in this namespace.
	KustomizeConfig: Kustomization: namespace: _NamespaceName
	// Grant the HTTPRoute access to route to this namespace.
	Resources: ReferenceGrant: (#ReferenceGrantBuilder & {Namespace: _NamespaceName}).ReferenceGrant
}
