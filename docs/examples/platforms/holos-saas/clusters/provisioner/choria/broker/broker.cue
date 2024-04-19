package holos

let Namespace = "jeff-holos"
let Broker = "choria-broker"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-platform-issuer": _

		metadata: name: "\(Namespace)-\(Broker)"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let SelectorLabels = {
	"app.kubernetes.io/instance": Broker
	"app.kubernetes.io/name":     Broker
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		Certificate: "\(Broker)-tls": #Certificate & {
			metadata: {
				name:      "\(Broker)-tls"
				namespace: Namespace
				labels:    SelectorLabels
			}
			spec: {
				commonName: "\(Broker).\(Namespace).svc.cluster.local"
				dnsNames: [
					Broker,
					"\(Broker).\(Namespace).svc",
					"\(Broker).\(Namespace).svc.cluster.local",
					"*.\(Broker)",
					"*.\(Broker).\(Namespace).svc",
					"*.\(Broker).\(Namespace).svc.cluster.local",
				]
				issuerRef: kind: "ClusterIssuer"
				issuerRef: name: "platform-issuer"
				secretName: metadata.name
				usages: ["signing", "key encipherment", "server auth", "client auth"]
			}
		}
	}
}
