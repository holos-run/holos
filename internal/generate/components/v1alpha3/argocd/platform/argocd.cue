package holos

// Manage the Component on every Cluster in the Platform
for Fleet in #Fleets {
	for Cluster in Fleet.clusters {
		#Platform: Components: "\(Cluster.name)/argocd-crds": {
			path:    "projects/platform/components/argocd/crds"
			cluster: Cluster.name
		}
		#Platform: Components: "\(Cluster.name)/argocd-backend": {
			path:    "projects/platform/components/argocd/backend"
			cluster: Cluster.name
		}
	}
}
