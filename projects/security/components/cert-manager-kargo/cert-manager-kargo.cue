package holos

import "path"

Parameters: {
	KargoProjectName: string @tag(KargoProjectName)
	KargoStageName:   string @tag(KargoStageName)
}

holos: Component.BuildPlan

// Manage a Kargo Project and promotion stages for cert-manager.  The use case
// is to watch for new helm chart versions and submit a PR against the main
// branch with the fully rendered manifests.
//
// This integration requires at least holos 0.101.7 to load external data from a
// yaml file.  Kargo will bump the chart version in the yaml file.
Component: #Kubernetes & {
	Resources: {
		let STAGE = "main"

		// The project is the same as the namespace, we adopt the namespace with the
		// kargo.akuity.io/project: "true" label, configured by the namespaces
		// component.
		Project: (CertManager.namespace): spec: promotionPolicies: [{
			stage:                STAGE
			autoPromotionEnabled: true
		}]

		Warehouse: "cert-manager": {
			metadata: name:      "cert-manager"
			metadata: namespace: CertManager.namespace
			spec: {
				// implicit value is Automatic
				freightCreationPolicy: "Automatic"
				// implicit value is 5m0s
				interval: "5m0s"
				subscriptions: [{
					chart: {
						// We leave semverConstraint empty to fetch the latest version
						// because the pipeline submits a pull request that must be manually
						// reviewed and approved.  The purpose is to automate the process of
						// showing the platform engineer what will change.
						name:    CertManager.chart.name
						repoURL: CertManager.chart.repository.url
					}
				}]
			}
		}

		let SRC_PATH = "./src"
		let DATAFILE = path.Join([SRC_PATH, CertManager.datafile], path.Unix)

		Stage: (STAGE): {
			metadata: name:      STAGE
			metadata: namespace: CertManager.namespace
			spec: {
				requestedFreight: [{
					origin: {
						kind: "Warehouse"
						name: Warehouse["cert-manager"].metadata.name
					}
					sources: direct: true
				}]
				promotionTemplate: spec: {
					steps: [
						{
							uses: "git-clone"
							config: {
								repoURL: Organization.RepoURL
								// Unlike the Kargo Quickstart, we aren't promoting into a
								// different branch, we're going to submit a PR to main, so we
								// only need to checkout main.
								checkout: [{
									branch: "main"
									path:   SRC_PATH
								}]
							}
						},
						{
							uses: "yaml-update"
							as:   "update"
							config: {
								path: DATAFILE
								updates: [{
									key: "CertManager.chart.version"
									// https://docs.kargo.io/references/expression-language/#chartfrom
									value: "${{ chartFrom('\(CertManager.chart.repository.url)', '\(CertManager.chart.name)', warehouse('cert-manager')).Version }}"
								}]
							}
						},
						{
							// https://docs.kargo.io/references/promotion-steps#git-commit
							uses: "git-commit"
							as:   "commit"
							config: {
								path:    SRC_PATH
								message: "cert-manager: update to ${{ chartFrom('\(CertManager.chart.repository.url)', '\(CertManager.chart.name)', warehouse('cert-manager')).Version }}"
							}
						},
						{
							// https://docs.kargo.io/references/promotion-steps#git-push
							uses: "git-push"
							as:   "push"
							config: {
								path:                 SRC_PATH
								generateTargetBranch: true
							}
						},
						{
							// https://docs.kargo.io/references/promotion-steps#git-open-pr
							uses: "git-open-pr"
							as:   "pull"
							config: {
								repoURL:      Organization.RepoURL
								sourceBranch: "${{ outputs.push.branch }}"
								targetBranch: "main"
							}
						},
						{
							uses: "git-wait-for-pr"
							as:   "merge-pr"
							config: {
								repoURL:  Organization.RepoURL
								prNumber: "${{ outputs.pull.prNumber }}"
							}
						},
						{
							uses: "argocd-update"
							// Do not update the target revision, let it sync against main.
							config: apps: [{name: "\(ProjectName)-cert-manager"}]
						},
					]
				}
			}
		}

	}
}
