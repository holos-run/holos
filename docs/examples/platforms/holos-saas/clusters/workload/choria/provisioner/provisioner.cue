package holos

let Namespace = "jeff-holos"
let Provisioner = "choria-provisioner"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-stores": _

		metadata: name: "\(Namespace)-\(Provisioner)"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let SelectorLabels = {
	"app.kubernetes.io/instance": Provisioner
	"app.kubernetes.io/name":     Provisioner
}

let Metadata = {
	name:      Provisioner
	namespace: Namespace
	labels:    SelectorLabels
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: "\(Provisioner)-tls": #ExternalSecret & {
			metadata: name:      "\(Provisioner)-tls"
			metadata: namespace: Namespace
		}
		ExternalSecret: "\(Provisioner)": #ExternalSecret & {
			metadata: name:      Provisioner
			metadata: namespace: Namespace
		}
		ServiceAccount: "\(Provisioner)": #ServiceAccount & {
			metadata: Metadata
		}
		Deployment: "\(Provisioner)": {
			metadata: Metadata
			spec: {
				selector: matchLabels: SelectorLabels
				template: metadata: labels: SelectorLabels
				template: spec: {
					containers: [
						{
							name: Provisioner
							command: ["bash", "/etc/provisioner/entrypoint"]
							// skopeo inspect docker://registry.choria.io/choria/provisioner | jq .RepoTags
							image:           "registry.choria.io/choria/provisioner:0.15.1"
							imagePullPolicy: "IfNotPresent"
							resources: {}
							securityContext: {}
							volumeMounts: [
								{
									mountPath: "/etc/provisioner"
									name:      Provisioner
								},
								{
									mountPath: "/etc/provisioner-tls"
									name:      "\(Provisioner)-tls"
								},
							]
						},
					]
					securityContext: {}
					serviceAccountName: Provisioner
					volumes: [
						{
							name: Provisioner
							secret: secretName: name
						},
						{
							name: "\(Provisioner)-tls"
							secret: secretName: name
						},
					]
				}
			}
		}
	}
}
