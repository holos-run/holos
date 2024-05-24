package holos

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

let NAME = "eso-creds-manager"
let REFRESHER = "eso-creds-refresher"
let EMAIL = _Platform.Model.eso.gcpServiceAccount

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      NAME
	Namespace: "holos-system"

	Resources: {
		// For Job eso-creds-refresher
		ClusterRole: "\(REFRESHER)": rbacv1.#ClusterRole & {
			metadata: name: REFRESHER
			metadata: annotations: "holos.run/description": "For \(EMAIL) assumed by Job eso-creds-refresher via workload identity."
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

		// Bind the Role to the Google IAM Service Account which the
		// eso-creds-refresher Job assumes via workload identity.
		ClusterRoleBinding: "\(REFRESHER)": rbacv1.#ClusterRoleBinding & {
			metadata: name: REFRESHER
			metadata: annotations: "holos.run/description": "For \(EMAIL) assumed by Job eso-creds-refresher via workload identity."
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "ClusterRole"
				name:     REFRESHER
			}
			subjects: [
				{
					apiGroup: ""
					kind:     "User"
					name:     EMAIL
				},
			]
		}
	}

	// Manage resources in every managed namespace
	for NS in _Namespaces {
		for obj in (#PlatformNamespaceObjects & {_ns: NS.metadata}).objects {
			Resources: "\(obj.kind)": "\(NS.metadata.name)/\(obj.metadata.name)": obj
		}
	}
}

// #PlatformNamespaceObjects defines the api objects necessary for eso SecretStores in external clusters to access secrets in a given namespace in the provisioner cluster.
#PlatformNamespaceObjects: {
	_ns: {
		name: string
		labels?: {[string]: string}
		annotations?: {[string]: string}
	}

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

// #IAMServiceAccount represents the k8s resources for the iam service account
// used by the eso-creds-refresher Job to get credentials for each of the
// eso-reader and eso-writer service accounts defined by #ESOServiceAccount
#IAMServiceAccount: {
	namespace: string
	let _namespace = namespace

	// Role for the eso-creds-refresher
	role: rbacv1.#Role & {
		metadata: {
			name:      REFRESHER
			namespace: _namespace
			annotations: "holos.run/description": "For \(EMAIL) assumed by Job \(REFRESHER) via workload identity."
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
	binding: rbacv1.#RoleBinding & {
		metadata: {
			name:      REFRESHER
			namespace: _namespace
			annotations: "holos.run/description": "For \(EMAIL) assumed by Job \(REFRESHER) via workload identity."
		}
		roleRef: {
			apiGroup: "rbac.authorization.k8s.io"
			kind:     "Role"
			name:     REFRESHER
		}
		subjects: [
			{
				kind: "User"
				name: EMAIL
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
	ksa: corev1.#ServiceAccount & {
		metadata: {
			name:      _name
			namespace: _namespace
			annotations: "holos.run/description": "Job \(REFRESHER) creates short lived tokens for this ksa for a SecretStore in workload clusters."
			// Selector labels for the eso-creds-refresher Job
			labels: "holos.run/job.name": "eso-creds-refresher"
			labels: "holos.run/iam.role": mode
		}
	}

	// role is the Role of the ServiceAccount used by the SecretStore
	role: rbacv1.#Role & {
		metadata: {
			name:      _name
			namespace: _namespace
			annotations: "holos.run/description": "For use by a SecretStore in workload clusters."
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
	binding: rbacv1.#RoleBinding & {
		metadata: {
			name:      _name
			namespace: _namespace
			annotations: "holos.run/description": "For use by a SecretStore in workload clusters."
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
