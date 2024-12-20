package holos

import "path"

Parameters: {
	KargoProjectName: string           @tag(KargoProjectName)
	KargoStageName:   string | *"main" @tag(KargoStageName)
	// The datafile where the version is stored
	KargoDataFile: string @tag(KargoDataFile)
	// The key in the data file where the version is stored
	KargoDataKey: string @tag(KargoDataKey)
	GitRepoURL:   string @tag(GitRepoURL)
	ChartName:    string @tag(ChartName)
	ChartRepoURL: string @tag(ChartRepoURL)
}

holos: Component.BuildPlan

// Manage a Kargo Project and promotion stages for a cluster add-on.  The use
// case is to watch for new helm chart versions and submit a PR against the main
// branch with the fully rendered manifests.
//
// This integration requires at least holos 0.101.7 to load external data from a
// yaml file.  Kargo will bump the chart version in the yaml file.
Component: #Kubernetes & {
	Resources: {
		// The project is the same as the namespace, we adopt the namespace with the
		// kargo.akuity.io/project: "true" label, configured by the namespaces
		// component.
		Project: (Parameters.KargoProjectName): spec: promotionPolicies: [{
			stage:                Parameters.KargoStageName
			autoPromotionEnabled: true
		}]

		Warehouse: (Parameters.KargoProjectName): {
			metadata: name:      Parameters.KargoProjectName
			metadata: namespace: Parameters.KargoProjectName
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
						name:    Parameters.ChartName
						repoURL: Parameters.ChartRepoURL
					}
				}]
			}
		}

		let SRC_PATH = "./src"
		let DATAFILE = path.Join([SRC_PATH, Parameters.KargoDataFile], path.Unix)

		Stage: (Parameters.KargoStageName): {
			metadata: name:      Parameters.KargoStageName
			metadata: namespace: Parameters.KargoProjectName
			spec: {
				requestedFreight: [{
					origin: {
						kind: "Warehouse"
						name: Warehouse[Parameters.KargoProjectName].metadata.name
					}
					sources: direct: true
				}]
				promotionTemplate: spec: {
					steps: [
						{
							uses: "git-clone"
							config: {
								repoURL: Parameters.GitRepoURL
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
									key: Parameters.KargoDataKey
									// https://docs.kargo.io/references/expression-language/#chartfrom
									value: "${{ chartFrom('\(Parameters.ChartRepoURL)', '\(Parameters.ChartName)', warehouse('\(Parameters.KargoProjectName)')).Version }}"
								}]
							}
						},
						{
							// https://docs.kargo.io/references/promotion-steps#git-commit
							uses: "git-commit"
							as:   "commit"
							config: {
								path:    SRC_PATH
								message: "\(Parameters.KargoProjectName): update to ${{ chartFrom('\(Parameters.ChartRepoURL)', '\(Parameters.ChartName)', warehouse('\(Parameters.KargoProjectName)')).Version }}"
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
								repoURL:      Parameters.GitRepoURL
								sourceBranch: "${{ outputs.push.branch }}"
								targetBranch: "main"
							}
						},
					]
				}
			}
		}

	}
}
