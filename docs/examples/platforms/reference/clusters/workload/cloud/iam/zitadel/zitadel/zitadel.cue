package holos

#InputKeys: component: "zitadel"

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
}
