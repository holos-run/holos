package holos

holos: Component.BuildPlan

Component: #Kustomize & {
	KustomizeConfig: {
		Files: "argocd-crds.\(ArgoCD.Version).yaml": _
		// Resources: {
		// 	"https://raw.githubusercontent.com/argoproj/argo-cd/v\(ArgoCD.Version)/manifests/crds/application-crd.yaml":    _
		// 	"https://raw.githubusercontent.com/argoproj/argo-cd/v\(ArgoCD.Version)/manifests/crds/applicationset-crd.yaml": _
		// 	"https://raw.githubusercontent.com/argoproj/argo-cd/v\(ArgoCD.Version)/manifests/crds/appproject-crd.yaml":     _
		// }
	}
}
