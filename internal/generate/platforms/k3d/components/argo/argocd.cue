package holos

_ArgoCD: {
	metadata: name:      "argocd"
	metadata: namespace: "argocd"

	hostname: "argocd.\(_Platform.Model.org.domain)"

	// issuerHost is the hostname portion of issuerURL
	issuerHost: _AuthProxy.issuerHost
	// issuerURL is the oidc id provider issuer, zitadel for this platform.
	issuerURL: "https://" + issuerHost

	// clientID is the client id of the authproxy in the id provider (zitadel).
	clientID: _Platform.Model.argocd.clientID

	// scopesList represents a list of scopes
	// Omit urn:zitadel:iam:org:domain:primary:example.com scope because members
	// of the Holos and the Open Infrastructure Services orgs may access ArgoCD.
	scopesList: ["openid", "profile", "email", "groups"]
}
