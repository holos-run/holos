package holos

let Namespace = "jeff-holos"
let Broker = "choria-broker"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		_dependsOn: "prod-secrets-stores": _

		metadata: name: "\(Namespace)-\(Broker)"
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let SelectorLabels = {
	"app.kubernetes.io/part-of": "choria"
	"app.kubernetes.io/name":    Broker
}

let Metadata = {
	name:      Broker
	namespace: Namespace
	labels:    SelectorLabels
}

let OBJECTS = #APIObjects & {
	apiObjects: {
		ExternalSecret: "\(Broker)-tls": #ExternalSecret & {
			metadata: name:      "\(Broker)-tls"
			metadata: namespace: Namespace
		}
		ExternalSecret: "\(Broker)": #ExternalSecret & {
			metadata: name:      Broker
			metadata: namespace: Namespace
		}
		StatefulSet: "\(Broker)": {
			metadata: Metadata
			spec: {
				selector: matchLabels: SelectorLabels
				serviceName: Broker
				template: metadata: labels: SelectorLabels
				template: spec: {
					containers: [
						{
							name: Broker
							command: ["choria", "broker", "run", "--config", "/etc/choria/broker.conf"]
							image:           "registry.choria.io/choria/choria:0.28.0"
							imagePullPolicy: "IfNotPresent"
							ports: [
								{
									containerPort: 4222
									name:          "tcp-nats"
									protocol:      "TCP"
								},
								{
									containerPort: 4333
									name:          "https-wss"
									protocol:      "TCP"
								},
								{
									containerPort: 5222
									name:          "tcp-cluster"
									protocol:      "TCP"
								},
								{
									containerPort: 8222
									name:          "http-stats"
									protocol:      "TCP"
								},
							]
							livenessProbe: httpGet: {
								path: "/healthz"
								port: "http-stats"
							}
							readinessProbe: livenessProbe
							resources: {}
							securityContext: {}
							volumeMounts: [
								{
									mountPath: "/etc/choria"
									name:      Broker
								},
								{
									mountPath: "/etc/choria-tls"
									name:      "\(Broker)-tls"
								},
							]
						},
					]
					securityContext: {}
					serviceAccountName: Broker
					volumes: [
						{
							name: Broker
							secret: secretName: Broker
						},
						{
							name: "\(Broker)-tls"
							secret: secretName: "\(Broker)-tls"
						},
					]
				}
			}
		}
		ServiceAccount: "\(Broker)": #ServiceAccount & {
			metadata: Metadata
		}
		Service: "\(Broker)": #Service & {
			metadata: Metadata
			spec: {
				type:      "ClusterIP"
				clusterIP: "None"
				selector:  SelectorLabels
				ports: [
					{
						name:        "tcp-nats"
						appProtocol: "tcp"
						port:        4222
						protocol:    "TCP"
						targetPort:  "tcp-nats"
					},
					{
						name:        "tcp-cluster"
						appProtocol: "tcp"
						port:        5222
						protocol:    "TCP"
						targetPort:  "tcp-cluster"
					},
					{
						name:        "https-wss"
						appProtocol: "https"
						port:        443
						protocol:    "TCP"
						targetPort:  "https-wss"
					},
				]
			}
		}
		DestinationRule: "\(Broker)-wss": #DestinationRule & {
			_decriptions: "Configures Istio to connect to Choria using a cert issued by the Platform Issuer"
			metadata:     Metadata
			spec: host: "\(Broker).\(Namespace).svc.cluster.local"
			spec: trafficPolicy: tls: {
				credentialName: "istio-ingress-mtls-cert"
				mode:           "MUTUAL"
				// subjectAltNames is important, otherwise istio will fail to verify the
				// choria broker upstream server.  make sure this matches a value
				// present in the choria broker's cert.
				//
				//  kubectl get secret choria-broker-tls -o json | jq --exit-status
				//  '.data | map_values(@base64d)' | jq .\"tls.crt\" -r | openssl x509
				//  -text -noout -in -
				subjectAltNames: [spec.host]
			}
		}
		VirtualService: "\(Broker)-wss": #VirtualService & {
			metadata: name:      "\(Broker)-wss"
			metadata: namespace: Namespace
			spec: {
				gateways: ["istio-ingress/default"]
				hosts: ["jeff.provision.dev.\(#ClusterName).holos.run"]
				http: [
					{
						route: [
							{
								destination: {
									host: "\(Broker).\(Namespace).svc.cluster.local"
									port: "number": 443
								}
							},
						]
					},
				]
			}
		}
	}
}
