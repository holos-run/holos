package holos

#InputKeys: component: "zitadel"

// Upstream helm chart doesn't specify the namespace field for all resources.
#Kustomization: spec: targetNamespace: #TargetNamespace

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name:    "zitadel"
		version: "7.9.0"
		repository: {
			name: "zitadel"
			url:  "https://charts.zitadel.com"
		}
	}
	values: #Values

	apiObjects: {
		ExternalSecret: masterkey: #ExternalSecret & {
			_name: "zitadel-masterkey"
		}
	}
}
