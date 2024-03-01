package holos

// https://cert-manager.io/docs/

#TargetNamespace: "cert-manager"

#InputKeys: {
	project: "mesh"
	component: "certmanager"
	service: "cert-manager"
}

#HelmChart & {
	values: installCrds: true
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
