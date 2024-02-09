package holos

// e.g. prod-secrets-namespaces
metadata: name: "\(#InputKeys.stage)-\(#InputKeys.project)-namespaces"
objects: [
	#Namespace & {
		metadata: name: "external-secrets"
	}
]
