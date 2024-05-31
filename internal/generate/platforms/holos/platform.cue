package holos

// _Fleets represent the clusters in the platform.
_Fleets: {
	management: clusters: management: _
	workload: clusters: aws1:         _
	workload: clusters: aws2: primary: true
}

// Map all clusters in all fleets into the _Clusters struct.
for Fleet in _Fleets {
	for Cluster in Fleet.clusters {
		_Clusters: "\(Cluster.name)": Cluster
	}
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
	// Admin projects accessible at *.admin.<cluster>.<org.domain>
	holos: spec: namespaces: "holos-system": metadata: labels: _Selector.GrantSubdomainAdmin.matchLabels
	argocd: spec: namespaces: argocd: metadata: labels:        _Selector.GrantSubdomainAdmin.matchLabels

	// Sync secrets from the management cluster to workload clusters.
	"external-secrets": spec: namespaces: "external-secrets": _

	// Istio service mesh
	istio: spec: namespaces: {
		"istio-system":   _
		"istio-gateways": _
	}

	// cert-manager primarily for the management cluster but also used in workload
	// clusters to provide cluster scoped ca cert.
	certificates: spec: namespaces: "cert-manager": _

	// Postgres Operator
	pgo: spec: namespaces: "postgres-operator": _

	login: spec: {
		// Namespace for zitadel.
		namespaces: zitadel: metadata: labels: _Selector.GrantSubdomainLogin.matchLabels
		// Certificate for login.example.com and *.login.example.com
		certificates: "login.\(_Platform.Model.org.domain)": #IngressCertificate & {
			metadata: name: string
			spec: dnsNames: [metadata.name, "*." + metadata.name]
		}
	}
}

// Manage certificates for admin services in workload clusters.
for Cluster in _Fleets.workload.clusters {
	// Issue a wildcard cert for all admin interfaces.   We need to verify this is
	// well-behaved with Istio and HTTP2.
	let Admin = "admin.\(Cluster.name).\(_Platform.Model.org.domain)"
	_Projects: holos: spec: certificates: "\(Admin)": #IngressCertificate & {
		metadata: name: Admin
		spec: dnsNames: [Admin, "*.\(Admin)"]
	}

	// Issue a dedicated cert for argocd.  This may be removed if the wildcard
	// works with the Gateway API.
	let Name = "argocd.\(Admin)"
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
		// Provision a root ca and certs on the management cluster for the zitadel
		// database in a workload cluster.
		"\(Cluster.name)/zitadel-certs": {
			path:    "components/login/zitadel-certs"
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
		"\(Cluster.name)/gateway-api": {
			path:    "components/gateway-api"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-base": {
			path:    "components/istio/base"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istio-cni": {
			path:    "components/istio/mesh/cni"
			cluster: Cluster.name
		}
		"\(Cluster.name)/istiod": {
			path:    "components/istio/mesh/istiod"
			cluster: Cluster.name
		}
		"\(Cluster.name)/gateway": {
			path:    "components/istio/mesh/gateway"
			cluster: Cluster.name
		}
		"\(Cluster.name)/httpbin": {
			path:    "components/istio/mesh/httpbin"
			cluster: Cluster.name
		}
		"\(Cluster.name)/postgres-crds": {
			path:    "components/pgo/crds"
			cluster: Cluster.name
		}
		"\(Cluster.name)/postgres-operator": {
			path:    "components/pgo/controller"
			cluster: Cluster.name
		}
		// ZITADEL components
		"\(Cluster.name)/zitadel-secrets": {
			path:    "components/login/zitadel-secrets"
			cluster: Cluster.name
		}
		"\(Cluster.name)/zitadel-database": {
			path:    "components/login/zitadel-database"
			cluster: Cluster.name
		}
		"\(Cluster.name)/zitadel-server": {
			path:    "components/login/zitadel-server"
			cluster: Cluster.name
		}
		// ArgoCD components
		"\(Cluster.name)/argocd": {
			path:    "components/argocd"
			cluster: Cluster.name
		}
	}
}

// _Selector represents label selectors
_Selector: {
	// GrantSubdomainAdmin represents the label selector to grant HTTPRoute
	// attachment for the admin subdomain. An admin service is defined as a
	// service accessible at a host matching *.admin.<cluster>.<org.domain>  Used
	// by Gateway API to grant HTTPRoute access to Namespaces that contain admin
	// services.
	GrantSubdomainAdmin: matchLabels: "grant.holos.run/subdomain.admin": "true"
	// GrantSubdomainLogin represents the label selector to grant HTTPRoute
	// attachment for the login subdomain; `*.login.<org.domain> and
	// login.<org.domain>`
	GrantSubdomainLogin: matchLabels: "grant.holos.run/subdomain.login": "true"
	// GrantSubdomainApp represents the label selector to grant HTTPRoute
	// attachment for the app subdomain; `*.login.<org.domain> and
	// login.<org.domain>`
	GrantSubdomainApp: matchLabels: "grant.holos.run/subdomain.app": "true"
}
