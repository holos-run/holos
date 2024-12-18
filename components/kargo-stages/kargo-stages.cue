@if(!NoKargo)

package holos

// _image represents the image to promote through the stages.  For example:
// "us-central1-docker.pkg.dev/bank-of-anthos-ci/bank-of-anthos/frontend"
_image:            string @tag(image)
_semverConstraint: string @tag(semverConstraint)
_ProjectName:      string @tag(ProjectName)

holos: Component.BuildPlan

// This component represents a basic, starter promotion process for a Holos
// Project.  Project owners are expected to use this as a starting point, copy
// the component into their own project definition and compose the project in
// accordingly.
Component: #Kubernetes & {
	Resources: {
		// Place all of the resources in the Kargo Project namespace.
		[_]: [_]: metadata: namespace: _ProjectName

		Warehouse: (_ProjectName): {
			spec: {
				subscriptions: [{
					image: {
						repoURL:          _image
						semverConstraint: _semverConstraint
						discoveryLimit:   5
					}
				}]
			}
		}

		for STAGE in KargoProjects[_ProjectName].stages {
			// NOTE: This assumes a simple structure where there is one component
			// named the same as the project which is managed in each stage.  This
			// hols true for podinfo.  There is a component named podinfo managed in
			// each stage with the stage name given as a prefix.  We promote through
			// those stages.
			//
			// In practice, we'd have a well defined structure of which components in
			// the project have promotion pipelines and build them from that
			// structure.  This structure would vary over on the promotable
			// components, images, and stages within the Project, likely defined as a
			// field of the #KargoProject definition.
			let ComponentName = "\(STAGE.name)-\(_ProjectName)"
			let OutPath = "deploy/projects/\(_ProjectName)/components/\(ComponentName)"
			let BRANCH = "project/\(_ProjectName)/component/\(ComponentName)"

			Stage: (ComponentName): {
				spec: {
					// The requested freight is a static structure where users actually
					// define how artifacts are promoted with kargo.  Currently this is a
					// static mapping for nonprod stages and all prod tier stages promote
					// from the uat stage.
					//
					// See projects.podinfo.cue for an example of how this is configured.
					requestedFreight: KargoProjects[_ProjectName].promotions[ComponentName].requestedFreight
					promotionTemplate: spec: {
						let SRC = "./src"
						let OUT = "./out"
						steps: [
							{
								uses: "git-clone"
								config: {
									repoURL: Organization.RepoURL
									checkout: [
										{
											branch: "main"
											path:   SRC
										},
										{
											// ComponentName has the stage prefix, so no need to also
											// scope to the stage name.
											branch: BRANCH
											path:   OUT
										},
									]
								}
							},
							{
								uses: "git-clear"
								config: path: OUT
							},
							{
								uses: "kustomize-set-image"
								as:   "update-image"
								config: {
									path: "\(SRC)/\(OutPath)"
									images: [{image: _image}]
								}
							},
							{
								uses: "kustomize-build"
								config: {
									path:    "\(SRC)/\(OutPath)"
									outPath: "\(OUT)/\(ComponentName).gen.yaml"
								}
							},
							{
								uses: "git-commit"
								as:   "commit"
								config: {
									path: OUT
									messageFromSteps: ["update-image"]
								}
							},
							{
								uses: "git-push"
								config: {
									path:         OUT
									targetBranch: BRANCH
								}
							},
							{
								uses: "argocd-update"
								config: {
									apps: [{
										name: "\(ProjectName)-\(ComponentName)"
										sources: [{
											repoURL:               Organization.RepoURL
											desiredCommitFromStep: "commit"
										}]
									}]
								}
							},
						]
					}
				}
			}
		}
	}
}
