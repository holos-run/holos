package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).Output

let Objects = {
	Name:      "httpbin-backend"
	Namespace: _HTTPBin.metadata.namespace

	// Constrain the metadata of all component resources.
	Resources: [_]: [_]: metadata: _HTTPBin.metadata

	// Grant the Gateway ns the ability to refer to the Service from HTTPRoutes.
	Resources: ReferenceGrant: (#IstioGatewaysNamespace): #ReferenceGrant

	Resources: {
		Deployment: httpbin: {
			metadata: labels: app: _HTTPBin.metadata.name
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
						name:  _HTTPBin.metadata.name
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
			metadata: labels: Deployment.httpbin.metadata.labels
			spec: {
				selector: Deployment.httpbin.spec.selector.matchLabels
				_ports: http: {
					port:       _HTTPBin.port
					targetPort: Deployment.httpbin.spec.template.spec.containers[0].ports[0].containerPort
					protocol:   "TCP"
					name:       "http"
				}
				ports: [for x in _ports {x}]
			}
		}
	}
}
