package holos

// e.g. prod-secrets-namespaces
name: "\(#InputKeys.stage)-\(#InputKeys.project)-namespaces"
objects: [
	#Namespace & {
		metadata: name: "external-secrets"
	}
]
