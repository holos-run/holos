package holos

// _Fleets represent the clusters in the platform.
_Fleets: {
	// The k3d platform has a single cluster named "workload" for simplicity.
	workload: clusters: workload: primary: true
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

	// cert-manager primarily for the management cluster to manage tls certs but
	// also used in workload clusters to generate mTLS certs.
	certificates: spec: namespaces: "cert-manager": _

	// Postgres Operator
	pgo: spec: namespaces: "postgres-operator": _

	// Holos local development
	holosapp: spec: {
		namespaces: "dev-holos": _
		let Subdomain = "app.\(_Platform.Model.org.domain)"
		certificates: "\(Subdomain)": #IngressCertificate
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
			"\(Cluster.name)/local-ca": {
				path:    "components/cert-local-ca"
				cluster: Cluster.name
			}
			"\(Cluster.name)/certificates": {
				path:    "components/certificates"
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

			// Auth Proxy for platform services
			"\(Cluster.name)/authproxy": {
				path:    "components/istio/mesh/iap/authproxy"
				cluster: Cluster.name
			}
			"\(Cluster.name)/authroutes": {
				path:    "components/istio/mesh/iap/authroutes"
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
			"\(Cluster.name)/argo-authpolicy": {
				path:    "components/argo/authpolicy"
				cluster: Cluster.name
			}
			"\(Cluster.name)/argo-routes": {
				path:    "components/argo/routes"
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
	issuerHost: "login.holos.run"
	// issuerURL is the oidc id provider issuer, zitadel for this platform.
	issuerURL: "https://" + issuerHost

	// clientID is the client id of the authproxy in the id provider (zitadel).
	clientID: "275804490387516853@holos_quickstart"
	// projectID is the zitadel project id of Holos Platform project in the id
	// provider (zitadel).
	projectID: "275571128859132936"

	// idTokenHeader is the header where the authproxy places the id token on
	// successful authentication.  Useful for services in the mesh to validate and
	// assert for authorization.
	idTokenHeader: "x-oidc-id-token"

	// provider is the istio meshconfig extauthz provider of the authproxy
	provider: "default-gateway-authproxy"
}
