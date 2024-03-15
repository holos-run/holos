package holos

#TargetNamespace: #ARCSystemNamespace
#InputKeys: component: "arc-system"

#HelmChart & {
	values:    #Values & #DefaultSecurityContext
	namespace: #TargetNamespace
	chart: {
		// Match the gha-base-name in the chart _helpers.tpl to avoid long full names.
		release: "gha-rs-controller"
		name:    "oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller"
		version: "0.8.3"
	}
}
