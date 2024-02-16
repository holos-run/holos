package holos

import ( "list"

	// This package is intended to be applied only to the provisioner cluster.
	// Resources:
	// - eso-creds-refresher ksa to issue tokens for secret-provisioner ksa in each ns.

	// For each namespace:
	// - Namespaces (e.g. flux-system)
	// - ServiceAccount secret-provisioner
	// - Secret
)

#TargetNamespace: "default"

#InputKeys: {
	cluster:   "provisioner"
	project:   "secrets"
	component: "namespaces"
}

// CredsRefresherName represents the name of the external secrets operator creds refesher job.
let CredsRefresherName = "eso-creds-refresher"

// ProjectID is the google cloud project id containing the provisioner cluster and iam service accounts.
let ProjectID = #InputKeys.gcpProjectID

// IAMServiceAccount is the iam gcp service account the eso creds refresher job uses to access the provisioner api
let IAMServiceAccount = CredsRefresherName + "@" + ProjectID + ".iam.gserviceaccount.com"

// #CredsRefresher defines the rbac policy for the job that refreshes credentials used by eso SecretStore resources in clusters other than the provisioner cluster.
#CredsRefresher: {
	objects: [
		// Allow the IAM Service Account to determine which service account tokens to refresh.
		// List namespaces.
		// List service accounts.
		#ClusterRoleBinding & {
			metadata: name: CredsRefresherName
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "ClusterRole"
				name:     "view"
			}
			subjects: [
				{
					apiGroup: "rbac.authorization.k8s.io"
					kind:     "User"
					name:     IAMServiceAccount
				},
			]
		},
	]
}

// #PlatformNamespaceObjects defines the api objects necessary for eso SecretStores in external clusters to access secrets in a given namespace in the provisioner cluster.
#PlatformNamespaceObjects: {
	_ns: #PlatformNamespace

	objects: [
		#Namespace & {
			metadata: name: _ns.name
		},
		// Role for the eso-creds-refresher
		#Role & {
			metadata: {
				name:      CredsRefresherName
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
				name:      CredsRefresherName
				namespace: _ns.name
			}
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "Role"
				name:     CredsRefresherName
			}
			subjects: [
				{
					apiGroup: "rbac.authorization.k8s.io"
					kind:     "User"
					name:     IAMServiceAccount
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
				name:     CredsRefresherName
			}
			subjects: [
				{
					apiGroup: "rbac.authorization.k8s.io"
					kind:     "User"
					name:     IAMServiceAccount
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

objects: list.FlattenN(_objects, 1)

_objects: [
	(#CredsRefresher & {project:                                      #InputKeys.gcpProjectID}).objects,
	for ns in #PlatformNamespaces {(#PlatformNamespaceObjects & {_ns: ns}).objects},
]

// No flux kustomization
ksObjects: []

{} & #KubernetesObjects
