package holos

// Platform wide configuration
#ExternalSecrets: {
	Version:   "{{ .Version }}"
	Namespace: "external-secrets"
}

// Register the namespace
#Namespaces: (#ExternalSecrets.Namespace): _
