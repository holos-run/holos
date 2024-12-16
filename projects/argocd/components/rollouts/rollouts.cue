package holos

holos: Component.BuildPlan

Component: #Kustomize & {
	KustomizeConfig: Kustomization: namespace: ArgoRollouts.Namespace
	KustomizeConfig: Files: "rollouts.yaml":   _
}
