package holos

let Name = "httpbin"
let SecretName = #InputKeys.cluster + "-" + Name
let MatchLabels = {app: Name} & #SelectorLabels
let Metadata = {
	name:      Name
	namespace: #TargetNamespace
	labels: app: Name
}

#InputKeys: component: Name

#TargetNamespace: "istio-ingress"
#DependsOn:       _IngressGateway

#KubernetesObjects & {
	apiObjects: {
		Certificate: httpbin: #HTTP01Cert & {
			_name:   Name
			_secret: SecretName
		}
		Deployment: httpbin: #Deployment & {
			metadata: Metadata
			spec: selector: matchLabels: MatchLabels
			spec: template: {
				metadata: labels: MatchLabels
				metadata: labels: #CommonLabels
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
						runAsUser:                1337
						runAsGroup:               1337
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
	}
}
