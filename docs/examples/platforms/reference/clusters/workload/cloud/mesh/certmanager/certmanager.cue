package holos

// https://cert-manager.io/docs/

#TargetNamespace: "cert-manager"

#InputKeys: {
	component: "certmanager"
	service:   "cert-manager"
}

#HelmChart & {
	values: installCRDs: true
	namespace: #TargetNamespace
	chart: {
		name:    "cert-manager"
		version: "1.14.3"
		repository: {
			name: "jetstack"
			url:  "https://charts.jetstack.io"
		}
	}
}
