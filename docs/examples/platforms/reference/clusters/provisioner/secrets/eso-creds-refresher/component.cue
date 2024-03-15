package holos

// This component is intended to be applied only to the provisioner cluster.
// Resources:
// - RBAC for User eso-creds-refresher IAM Service Account to issue tokens for eso-reader and eso-writer in each namespace.

// For each namespace:
// - Namespace
// - ServiceAccount eso-reader, eso-writer

// No flux kustomization
ksObjects: []

#KubernetesObjects & {
	apiObjects: {
		let role = #CredsRefresherIAM.role
		let binding = #CredsRefresherIAM.binding
		ClusterRole: "\(role.metadata.name)":           role
		ClusterRoleBinding: "\(binding.metadata.name)": binding
		for ns in #PlatformNamespaces {
			for obj in (#PlatformNamespaceObjects & {_ns: ns}).objects {
				let Kind = obj.kind
				let Name = obj.metadata.name
				"\(Kind)": "\(ns.name)/\(Name)": obj
			}
		}

		for nsName, ns in #ManagedNamespaces {
			for obj in (#PlatformNamespaceObjects & {_ns: ns.namespace.metadata}).objects {
				let Kind = obj.kind
				let Name = obj.metadata.name
				"\(Kind)": "\(nsName)/\(Name)": obj
			}
		}
	}
}

#InputKeys: {
	cluster:   "provisioner"
	project:   "secrets"
	component: "eso-creds-refresher"
}

// #CredsRefresherIAM defines the rbac policy for the job that refreshes credentials used by eso SecretStore resources in clusters other than the provisioner cluster.
#CredsRefresherIAM: {
	let _name = #CredsRefresher.name

	// Allow the IAM Service Account to determine which service account tokens to refresh.
	// List namespaces.
	// List service accounts.
	role: #ClusterRole & {
		metadata: name: _name
		rules: [
			{
				apiGroups: [""]
				resources: ["namespaces"]
				verbs: ["list"]
			},
			{
				apiGroups: [""]
				resources: ["serviceaccounts"]
				verbs: ["list", "get"]
			},
		]
	}
	binding: #ClusterRoleBinding & {
		metadata: name: _name
		roleRef: {
			apiGroup: "rbac.authorization.k8s.io"
			kind:     "ClusterRole"
			name:     _name
		}
		subjects: [
			{
				apiGroup: ""
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
				kind: "User"
				name: #CredsRefresher.iamServiceAccount
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
			// This label is important, the Job running on the other clusters selects this label to determine which service accounts to issue tokens for.
			labels: "holos.run/job.name": #CredsRefresher.name
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
				kind:      "ServiceAccount"
				name:      _name
				namespace: _namespace
			},
		]
	}
}
