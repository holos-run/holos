package holos

// GitHub Actions Runner Controller

#TargetNamespace: "arc-system"

#InputKeys: project:   "github"
#InputKeys: component: "arc"

#Kustomization: spec: targetNamespace: #TargetNamespace
#DependsOn: Namespaces: name:          "prod-secrets-namespaces"

#HelmChart & {
	values: #Values & #DefaultSecurityContext
	namespace: #TargetNamespace
	chart: {
		// Match the gha-base-name in the chart _helpers.tpl to avoid long full names.
		release: "gha-rs-controller"
		name:    "oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller"
		version: "0.8.3"
	}
}
