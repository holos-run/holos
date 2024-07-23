package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "argo-authpolicy"
	Namespace: _AuthProxy.metadata.namespace

	let Selector = {matchLabels: "istio.io/gateway-name": "default"}

	Resources: [_]: [NAME=string]: {
		metadata: _IAP.metadata
		metadata: name:      NAME
		metadata: namespace: Namespace
	}

	// Auth policy resources represent the RequestAuthentication and
	// AuthorizationPolicy resources in the istio-gateways namespace governing the
	// default Gateway.
	Resources: {
		AuthorizationPolicy: "\(Name)-allow-argocd": {
			_description: "Allow argocd access"

			spec: {
				action:   "ALLOW"
				selector: Selector
				rules: [
					{
						to: [{
							// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
							operation: hosts: [
								"argocd.\(_Platform.Model.org.domain)",
								"argocd.\(_Platform.Model.org.domain):*",
							]
						}]
						if _Platform.Model.rbac.mode == "real" {
							when: [
								// Must be issued by the platform identity provider.
								{
									key: "request.auth.principal"
									values: [_AuthProxy.issuerURL + "/*"]
								},
								// Must be intended for an app within the Holos Platform ZITADEL project.
								{
									key: "request.auth.audiences"
									values: [_AuthProxy.projectID]
								},
								// Must be presented by the istio ExtAuthz auth proxy.
								{
									key: "request.auth.presenter"
									values: [_AuthProxy.clientID]
								},
							]

						}
						if _Platform.Model.rbac.mode == "fake" {
							when: [
								{
									// bypass the external authorizer when the user did not sign up
									key: "request.headers[user-agent]"
									values: [_AuthorizedUserAgent]
								},
							]
						}
					},
				]
			}
		}
	}
}
