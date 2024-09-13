package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Objects = {
	Name:      "zitadel-secrets"
	Namespace: _ZitadelNamespace

	Resources: {
		ExternalSecret: [_]: #ExternalCert & {metadata: namespace: Namespace}
		ExternalSecret: "\(_DBName)-primary-tls": _
		ExternalSecret: "\(_DBName)-repl-tls":    _
		ExternalSecret: "\(_DBName)-client-tls":  _
		ExternalSecret: "\(_DBName)-root-ca":     _
	}
}
