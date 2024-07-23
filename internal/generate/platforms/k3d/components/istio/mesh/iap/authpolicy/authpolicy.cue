package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Mode = _Platform.Model.rbac.mode

let Objects = {
	Name:      "authpolicy"
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
		RequestAuthentication: (Name): {
			spec: jwtRules: [{
				audiences: [
					"\(_AuthProxy.projectID)",
					_HolosCLIClientID,
				]
				forwardOriginalToken: true
				fromHeaders: [{name: _AuthProxy.idTokenHeader}]
				issuer: _AuthProxy.issuerURL
			}]
			spec: selector: Selector
		}

		AuthorizationPolicy: "\(Name)-custom": {
			_description: "Route all requests through the auth proxy by default"

			spec: {
				action: "CUSTOM"
				provider: name: _AuthProxy.provider
				rules: [
					{
						to: [
							{
								// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
								operation: notHosts: [
									// Never send requests for the login service through the
									// authorizer, would block login.
									_AuthProxy.issuerHost,
									_AuthProxy.issuerHost + ":*",
								]
								operation: notPaths: [
									"/holos/authproxy/{**}",
								]
							},
						]
						if Mode == "real" {
							when: [
								{
									// bypass the external authorizer when the id token is already in the request.
									// the RequestAuthentication rule will verify the token.
									key: "request.headers[\(_AuthProxy.idTokenHeader)]"
									notValues: ["*"]
								},
							]
						}
						if Mode == "fake" {
							when: [
								{
									// bypass the external authorizer when the user did not sign up
									key: "request.headers[user-agent]"
									notValues: [_AuthorizedUserAgent]
								},
							]
						}
					},
				]
				selector: Selector
			}
		}

		AuthorizationPolicy: "\(Name)-allow-nothing": {
			_description: "Allow nothing"

			spec: {
				action:   "ALLOW"
				selector: Selector
			}
		}

		if Mode == "real" {
			AuthorizationPolicy: "\(Name)-allow-backstage": {
				_description: "Allow backstage access"

				spec: {
					action:   "ALLOW"
					selector: Selector
					rules: [
						{
							to: [{
								// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
								operation: hosts: [
									"backstage.\(_Platform.Model.org.domain)",
									"backstage.\(_Platform.Model.org.domain):*",
								]
							}]
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
						},
					]
				}
			}

			AuthorizationPolicy: "\(Name)-allow-httpbin": {
				_description: "Allow httpbin authenticated access"

				spec: {
					action:   "ALLOW"
					selector: Selector
					rules: [
						{
							to: [{
								// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
								operation: hosts: [
									"httpbin.\(_Platform.Model.org.domain)",
									"httpbin.\(_Platform.Model.org.domain):*",
								]
							}]
							when: [
								// Token must be issued by the platform identity provider.
								{
									key: "request.auth.principal"
									values: [_AuthProxy.issuerURL + "/*"]
								},
								// Token must be intended for one of these audiences
								{
									key: "request.auth.audiences"
									values: [_HolosCLIClientID, _AuthProxy.projectID]
								},
							]
						},
					]
				}
			}

			AuthorizationPolicy: "\(Name)-allow-authproxy": {
				_description: "Allow userinfo unauthenicated access"

				spec: {
					action:   "ALLOW"
					selector: Selector
					// Refer to https://oauth2-proxy.github.io/oauth2-proxy/features/endpoints/
					rules: [{to: [{operation: paths: [
						"/holos/authproxy/sign_in",
						"/holos/authproxy/sign_out",
						"/holos/authproxy/start",
						"/holos/authproxy/callback",
						"/holos/authproxy/userinfo",
						"/holos/authproxy/auth",
						"/holos/authproxy/static/{**}",
					]}]}]
				}
			}

		}

		if Mode == "fake" {
			AuthorizationPolicy: "\(Name)-allow-httpbin": {
				_description: "Allow httpbin authenticated access"

				spec: {
					action:   "ALLOW"
					selector: Selector
					rules: [
						{
							to: [{
								// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
								operation: hosts: [
									"httpbin.\(_Platform.Model.org.domain)",
									"httpbin.\(_Platform.Model.org.domain):*",
								]
							}]
							when: [
								{
									// bypass the external authorizer when the user did not sign up
									key: "request.headers[user-agent]"
									values: [_AuthorizedUserAgent]
								},
							]
						},
					]
				}
			}
		}
	}
}
