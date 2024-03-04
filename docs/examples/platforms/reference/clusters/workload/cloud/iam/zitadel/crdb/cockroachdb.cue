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
		Issuer: {
			// https://github.com/cockroachdb/helm-charts/blob/3dcf96726ebcfe3784afb526ddcf4095a1684aea/README.md?plain=1#L196-L201
			cockroachdb: #Issuer & {
				metadata: name:      #ComponentName
				metadata: namespace: #TargetNamespace
				spec: selfSigned: {}
			}
		}
	}
}
