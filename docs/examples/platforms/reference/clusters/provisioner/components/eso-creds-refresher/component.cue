package holos

// This component is intended to be applied only to the provisioner cluster.
// Resources:
// - RBAC for User eso-creds-refresher IAM Service Account to issue tokens for eso-reader and eso-writer in each namespace.

// For each namespace:
// - Namespace
// - ServiceAccount eso-reader, eso-writer

import "list"

// objects are kubernetes api objects to apply.
objects: list.FlattenN(_objects, 1)

_objects: [
	#CredsRefresherIAM.objects,
	for ns in #PlatformNamespaces {(#PlatformNamespaceObjects & {_ns: ns}).objects},
]

// No flux kustomization
ksObjects: []

{} & #KubernetesObjects

#InputKeys: {
	cluster:   "provisioner"
	project:   "secrets"
	component: "namespaces"
}

// #CredsRefresherIAM defines the rbac policy for the job that refreshes credentials used by eso SecretStore resources in clusters other than the provisioner cluster.
#CredsRefresherIAM: {
	objects: [
		// Allow the IAM Service Account to determine which service account tokens to refresh.
		// List namespaces.
		// List service accounts.
		#ClusterRoleBinding & {
			metadata: name: #CredsRefresher.name
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "ClusterRole"
				name:     "view"
			}
			subjects: [
				{
					apiGroup: "rbac.authorization.k8s.io"
					kind:     "User"
					name:     #CredsRefresher.iamServiceAccount
				},
			]
		},
	]
}

// #PlatformNamespaceObjects defines the api objects necessary for eso SecretStores in external clusters to access secrets in a given namespace in the provisioner cluster.
#PlatformNamespaceObjects: {
	_ns: #PlatformNamespace

	objects: [
		// Role for the eso-creds-refresher
		#Role & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: _ns.name
			}
			rules: [
				{
					apiGroups: [""]
					resources: ["serviceaccounts/token"]
					verbs: ["create"]
				},
				{
					apiGroups: [""]
					resources: ["serviceaccounts"]
					verbs: ["list"]
				},
			]
		},
		#RoleBinding & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: _ns.name
			}
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "Role"
				name:     #CredsRefresher.name
			}
			subjects: [
				{
					apiGroup: "rbac.authorization.k8s.io"
					kind:     "User"
					name:     #CredsRefresher.iamServiceAccount
				},
			]
		},
		// Used by SecretStores in non-provisioner clusters to fetch secrets from the provisioner cluster.
		#ServiceAccount & {
			metadata: {
				name:      "eso-reader"
				namespace: _ns.name
				annotations: "holos.run/description": "Used by SecretStore to pull ExternalSecrets"
			}
		},
		#Role & {
			metadata: {
				name:      "eso-reader"
				namespace: _ns.name
			}
			rules: [
				{
					apiGroups: [""]
					resources: ["secrets"]
					verbs: ["get", "list", "watch"]
				},
			]
		},
		#RoleBinding & {
			metadata: {
				name:      "eso-reader"
				namespace: _ns.name
			}
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "Role"
				name:     #CredsRefresher.name
			}
			subjects: [
				{
					apiGroup: "rbac.authorization.k8s.io"
					kind:     "User"
					name:     #CredsRefresher.iamServiceAccount
				},
			]
		},
		// Used by SecretStores in non-provisioner clusters to push secrets to the provisioner cluster.
		#ServiceAccount & {
			metadata: {
				name:      "eso-writer"
				namespace: _ns.name
				annotations: "holos.run/description": "Used by SecretStore with PushSecrets"
			}
		},
		#Role & {
			metadata: {
				name:      "eso-writer"
				namespace: _ns.name
			}
			rules: [
				{
					apiGroups: [""]
					resources: ["secrets"]
					verbs: ["*"]
				},
			]
		},
	]
}
