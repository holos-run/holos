package holos

// PlatformNamespace is a namespace to manage for Secret provisioning, SecretStore, etc...
#PlatformNamespace: {
	name: string
}

// #PlatformNamespaces is a list of namespaces to manage across the platform.
#PlatformNamespaces: [...#PlatformNamespace]

// #CredsRefresher defines attributes for the external secrets credential refresher job.
#CredsRefresher: {
	// name is the name of the job, component, and related resources.
	name: "eso-creds-refresher"
	// namespace is the namespace where the job runs.
	namespace: "holos-system"
	// project is the gcp project id hosting the provisioner cluster.
	project: #InputKeys.gcpProjectID
	// iamsa is the iam service account email address.
	iamServiceAccount: name + "@" + project + ".iam.gserviceaccount.com"
}
