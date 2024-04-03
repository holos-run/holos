package holos

let Name = "httpbin"
let ComponentName = "\(#InstancePrefix)-\(Name)"

let MatchLabels = {
	app:                          Name
	"app.kubernetes.io/instance": ComponentName
}
let Metadata = {
	name:      Name
	namespace: #TargetNamespace
	labels: app: Name
}

#InputKeys: component: Name

#TargetNamespace: "istio-ingress"

let Cert = #PlatformCerts["\(#ClusterName)-httpbin"]

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-namespaces":       _
		_dependsOn: "\(#InstancePrefix)-istio-base": _
		_dependsOn: "\(#InstancePrefix)-ingress":    _

		metadata: name: ComponentName

		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: "\(Cert.spec.secretName)": _
		Deployment: httpbin: #Deployment & {
			metadata: Metadata
			spec: selector: matchLabels: MatchLabels
			spec: template: {
				metadata: labels: MatchLabels
				metadata: labels: #IstioSidecar
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
		Service: httpbin: #Service & {
			metadata: Metadata
			spec: selector: MatchLabels
			spec: ports: [
				{port: 80, targetPort: 8080, protocol: "TCP", name: "http"},
			]
		}
		VirtualService: httpbin: #VirtualService & {
			metadata: Metadata
			spec: hosts: [for host in Cert.spec.dnsNames {host}]
			spec: gateways: ["istio-ingress/default"]
			spec: http: [{route: [{destination: host: Name}]}]
		}
	}
}
