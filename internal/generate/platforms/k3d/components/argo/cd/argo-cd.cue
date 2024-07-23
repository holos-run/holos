package holos

import (
	"encoding/yaml"
	"strings"
)

// Produce a helm chart build plan.
(#Helm & Chart).Output

let Chart = {
	Name:      "argo-cd"
	Namespace: _ArgoCD.metadata.namespace
	Version:   "7.1.1"

	Chart: chart: release: _ArgoCD.metadata.name
	// Upstream uses a Kubernetes Job to create the argocd-redis Secret.  Enable
	// hooks to enable the Job.
	Chart: enableHooks: true

	Repo: name: "argocd"
	Repo: url:  "https://argoproj.github.io/argo-helm"

	Resources: [_]: [_]: metadata: namespace: Namespace
	// Grant the Gateway namespace the ability to refer to the backend service
	// from HTTPRoute resources.
	Resources: ReferenceGrant: (#IstioGatewaysNamespace): #ReferenceGrant

	EnableKustomizePostProcessor: true
	// Force all resources into the component namespace, some resources in the
	// helm chart may not specify the namespace so they may get mis-applied
	// depending on the kubectl (client-go) context.
	KustomizeFiles: "kustomization.yaml": namespace: Namespace

	// Patch the backend with the service mesh sidecar.
	KustomizePatches: {
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

	Values: #Values & {
		kubeVersionOverride: "1.29.0"
		// handled in the argo-crds component
		crds: install:  false
		global: domain: _ArgoCD.hostname
		dex: enabled:   false
		// for integration with istio
		configs: params: "server.insecure": true
		configs: cm: {
			"admin.enabled": false
			if _Platform.Model.rbac.mode == "real" {
				"oidc.config":             yaml.Marshal(OIDCConfig)
				"users.anonymous.enabled": "false"
			}
			if _Platform.Model.rbac.mode == "fake" {
				"oidc.config":             "{}"
				"users.anonymous.enabled": "true"
			}
		}

		// Refer to https://argo-cd.readthedocs.io/en/stable/operator-manual/rbac/
		let Policy = [
			"g, argocd-view, role:readonly",
			"g, prod-cluster-view, role:readonly",
			"g, prod-cluster-edit, role:readonly",
			"g, prod-cluster-admin, role:admin",
			"g, \(_Platform.Model.rbac.sub), role:admin",
		]

		configs: rbac: "policy.csv": strings.Join(Policy, "\n")

		if _Platform.Model.rbac.mode == "fake" {
			configs: rbac: "policy.default": "role:admin"
		}
	}
}

let IstioInject = [{op: "add", path: "/spec/template/metadata/labels/sidecar.istio.io~1inject", value: "true"}]

let OIDCConfig = {
	name:            "Holos Platform"
	issuer:          _ArgoCD.issuerURL
	clientID:        _ArgoCD.clientID
	requestedScopes: _ArgoCD.scopesList
	// Set redirect uri to https://argocd.example.com/pkce/verify
	enablePKCEAuthentication: true
	// groups is essential for rbac
	requestedIDTokenClaims: groups: essential: true
}
