package holos

// Produce a Kustomize BuildPlan for Holos
holos: Kustomize.BuildPlan

// https://github.com/mccutchen/go-httpbin/blob/v2.15.0/kustomize/README.md
Kustomize: #Kustomize & {
	KustomizeConfig: {
		// Files tells Holos to copy the file from the component path to the
		// temporary directory Holos uses for BuildPlan execution.
		Files: {
			"httpbin.yaml": _
		}
		CommonLabels: {
			"app.kubernetes.io/name": "httpbin"
		}
		// Kustomization represents a kustomization.yaml file in CUE.  Holos
		// marshals this field into a `kustomization.yaml` while processing a
		// BuildPlan.  See
		// https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/
		Kustomization: {
			images: [{name: "mccutchen/go-httpbin"}]
			// Use a hidden field to compose patches easily with a struct.  Hidden
			// fields are not included in exported structures.
			_patches: {}
			// Convert the hidden struct to a list.
			patches: [for x in _patches {x}]
		}
	}
}
