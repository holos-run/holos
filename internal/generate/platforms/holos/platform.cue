package holos

// _Fleets represent the clusters in the platform.
_Fleets: {
	management: clusters: management: _
	workload: clusters: aws1:         _
	workload: clusters: aws2:         _
}

// Namespaces to manage.
_Namespaces: "holos-system": _

for project in _Projects {
	// Include all project namespaces in the platform namespaces.
	_Namespaces: project.spec.namespaces
	// Include project certificates
	_Certificates: project.spec.certificates
}

// Projects to manage.
_Projects: {
	holos: spec: namespaces: "holos-system":                  _
	"external-secrets": spec: namespaces: "external-secrets": _
	istio: spec: namespaces: {
		"istio-system":  _
		"istio-ingress": _
	}
	certificates: spec: namespaces: "cert-manager": _
	argocd: spec: namespaces: argocd:               _
}

// Manage certificates for admin services in workload clusters.
for Cluster in _Fleets.workload.clusters {
	let Name = "argocd.admin.\(Cluster.name).\(_Platform.Model.org.domain)"
	_Projects: argocd: spec: certificates: "\(Name)": #IngressCertificate & {metadata: name: Name}
}

// Platform components to manage.
_Platform: Components: {
	// Components to manage on all clusters.
	for Fleet in _Fleets {
		for Cluster in Fleet.clusters {
			"\(Cluster.name)/namespaces": {
				path:    "components/namespaces"
				cluster: Cluster.name
			}
			"\(Cluster.name)/cert-manager": {
				path:    "components/cert-manager"
				cluster: Cluster.name
			}
		}
	}

	// Components to manage on the management cluster.
	for Cluster in _Fleets.management.clusters {
		"\(Cluster.name)/eso-creds-manager": {
			path:    "components/eso-creds-manager"
			cluster: Cluster.name
		}
		"\(Cluster.name)/cert-letsencrypt": {
			path:    "components/cert-letsencrypt"
			cluster: Cluster.name
		}
		"\(Cluster.name)/certificates": {
			path:    "components/certificates"
			cluster: Cluster.name
		}
	}

	// Components to manage on workload clusters.
	for Cluster in _Fleets.workload.clusters {
		"\(Cluster.name)/external-secrets": {
			path:    "components/external-secrets"
			cluster: Cluster.name
		}
		"\(Cluster.name)/eso-creds-refresher": {
			path:    "components/eso-creds-refresher"
			cluster: Cluster.name
		}
		"\(Cluster.name)/secretstores": {
			path:    "components/secretstores"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-base": {
			path:    "components/istio/base"
			cluster: Cluster.name
		}
		"\(Cluster.name)/argocd": {
			path:    "components/argocd"
			cluster: Cluster.name
		}
	}
}
