package holos

// objects are kubernetes api objects to apply
objects: #CredsRefresherService.objects

#InputKeys: {
	project:   "secrets"
	component: "eso-creds-refresher"
}

#TargetNamespace: #CredsRefresher.namespace

// #CredsRefresherService defines the job that refreshes credentials used by eso SecretStore resources.
#CredsRefresherService: {
	objects: [
		// Kubernetes ServiceAccount used by the Job.
		#ServiceAccount & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
				annotations: "holos.run/description": "Refreshes provisioner cluster eso-reader and eso-writer service account credentials."
			}
		},
		// Job needs to read and write secrets across the cluster.
		#ClusterRole & {
			metadata: {
				name: #CredsRefresher.name
			}
			rules: [
				{
					apiGroups: [""]
					resources: ["secrets"]
					verbs: ["*"]
				},
			]
		},
		// Bind the Role to the ServiceAccount for the Job.
		#ClusterRoleBinding & {
			metadata: {
				name: #CredsRefresher.name
			}
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "ClusterRole"
				name:     #CredsRefresher.name
			}
			subjects: [
				{
					kind:      "ServiceAccount"
					name:      #CredsRefresher.name
					namespace: #CredsRefresher.namespace
				},
			]
		},
	]
}

{} & #KubernetesObjects
