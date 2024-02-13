package holos

#Kustomization: spec: dependsOn: [{name: "\(#InstancePrefix)-namespaces"}]

#HelmChart & {
	values: installCrds: true
	namespace: #TargetNamespace
	chart: {
		name: "external-secrets"
		version: "0.9.12"
		repository: {
			name: "external-secrets"
		  url: "https://charts.external-secrets.io"
		}
	}
}
