package holos

// Refer to https://kubernetes.io/docs/concepts/security/pod-security-standards/
let Restricted = {
	labels: "pod-security.kubernetes.io/enforce":         "restricted"
	labels: "pod-security.kubernetes.io/enforce-version": "latest"
}
let Privileged = {
	labels: "pod-security.kubernetes.io/enforce":         "privileged"
	labels: "pod-security.kubernetes.io/enforce-version": "latest"
}

// #PlatformNamespaces is the union of all namespaces across all cluster types.  Namespaces are created in all clusters regardless of if they're
// used within the cluster or not.  The is important for security and consistency with IAM, RBAC, and Secrets sync between clusters.
#PlatformNamespaces: [
	{name: "external-secrets"},
	{name: "holos-system"},
	{name: "flux-system"},
	{name: "ceph-system"} & Privileged,
	{name: "istio-system"} & Privileged,
	{name: "istio-ingress"} & Restricted,
	{name: "cert-manager"},
	{name: "argocd"},
]
