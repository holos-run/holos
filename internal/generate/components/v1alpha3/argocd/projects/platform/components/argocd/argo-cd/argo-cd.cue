package holos

import "strings"

// Produce a helm chart build plan.
(#Helm & Chart).BuildPlan

let Chart = {
	Name:      "argo-cd"
	Namespace: #ArgoCD.Namespace
	Version:   "7.5.2"

	Repo: name: "argocd"
	Repo: url:  "https://argoproj.github.io/argo-helm"

	Chart: chart: name:    "argo-cd"
	Chart: chart: release: "argo-cd"
	// Upstream uses a Kubernetes Job to create the argocd-redis Secret.  Enable
	// hooks to enable the Job.
	Chart: enableHooks: true

	Resources: [_]: [_]: metadata: namespace: Namespace
	// Grant the Gateway namespace the ability to refer to the backend service
	// from HTTPRoute resources.
	Resources: ReferenceGrant: (#Istio.Gateway.Namespace): #ReferenceGrant

	EnableKustomizePostProcessor: true
	// Force all resources into the component namespace, some resources in the
	// helm chart may not specify the namespace so they may get mis-applied
	// depending on the kubectl (client-go) context.
	KustomizeFiles: "kustomization.yaml": namespace: Namespace

	Values: #Values & {
		kubeVersionOverride: "1.29.0"
		// handled in the argo-crds component
		crds: install:  false
		global: domain: "argocd.\(#Platform.Domain)"
		dex: enabled:   false
		// the platform handles mutual tls to the backend
		configs: params: "server.insecure": true

		configs: cm: {
			"admin.enabled":           false
			"oidc.config":             "{}"
			"users.anonymous.enabled": "true"
		}

		// Refer to https://argo-cd.readthedocs.io/en/stable/operator-manual/rbac/
		let Policy = [
			"g, argocd-view, role:readonly",
			"g, prod-cluster-view, role:readonly",
			"g, prod-cluster-edit, role:readonly",
			"g, prod-cluster-admin, role:admin",
		]

		configs: rbac: "policy.csv":     strings.Join(Policy, "\n")
		configs: rbac: "policy.default": "role:admin"
	}
}
