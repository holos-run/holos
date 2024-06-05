package holos

_AppInfo: #AppInfo & {
	metadata: name: "holos"
	metadata: labels: {
		"app.holos.run/environment":  spec.env
		"app.holos.run/name":         metadata.name
		"app.holos.run/component":    spec.component
		"render.holos.run/component": status.component
	}

	spec: env:       string
	spec: port:      3000
	spec: component: "app" | "infra" | "routes"

	spec: dns: segments: name: ["app"]

	status: component: spec.env + "-" + metadata.name + "-" + spec.component
}
