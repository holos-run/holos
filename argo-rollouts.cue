@if(!NoArgoRollouts && !NoArgoCD)
package holos

ArgoRollouts: {
	Namespace: "argo-rollouts"
}

// Register the ArgoCD Project namespaces and components
Projects: {
	argocd: {
		namespaces: (ArgoRollouts.Namespace): _
		_components: {
			"rollouts-crds": _
			"rollouts":      _
		}
	}
}
