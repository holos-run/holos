package holos

// Refer to https://cert-manager.io/docs/configuration/selfsigned/#bootstrapping-ca-issuers

let Namespace = "cert-manager"

spec: components: KubernetesObjectsList: [
	#KubernetesObjects & {
		metadata: name: "prod-platform-issuer"

		_dependsOn: "prod-mesh-certmanager": _
		apiObjectMap: OBJECTS.apiObjectMap
	},
]

let SelfSigned = "platform-selfsigned"
let PlatformIssuer = "platform-issuer"

let OBJECTS = #APIObjects & {
	apiObjects: {
		ClusterIssuer: {
			"\(SelfSigned)": #ClusterIssuer & {
				metadata: name: SelfSigned
				spec: selfSigned: {}
			}
		}
		Certificate: {
			"\(PlatformIssuer)": #Certificate & {
				metadata: name:      PlatformIssuer
				metadata: namespace: Namespace
				spec: {
					duration:   "999999h"
					isCA:       true
					commonName: PlatformIssuer
					secretName: PlatformIssuer
					privateKey: algorithm: "ECDSA"
					privateKey: size:      256
					issuerRef: {
						name:  SelfSigned
						kind:  "ClusterIssuer"
						group: "cert-manager.io"
					}
				}
			}
		}
		ClusterIssuer: {
			"\(PlatformIssuer)": #ClusterIssuer & {
				metadata: name: PlatformIssuer
				spec: ca: secretName: PlatformIssuer
			}
		}
	}
}
