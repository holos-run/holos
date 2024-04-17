package holos

// Certificate used by the ingress to connect to services using a platform
// issued certificate but which are not using istio sidecar injection.
// Examples are keycloak, vault, nats, choria, etc...

let Namespace = "istio-ingress"
let CertName = "istio-ingress-mtls-cert"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-platform-issuer": _

		metadata: name: CertName
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		Certificate: "\(CertName)": #Certificate & {
			metadata: {
				name:      CertName
				namespace: Namespace
			}
			spec: {
				secretName: metadata.name
				issuerRef: kind: "ClusterIssuer"
				issuerRef: name: "platform-issuer"
				commonName: "istio-ingress"
				dnsNames: [
					"istio-ingress",
					"istio-ingress.\(Namespace)",
					"istio-ingress.\(Namespace).svc",
					"istio-ingress.\(Namespace).svc.cluster.local",
				]
			}
		}
	}
}
