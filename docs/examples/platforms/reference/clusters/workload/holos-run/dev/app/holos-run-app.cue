package holos

let Namespace = "dev-holos"
let Holos = "holos"

// spec represents the output provided to holos
spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "dev-holos-app"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

// OBJECTS represents the kubernetes api objects to manage.
let OBJECTS = #APIObjects & {
	apiObjects: Deployment: holos: {
		metadata: {
			name:      Holos
			namespace: Namespace
			labels: app: Holos
		}
		spec: {
			selector: matchLabels: app: Holos
			template: metadata: labels: {
				app:                       Holos
				"sidecar.istio.io/inject": "true"
			}
			strategy: rollingUpdate: maxSurge:       1
			strategy: rollingUpdate: maxUnavailable: 0
			template: {
				spec: {
					serviceAccountName: Holos
					securityContext: seccompProfile: type: "RuntimeDefault"
					containers: [
						{
							name:            Holos
							image:           "271053619184.dkr.ecr.us-east-2.amazonaws.com/holos-run/holos-server/holos:0.73.1"
							imagePullPolicy: "Always"
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
								cpu:    "0.25"
								memory: "256Mi"
							}
							resources: requests: resources.limits
						},
					]
				}
			}
		}
	}
}
