package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "authpolicy"
	Namespace: _AuthProxy.metadata.namespace

	let Metadata = _IAP.metadata
	let Selector = {matchLabels: "istio.io/gateway-name": "default"}

	Resources: [_]: [_]: metadata: namespace: Namespace

	// Auth policy resources represent the RequestAuthentication and
	// AuthorizationPolicy resources in the istio-gateways namespace governing the
	// default Gateway.
	Resources: {
		RequestAuthentication: (Name): {
			metadata: Metadata & {name: Name}
			spec: jwtRules: [{
				audiences: ["\(_AuthProxy.projectID)"]
				forwardOriginalToken: true
				fromHeaders: [{name: _AuthProxy.idTokenHeader}]
				issuer: _AuthProxy.issuerURL
			}]
			spec: selector: Selector
		}

		AuthorizationPolicy: "\(Name)-custom": {
			_description: "Route all requests through the auth proxy by default"

			metadata: Metadata & {name: "\(Name)-custom"}
			spec: {
				action: "CUSTOM"
				provider: name: _AuthProxy.provider
				rules: [
					{
						to: [{
							// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
							operation: notHosts: [
								// Never send requests for the login service through the
								// authorizer, would block login.
								_AuthProxy.issuerHost,
								_AuthProxy.issuerHost + ":*",
							]
						}]
						when: [
							{
								// bypass the external authorizer when the id token is already in the request.
								// the RequestAuthentication rule will verify the token.
								key: "request.headers[\(_AuthProxy.idTokenHeader)]"
								notValues: ["*"]
							},
						]
					},
				]
				selector: Selector
			}
		}
	}
}
