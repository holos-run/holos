package holos

let Namespace = "jeff-holos"
let Provisioner = "choria-provisioner"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-platform-issuer": _

		metadata: name: "\(Namespace)-\(Provisioner)"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let SelectorLabels = {
	"app.kubernetes.io/instance": Provisioner
	"app.kubernetes.io/name":     Provisioner
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		Certificate: "\(Provisioner)-tls": #Certificate & {
			metadata: {
				name:      "\(Provisioner)-tls"
				namespace: Namespace
				labels:    SelectorLabels
			}
			spec: {
				commonName: "\(Provisioner).\(Namespace).svc.cluster.local"
				dnsNames: [
					Provisioner,
					"\(Provisioner).\(Namespace).svc",
					"\(Provisioner).\(Namespace).svc.cluster.local",
					"*.\(Provisioner)",
					"*.\(Provisioner).\(Namespace).svc",
					"*.\(Provisioner).\(Namespace).svc.cluster.local",
				]
				issuerRef: kind: "ClusterIssuer"
				issuerRef: name: "platform-issuer"
				secretName: metadata.name
				usages: ["signing", "key encipherment", "server auth", "client auth"]
			}
		}
	}
}
