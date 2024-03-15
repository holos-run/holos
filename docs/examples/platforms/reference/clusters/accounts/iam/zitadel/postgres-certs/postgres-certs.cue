package holos

#InputKeys: component: "postgres-certs"

let SecretNames = {
	[Name=_]: {name: Name}
	"\(_DBName)-primary-tls": _
	"\(_DBName)-repl-tls":    _
	"\(_DBName)-client-tls":  _
	"\(_DBName)-root-ca":     _
}

#Kustomization: spec: targetNamespace: #TargetNamespace
#Kustomization: spec: healthChecks: [
	for s in SecretNames {
		apiVersion: "external-secrets.io/v1beta1"
		kind:       "ExternalSecret"
		name:       s.name
		namespace:  #TargetNamespace
	},
]

#KubernetesObjects & {
	apiObjects: {
		for s in SecretNames {
			ExternalSecret: "\(s.name)": _
		}
	}
}
