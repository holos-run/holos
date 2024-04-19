package holos

import ap "security.istio.io/authorizationpolicy/v1"

// #AuthPolicyRules represents AuthorizationPolicy rules for hosts that need
// specialized treatment.  Entries in this struct are excluded from
// AuthorizationPolicy/authproxy-custom in the istio-ingress namespace.  Entries
// are added to their own AuthorizationPolicy.
#AuthPolicyRules: {
	// AuthProxySpec represents the identity provider configuration
	AuthProxySpec: #AuthProxySpec & #Platform.authproxy

	// Hosts are hosts that need specialized treatment
	hosts: {
		[Name=_]: {
			// name is the fully qualifed hostname, a Host: header value.
			name: Name
			// slug is the resource name prefix
			slug: string
			// NoAuthorizationPolicy disables an AuthorizationPolicy for the host
			NoAuthorizationPolicy: true | *false

			// Refer to https://istio.io/latest/docs/reference/config/security/authorization-policy/#Rule
			spec: ap.#AuthorizationPolicySpec & {
				action: "CUSTOM"
				provider: name: AuthProxySpec.provider
				selector: matchLabels: istio: "ingressgateway"
			}
		}
	}

	objects: #APIObjects & {
		for Host in hosts {
			if Host.NoAuthorizationPolicy == false {
				apiObjects: {
					AuthorizationPolicy: "\(Host.slug)-custom": {
						metadata: namespace: "istio-ingress"
						metadata: name:      "\(Host.slug)-custom"
						spec: Host.spec
					}
				}
			}
		}
	}
}
