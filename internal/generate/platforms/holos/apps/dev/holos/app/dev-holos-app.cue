package holos

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Image = "quay.io/holos-run/holos:v0.83.1-7-gd9fe32b"

_AppInfo: spec: component: "app"

let Objects = {
	Name:      _AppInfo.status.component
	Namespace: _AppInfo.metadata.namespace

	Resources: [_]: [_]: metadata: namespace: Namespace

	let Metadata = _AppInfo.metadata

	Resources: {
		let MatchLabels = {"app.kubernetes.io/component": "server"}

		// Grant the Gateway ns the ability to refer to the Service from HTTPRoutes.
		ReferenceGrant: (#IstioGatewaysNamespace): #ReferenceGrant

		Service: holos: {
			apiVersion: "v1"
			metadata:   Metadata
			spec: {
				type:     "ClusterIP"
				selector: MatchLabels
				ports: [{
					appProtocol: "http2"
					name:        "http"
					port:        _AppInfo.spec.port
					protocol:    "TCP"
					targetPort:  _AppInfo.spec.port
				}, {
					appProtocol: "http"
					name:        "metrics"
					port:        9090
					protocol:    "TCP"
					targetPort:  9090
				}]
			}
		}
		Deployment: holos: {
			metadata: Metadata
			metadata: labels: MatchLabels
			spec: {
				selector: matchLabels: MatchLabels

				template: metadata: labels: Metadata.labels
				template: metadata: labels: MatchLabels
				template: metadata: labels: "sidecar.istio.io/inject": "true"

				strategy: rollingUpdate: maxSurge:       1
				strategy: rollingUpdate: maxUnavailable: 0
				template: {
					spec: {
						serviceAccountName: Metadata.name
						securityContext: seccompProfile: type: "RuntimeDefault"
						containers: [
							{
								name:            Metadata.name
								image:           Image
								imagePullPolicy: "IfNotPresent"
								command: [
									"/app/bin/holos",
									"server",
									"--log-format=json",
									"--oidc-issuer=\(_AuthProxy.issuerURL)",
									"--oidc-audience=\(_AuthProxy.projectID)",
								]
								env: [
									{
										name:  "TZ"
										value: "America/Los_Angeles"
									},
									{
										name: "DATABASE_URL"
										valueFrom: secretKeyRef: {
											key:  "uri"
											name: "holos-pguser-holos"
										}
									},
								]
								ports: [
									{
										containerPort: 3000
										name:          "http"
										protocol:      "TCP"
									},
								]
								securityContext: capabilities: drop: ["ALL"]
								securityContext: allowPrivilegeEscalation: false
								securityContext: runAsNonRoot:             true
								resources: limits: {
									cpu:    "0.5"
									memory: "512Mi"
								}
								resources: requests: resources.limits
							},
						]
					}
				}
			}
		}
	}
}
