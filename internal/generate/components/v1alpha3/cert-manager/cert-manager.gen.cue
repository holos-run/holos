package holos

// Platform wide configuration
#CertManager: {
	Version:   "{{ .Version }}"
	Namespace: "{{ .Namespace }}"
}

// Register the namespace
#Namespaces: (#CertManager.Namespace): _

// Manage the component on every cluster in the platform
for Fleet in #Fleets {
	for Cluster in Fleet.clusters {
		#Platform: Components: "\(Cluster.name)/{{ .Name }}": {
			path:    "components/cert-manager"
			cluster: Cluster.name
		}
	}
}
