package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "backstage-secrets"
	Namespace: "backstage"

	Resources: {
		ExternalSecret: [_]: #ExternalSecret & {metadata: namespace: Namespace}

		ExternalSecret: "\(_DBName)-primary-tls": #ExternalCert
		ExternalSecret: "\(_DBName)-repl-tls":    #ExternalCert
		ExternalSecret: "\(_DBName)-client-tls":  #ExternalCert
		ExternalSecret: "\(_DBName)-root-ca":     #ExternalCert

		ExternalSecret: "pgbackrest": _
	}
}
