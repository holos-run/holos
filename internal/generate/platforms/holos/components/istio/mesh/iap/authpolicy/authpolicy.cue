package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

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
				audiences: ["\(_AuthProxy.projectID)"]
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

		AuthorizationPolicy: "\(Name)-allow-nothing": {
			_description: "Allow nothing"

			spec: {
				action:   "ALLOW"
				selector: Selector
			}
		}

		AuthorizationPolicy: "\(Name)-allow-login": {
			_description: "Allow login"

			spec: {
				action:   "ALLOW"
				selector: Selector
				rules: [
					{
						to: [{
							// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
							operation: hosts: [
								// Allow requests to the login service
								_AuthProxy.issuerHost,
								_AuthProxy.issuerHost + ":*",
							]
						}]
					},
				]
			}
		}

		AuthorizationPolicy: "\(Name)-allow-admin": {
			_description: "Allow cluster admin roles"

			spec: {
				action:   "ALLOW"
				selector: Selector
				rules: [
					{
						to: [{
							// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
							operation: hosts: [
								// Allow authenticated users with cluster admin, edit, or view
								// roles to access admin interfaces.

								// TODO(jeff): The set of admin services should be defined in a
								// nice root-level struct somewhere, probably as part of the
								// _Projects struct.
								"argocd.admin.\(_ClusterName).\(_Platform.Model.org.domain)",
								"argocd.admin.\(_ClusterName).\(_Platform.Model.org.domain):*",
								"httpbin.admin.\(_ClusterName).\(_Platform.Model.org.domain)",
								"httpbin.admin.\(_ClusterName).\(_Platform.Model.org.domain):*",
								"backstage.admin.\(_ClusterName).\(_Platform.Model.org.domain)",
								"backstage.admin.\(_ClusterName).\(_Platform.Model.org.domain):*",
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
							// Must have one of the listed roles.
							AdminRoleGroups,
						]
					},
				]
			}
		}

		AuthorizationPolicy: "\(Name)-allow-holos-server": {
			_description: "Allow authenticated access to holos server"

			spec: {
				action:   "ALLOW"
				selector: Selector
				rules: [
					{
						to: [{
							// Refer to https://istio.io/latest/docs/ops/best-practices/security/#writing-host-match-policies
							operation: hosts: [
								// Allow authenticated users with cluster admin, edit, or view
								// roles to access admin interfaces.

								// TODO(jeff): The set of admin services should be defined in a
								// nice root-level struct somewhere, probably as part of the
								// _Projects struct.
								"app.\(_ClusterName).\(_Platform.Model.org.domain)",
								"app.\(_ClusterName).\(_Platform.Model.org.domain):*",
								"dev.app.\(_ClusterName).\(_Platform.Model.org.domain)",
								"dev.app.\(_ClusterName).\(_Platform.Model.org.domain):*",

								"app.\(_Platform.Model.org.domain)",
								"app.\(_Platform.Model.org.domain):*",
								"dev.app.\(_Platform.Model.org.domain)",
								"dev.app.\(_Platform.Model.org.domain):*",
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
	}
}

let AdminRoleGroups = {
	key: "request.auth.claims[groups]"
	values: [
		"prod-cluster-admin",
		"prod-cluster-edit",
		"prod-cluster-view",
	]
}
