package holos

holos: Component.BuildPlan

Component: #Kustomize & {
	KustomizeConfig: Files: "crds.yaml": _
}
