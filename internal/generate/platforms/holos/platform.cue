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
	holos: spec: namespaces: "holos-system": metadata: labels: _Selector.Admin.matchLabels

	"external-secrets": spec: namespaces: "external-secrets": _
	istio: spec: namespaces: {
		"istio-system":   _
		"istio-gateways": _
	}
	certificates: spec: namespaces: "cert-manager": _
	argocd: spec: namespaces: argocd:               _
	login: spec: {
		// Namespace zitadel with grant to attach HTTPRoute resources to the login
		// listeners.
		namespaces: zitadel: metadata: labels: _Selector.Login.matchLabels
		// Cert for *.login.example.com and login.example.com
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
		"\(Cluster.name)/istio-gateways": {
			path:    "components/istio/mesh/gateways"
			cluster: Cluster.name
		}
		"\(Cluster.name)/httpbin": {
			path:    "components/istio/mesh/httpbin"
			cluster: Cluster.name
		}
		"\(Cluster.name)/argocd": {
			path:    "components/argocd"
			cluster: Cluster.name
		}
	}
}

// _Selector represents label selectors
_Selector: {
	// Admin represents the label selector for an admin service. An admin service is
	// defined as a service accessible at a host matching
	// *.admin.<cluster>.<org.domain>  Used by Gateway API to grant HTTPRoute access
	// to Namespaces that contain admin services.
	Admin: matchLabels: "holos.run/admin.grant": "true"
	// Login represents the label selector for zitadel; `*.login.<org.domain> and
	// login.<org.domain>`
	Login: matchLabels: "holos.run/login.grant": "true"
}
