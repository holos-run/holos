package holos

#TargetNamespace: "arc-runner"
#InputKeys: component: "arc-runner"
#Kustomization: spec: targetNamespace: #TargetNamespace

let GitHubConfigSecret = "controller-manager"

// Just sync the external secret, don't configure the scale set
// Work around https://github.com/actions/actions-runner-controller/issues/3351
if #IsPrimaryCluster == false {
	spec: components: KubernetesObjectsList: [
		#KubernetesObjects & {
			metadata: name:                        "prod-github-arc-runner"
			_dependsOn: "prod-secrets-namespaces": _

			apiObjectMap: (#APIObjects & {
				apiObjects: ExternalSecret: "\(GitHubConfigSecret)": _
			}).apiObjectMap
		},
	]
}

// Put the scale set on the primary cluster.
if #IsPrimaryCluster == true {
	spec: components: HelmChartList: [
		#HelmChart & {
			_dependsOn: "prod-secrets-namespaces": _
			metadata: name:                        "prod-github-arc-runner"
			_values: {
				#Values
				controllerServiceAccount: name:      "gha-rs-controller"
				controllerServiceAccount: namespace: "arc-system"
				githubConfigSecret: GitHubConfigSecret
				githubConfigUrl:    "https://github.com/" + #Platform.org.github.orgs.primary.name
			}
			apiObjectMap: (#APIObjects & {apiObjects: ExternalSecret: "\(_values.githubConfigSecret)": _}).apiObjectMap
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
		},
	]
}
