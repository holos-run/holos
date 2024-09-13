package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Objects = {
	Name:      "{{ .Name }}"
	Namespace: #HTTPBin.Namespace

	// Constrain resources to the httpbin namespace
	Resources: [_]: [_]: metadata: namespace: #HTTPBin.Namespace

	Resources: {
		Deployment: httpbin: {
			metadata: name: "httpbin"
			metadata: labels: app: metadata.name
			spec: {
				selector: matchLabels: {
					app:                          metadata.labels.app
					"app.kubernetes.io/instance": app
				}

				template: {
					metadata: labels: selector.matchLabels
					spec: securityContext: seccompProfile: type: "RuntimeDefault"
					spec: containers: [{
						name:  "httpbin"
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
					port:       #HTTPBin.Port
					targetPort: Deployment.httpbin.spec.template.spec.containers[0].ports[0].containerPort
					protocol:   "TCP"
					name:       "http"
				}
				ports: [for x in _ports {x}]
			}
		}

		// Allow istio-ingress to refer to Services from HTTPRoutes
		ReferenceGrant: httpbin: #ReferenceGrant & {
			metadata: labels: Deployment.httpbin.metadata.labels
		}
	}
}
