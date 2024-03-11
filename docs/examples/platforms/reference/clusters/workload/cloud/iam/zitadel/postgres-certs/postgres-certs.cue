package holos

#InputKeys: component: "postgres-certs"
#KubernetesObjects & {
	apiObjects: {
		ExternalSecret: {
			"\(_DBName)-primary-tls": _
			"\(_DBName)-repl-tls":    _
			"\(_DBName)-client-tls":  _
			"\(_DBName)-root-ca":     _
		}
	}
}
