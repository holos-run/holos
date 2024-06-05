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
	holos: spec: namespaces: "holos-system": _
	argocd: spec: namespaces: argocd:        _

	// Sync secrets from the management cluster to workload clusters.
	"external-secrets": spec: namespaces: "external-secrets": _

	// Istio service mesh
	istio: spec: namespaces: {
		"istio-system": _
		"istio-gateways": {
			// All subdomains should be included here for the authproxy HTTPRoute to
			// attach to all listeners.
			metadata: labels: _Selector.GrantSubdomainAdmin.matchLabels
			metadata: labels: _Selector.GrantSubdomainLogin.matchLabels
			metadata: labels: _Selector.GrantSubdomainApp.matchLabels
		}
	}

	// cert-manager primarily for the management cluster but also used in workload
	// clusters to provide cluster scoped ca cert.
	certificates: spec: namespaces: "cert-manager": _

	// Postgres Operator
	pgo: spec: namespaces: "postgres-operator": _

	login: spec: {
		// Namespace for zitadel.
		namespaces: zitadel: _
		// Certificate for login.example.com and *.login.example.com
		let Subdomain = "login.\(_Platform.Model.org.domain)"
		certificates: "\(Subdomain)": #IngressCertificate
		certificates: "any.\(Subdomain)": #IngressCertificate & {
			spec: commonName: "*." + Subdomain
		}
	}

	holosapp: spec: {
		namespaces: "dev-holos":  _
		namespaces: "prod-holos": _
		namespaces: "jeff-holos": _

		let Subdomain = "app.\(_Platform.Model.org.domain)"
		certificates: "\(Subdomain)": #IngressCertificate
		certificates: "any.\(Subdomain)": #IngressCertificate & {
			spec: commonName: "*." + Subdomain
		}
	}
}

// Manage per-cluster certificates for services in workload clusters.
for Cluster in _Fleets.workload.clusters {
	// Issue a wildcard cert for all admin interfaces.   We need to verify this is
	// well-behaved with Istio and HTTP2.
	let CertPair = #ClusterCertPair & {cluster: Cluster.name}
	_Projects: holos: spec: (CertPair & {name: "admin"}).spec

	// Holos app certs
	_Projects: holosapp: spec: (CertPair & {name: "app"}).spec
}

#ClusterCertPair: {
	name:    string
	cluster: string
	let Subdomain = name + ".\(cluster).\(_Platform.Model.org.domain)"
	spec: certificates: (Subdomain): #IngressCertificate
	spec: certificates: "any.\(Subdomain)": #IngressCertificate & {
		spec: commonName: "*." + Subdomain
	}
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
		// ECR Credentials (ecr-creds-<account-number>)
		"\(Cluster.name)/ecr-creds-manager": {
			path:    "components/ecr-creds-manager"
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
		// Secret ecr-creds-<aws-account-number> in each Namespace to pull images
		// from the private ECR registry.
		"\(Cluster.name)/ecr-creds-refresher": {
			path:    "components/ecr-creds-refresher"
			cluster: Cluster.name
		}
		// We use HTTPRoute from the Kubernetes Gateway API v1 instead of
		// VirtualService from the Istio Gateway API.
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
		"\(Cluster.name)/httpbin-backend": {
			path:    "components/istio/mesh/httpbin/backend"
			cluster: Cluster.name
		}
		"\(Cluster.name)/httpbin-routes": {
			path:    "components/istio/mesh/httpbin/routes"
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
		"\(Cluster.name)/zitadel-routes": {
			path:    "components/login/zitadel-routes"
			cluster: Cluster.name
		}
		// Auth Proxy for platform services
		"\(Cluster.name)/authproxy": {
			path:    "components/istio/mesh/iap/authproxy"
			cluster: Cluster.name
		}
		"\(Cluster.name)/authpolicy": {
			path:    "components/istio/mesh/iap/authpolicy"
			cluster: Cluster.name
		}
		// ArgoCD components
		"\(Cluster.name)/argo-crds": {
			path:    "components/argo/crds"
			cluster: Cluster.name
		}
		"\(Cluster.name)/argo-cd": {
			path:    "components/argo/cd"
			cluster: Cluster.name
		}
		"\(Cluster.name)/argo-routes": {
			path:    "components/argo/routes"
			cluster: Cluster.name
		}
		"\(Cluster.name)/argo-creds": {
			path:    "components/argo/creds"
			cluster: Cluster.name
		}

		// Holos server
		"\(Cluster.name)/apps/dev-holos-infra": {
			path:    "apps/dev/holos/infra"
			cluster: Cluster.name
		}
		"\(Cluster.name)/apps/dev-holos-app": {
			path:    "apps/dev/holos/app"
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

// _AuthProxy represents the authproxy service.  The Service name is referenced
// in multiple components and namespaces.
_AuthProxy: {
	metadata: name:      "authproxy"
	metadata: namespace: #IstioGatewaysNamespace

	// pathPrefix represents the path for all hostnames where the authproxy and
	// authorization policy redirects to start the login flow.
	pathPrefix: "/holos/authproxy"
	// servicePort is the port oauth2-proxy listens on and the Service is
	// reachable at.
	servicePort: 4180

	// issuerHost is the hostname portion of issuerURL
	issuerHost: "login." + _Platform.Model.org.domain
	// issuerURL is the oidc id provider issuer, zitadel for this platform.
	issuerURL: "https://" + issuerHost

	// clientID is the client id of the authproxy in the id provider (zitadel).
	clientID: _Platform.Model.authproxy.clientID
	// projectID is the zitadel project id of Holos Platform project in the id
	// provider (zitadel).
	projectID: _Platform.Model.authproxy.projectID

	// idTokenHeader is the header where the authproxy places the id token on
	// successful authentication.  Useful for services in the mesh to validate and
	// assert for authorization.
	idTokenHeader: "x-oidc-id-token"

	// provider is the istio meshconfig extauthz provider of the authproxy
	provider: "default-gateway-authproxy"
}
