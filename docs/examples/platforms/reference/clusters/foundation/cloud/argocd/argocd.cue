package holos

import "encoding/yaml"

let ArgoCD = "argocd"
let Namespace = "prod-platform"

spec: components: HelmChartList: [
	#HelmChart & {
		_dependsOn: "prod-secrets-stores": _
		namespace: Namespace

		metadata: name: "\(namespace)-\(ArgoCD)"

		chart: {
			name:    "argo-cd"
			release: "argocd"
			version: "6.7.8"
			repository: {
				name: "argocd"
				url:  "https://argoproj.github.io/argo-helm"
			}
		}

		_values: #ArgoCDValues & {
			kubeVersionOverride: "1.29.0"
			global: domain: "argocd.\(#ClusterName).\(#Platform.org.domain)"

			configs: params: "server.insecure": true

			configs: cm: {
				"admin.enabled": false
				"oidc.config":   yaml.Marshal(OIDCConfig)
			}
		}

		// Holos overlay objects
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		// ExternalSecret: "deploy-key": _
		VirtualService: (ArgoCD): {
			metadata: name:      ArgoCD
			metadata: namespace: Namespace
			spec: hosts: [
				ArgoCD + ".\(#Platform.org.domain)",
				ArgoCD + ".\(#ClusterName).\(#Platform.org.domain)",
			]
			spec: gateways: ["istio-ingress/\(Namespace)"]
			spec: http: [{route: [{destination: {
				host: "argocd-server.\(Namespace).svc.cluster.local"
				port: number: 80
			}}]}]
		}
	}
}
let IstioInject = [{op: "add", path: "/spec/template/metadata/labels/sidecar.istio.io~1inject", value: "true"}]

#Kustomize: _patches: {
	mesh: {
		target: {
			group:   "apps"
			version: "v1"
			kind:    "Deployment"
			name:    "argocd-server"
		}
		patch: yaml.Marshal(IstioInject)
	}
}

// Probably shouldn't use the authproxy struct and should instead define an identity provider struct.
let AuthProxySpec = #AuthProxySpec & #Platform.authproxy

let OIDCConfig = {
	name:     "Holos Platform"
	issuer:   AuthProxySpec.issuer
	clientID: #Platform.argocd.clientID
	requestedIDTokenClaims: groups: essential: true
	requestedScopes: ["openid", "profile", "email", "groups", "urn:zitadel:iam:org:domain:primary:\(AuthProxySpec.orgDomain)"]
	enablePKCEAuthentication: true
}
