package holos

// #ArgoCD represents platform wide configuration
#ArgoCD: {
	Version:   "2.12.3"
	Namespace: "argocd"
}

// Register namespaces
#Namespaces: (#ArgoCD.Namespace): _

// Register the HTTPRoute to the backend Service
#HTTPRoutes: argocd: _backendRefs: "argocd-server": namespace: #ArgoCD.Namespace
