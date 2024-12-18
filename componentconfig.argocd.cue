@if(!NoArgoCD)
package holos

import (
	"path"
	app "argoproj.io/application/v1alpha1"
)

#ComponentConfig: {
	Name:          _
	OutputBaseDir: _
	// Application resources are Cluster scoped.  BuildPlan metadata.name values
	// are Project scoped.  Construct a unique cluster scoped named to resolve
	// conflicts within ArgoCD.
	_ArgoAppName: "\(ProjectName)-\(Name)"

	// Allow other aspects of the platform configuration to refer to
	// `Component._ArgoApplication` to get a handle on the Application resource.
	_ArgoApplication: app.#Application & {
		metadata: name:      _ArgoAppName
		metadata: namespace: "argocd"
		metadata: labels:    Labels
		spec: {
			destination: server: "https://kubernetes.default.svc"
			project: ProjectName
			source: {
				path:           string | *ResourcesPath
				repoURL:        Organization.RepoURL
				targetRevision: string | *"main"
			}
		}
	}

	let ArtifactPath = path.Join([OutputBaseDir, "gitops", "\(Name).application.gen.yaml"], path.Unix)
	let ResourcesPath = path.Join(["deploy", OutputBaseDir, "components", Name], path.Unix)

	// Add the argocd Application instance label to GitOps so resources are in sync.
	KustomizeConfig: CommonLabels: "argocd.argoproj.io/instance": _ArgoAppName

	// Labels for the Application itself.  We filter the argocd application
	// instance label so ArgoCD doesn't think the Application resource manages
	// itself.
	let Labels = {
		for k, v in KustomizeConfig.CommonLabels {
			if k != "argocd.argoproj.io/instance" {
				(k): v
			}
		}
	}

	Artifacts: "\(Name)-application": {
		artifact: ArtifactPath
		generators: [{
			kind:   "Resources"
			output: artifact
			resources: Application: (_ArgoAppName): _ArgoApplication
		}]
	}
}
