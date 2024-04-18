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
			dex: enabled:   false
			// for integration with istio
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
			spec: gateways: ["istio-ingress/default"]
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

let OAuthClient = #Platform.oauthClients.argocd.spec

let OIDCConfig = {
	name:            "Holos Platform"
	issuer:          OAuthClient.issuer
	clientID:        OAuthClient.clientID
	requestedScopes: OAuthClient.scopesList
	// Set redirect uri to https://argocd.example.com/pkce/verify
	enablePKCEAuthentication: true

	requestedIDTokenClaims: groups: essential: true
}
