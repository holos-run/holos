package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Host = "httpbin.admin.\(_ClusterName).\(_Platform.Model.org.domain)"

let Objects = {
	Name:      "httpbin"
	Namespace: "holos-system"

	Resources: {
		Deployment: httpbin: {
			metadata: {
				namespace: Namespace
				labels: app: Name
			}
			spec: {
				selector: matchLabels: {
					app:                          metadata.labels.app
					"app.kubernetes.io/instance": app
				}

				template: {
					metadata: labels: "sidecar.istio.io/inject": "true"
					metadata: labels: selector.matchLabels
					spec: securityContext: seccompProfile: type: "RuntimeDefault"
					spec: containers: [{
						name:  Name
						image: "quay.io/holos/mccutchen/go-httpbin"
						ports: [{containerPort: 8080}]
						securityContext: {
							seccompProfile: type: "RuntimeDefault"
							allowPrivilegeEscalation: false
							runAsNonRoot:             true
							runAsUser:                8192
							runAsGroup:               8192
							capabilities: drop: ["ALL"]
						}}]
				}
			}
		}
		Service: httpbin: {
			metadata: name:      string
			metadata: namespace: Namespace
			metadata: labels:    Deployment.httpbin.metadata.labels
			spec: {
				selector: Deployment.httpbin.spec.selector.matchLabels
				_ports: http: {
					port:       80
					targetPort: 8080
					protocol:   "TCP"
					name:       "http"
				}
				ports: [for x in _ports {x}]
			}
		}
		HTTPRoute: httpbin: {
			metadata: namespace: Namespace
			metadata: labels:    Deployment.httpbin.metadata.labels
			spec: hostnames: [Host]
			spec: parentRefs: [
				{
					name:      "default"
					namespace: #IstioGatewaysNamespace
				},
			]
			spec: rules: [
				{
					matches: [{path: {type: "PathPrefix", value: "/"}}]
					backendRefs: [
						{
							name: Service.httpbin.metadata.name
							port: Service.httpbin.spec._ports.http.port
						},
					]
				},
			]
		}
	}
}
