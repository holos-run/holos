package holos

#Values: {
	image: repository: "quay.io/holos/cockroachdb/cockroach"

	fullnameOverride: #ComponentName

	tls: {
		enabled: true
		certs: {
			// https://github.com/cockroachdb/helm-charts/blob/3dcf96726ebcfe3784afb526ddcf4095a1684aea/README.md?plain=1#L204-L215
			selfSigner: enabled: false
			certManager: true
			certManagerIssuer: {
				kind: "Issuer"
				name: #ComponentName
			}
		}
	}

	storage: persistentVolume: {
		enabled: true
		size:    "1Gi"
	}
}
