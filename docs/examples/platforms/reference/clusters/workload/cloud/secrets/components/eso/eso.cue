package holos

// Manages the External Secrets Operator from the official upstream Helm chart.

#TargetNamespace: "external-secrets"

#InputKeys: component: "eso"

#InputKeys: {
	project: "secrets"
	service: "eso"
}

#Kustomization: spec: {
	dependsOn: [{name: #InstancePrefix + "-namespaces"}]
	targetNamespace: #TargetNamespace
}

#HelmChart & {
	values: installCrds: true
	namespace: #TargetNamespace
	chart: {
		name:    "external-secrets"
		version: "0.9.12"
		repository: {
			name: "external-secrets"
			url:  "https://charts.external-secrets.io"
		}
	}
}
