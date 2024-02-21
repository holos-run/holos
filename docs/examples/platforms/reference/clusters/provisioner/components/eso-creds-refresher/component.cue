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
	#CredsRefresherIAM.binding,
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
	// Allow the IAM Service Account to determine which service account tokens to refresh.
	// List namespaces.
	// List service accounts.
	binding: #ClusterRoleBinding & {
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
	}
}

// #PlatformNamespaceObjects defines the api objects necessary for eso SecretStores in external clusters to access secrets in a given namespace in the provisioner cluster.
#PlatformNamespaceObjects: {
	_ns: #PlatformNamespace

	let iamsa = #IAMServiceAccount & {
		namespace: _ns.name
	}
	let reader = #ESOServiceAccount & {
		namespace: _ns.name
		mode:      "reader"
	}
	let writer = #ESOServiceAccount & {
		namespace: _ns.name
		mode:      "writer"
	}

	objects: [
		iamsa.role,
		iamsa.binding,
		reader.ksa,
		reader.role,
		reader.binding,
		writer.ksa,
		writer.role,
		writer.binding,
	]
}

// #IAMServiceAccount represents the iam service account used by the eso-creds-refresher Job to get credentials for #ESOServiceAccounts
#IAMServiceAccount: {
	namespace: string
	let _namespace = namespace

	// Role for the eso-creds-refresher
	role: #Role & {
		metadata: {
			name:      #CredsRefresher.name
			namespace: _namespace
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
	}
	binding: #RoleBinding & {
		metadata: {
			name:      #CredsRefresher.name
			namespace: _namespace
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
	}
}

// #ESOServiceAccount represents the ksa used by SecretStores in non-provisioner clusters to sync secrets with the provisioner cluster.
#ESOServiceAccount: {
	namespace: string
	mode:      string | "reader" | "writer"

	let _name = "eso-" + mode
	let _namespace = namespace

	// ksa is the ServiceAccount used by the SecretStore
	ksa: #ServiceAccount & {
		metadata: {
			name:      _name
			namespace: _namespace
			annotations: "holos.run/description": "Used by SecretStores in other clusters to sync secrets with the provisioner cluster."
		}
	}

	// role is the Role of the ServiceAccount used by the SecretStore
	role: #Role & {
		metadata: {
			name:      _name
			namespace: _namespace
		}
		rules: [
			{
				apiGroups: [""]
				resources: ["secrets"]
				if mode == "reader" {
					verbs: ["get", "list", "watch"]
				}
				if mode == "writer" {
					verbs: ["*"]
				}
			},
		]
	}
	binding: #RoleBinding & {
		metadata: {
			name:      _name
			namespace: _namespace
		}
		roleRef: {
			apiGroup: "rbac.authorization.k8s.io"
			kind:     "Role"
			name:     _name
		}
		subjects: [
			{
				apiGroup:  "rbac.authorization.k8s.io"
				kind:      "ServiceAccount"
				name:      _name
				namespace: _namespace
			},
		]
	}
}
