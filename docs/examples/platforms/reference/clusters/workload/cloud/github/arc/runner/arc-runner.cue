package holos

#InputKeys: component: "arc-runner"
#Kustomization: spec: targetNamespace: #TargetNamespace

let GithubOrg = #Platform.org.github.orgs.primary.name

#HelmChart & {
	values: {
		#Values
		controllerServiceAccount: name:      "gha-rs-controller"
		controllerServiceAccount: namespace: "arc-system"
		githubConfigSecret: "controller-manager"
		githubConfigUrl:    "https://github.com/\(GithubOrg)"
	}
	apiObjects: {
		ExternalSecret: controller: #ExternalSecret & {
			_name: values.githubConfigSecret
		}
	}
	chart: {
		// Match the gha-base-name in the chart _helpers.tpl to avoid long full names.
		release: "gha-rs"
		name:    "oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set"
	}
}
