package holos

#InputKeys: component: "crdb"

#HelmChart & {
	namespace: #TargetNamespace
	chart: {
		name:    "cockroachdb"
		version: "11.2.3"
		repository: {
			name: "cockroachdb"
			url:  "https://charts.cockroachdb.com/"
		}
	}
	values: #Values
	apiObjects: {
		ExternalSecret: node: #ExternalSecret & {_name: "cockroachdb-node"}
		ExternalSecret: root: #ExternalSecret & {_name: "cockroachdb-root"}
	}
}
