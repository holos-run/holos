package holos

#InputKeys: component: "dbtls"
#KubernetesObjects & {
	apiObjects: {
		ExternalSecret: {
			zitadel: #ExternalSecret & {
				_name: "zitadel-client"
			}
			primary: #ExternalSecret & {
				_name: "\(_DBName)-primary"
			}
			repl: #ExternalSecret & {
				_name: "\(_DBName)-repl"
			}
		}
	}
}
