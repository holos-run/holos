package holos

_IAP: {
	metadata: {
		name:      string
		namespace: _AuthProxy.metadata.namespace
		labels: "app.kubernetes.io/name":    name
		labels: "app.kubernetes.io/part-of": "default-gateway"
		...
	}
}
