package holos

#InputKeys: component: "arc-runner"
#Kustomization: spec: targetNamespace: #TargetNamespace

#HelmChart & {
	values: {
		#Values
		controllerServiceAccount: name:      "gha-rs-controller"
		controllerServiceAccount: namespace: "arc-system"
		githubConfigSecret: "controller-manager"
		githubConfigUrl:    "https://github.com/" + #Platform.org.github.orgs.primary.name
	}
	apiObjects: ExternalSecret: "\(values.githubConfigSecret)": _
	chart: {
		// Match the gha-base-name in the chart _helpers.tpl to avoid long full names.
		// NOTE: Unfortunately the INSTALLATION_NAME is used as the helm release
		// name and GitHub removed support for runner labels, so the only way to
		// specify which runner a workflow runs on is using this helm release name.
		// The quote is "Update the INSTALLATION_NAME value carefully. You will use
		// the installation name as the value of runs-on in your workflows."  Refer to
		// https://docs.github.com/en/actions/hosting-your-own-runners/managing-self-hosted-runners-with-actions-runner-controller/quickstart-for-actions-runner-controller
		release: "gha-rs"
		name:    "oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set"
	}
}
