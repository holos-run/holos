package holos

holos: Component.BuildPlan

Component: #Kustomize & {
	KustomizeConfig: {
		Resources: {
			"https://raw.githubusercontent.com/argoproj/argo-cd/v\(ArgoCD.Version)/manifests/crds/application-crd.yaml":    _
			"https://raw.githubusercontent.com/argoproj/argo-cd/v\(ArgoCD.Version)/manifests/crds/applicationset-crd.yaml": _
			"https://raw.githubusercontent.com/argoproj/argo-cd/v\(ArgoCD.Version)/manifests/crds/appproject-crd.yaml":     _
		}
	}
}
