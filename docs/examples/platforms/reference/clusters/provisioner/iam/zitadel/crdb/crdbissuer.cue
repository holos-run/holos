package holos

// Manage an Issuer for cockroachdb for zitadel.
// For the iam login service, zitadel connects to cockroach db using tls certs for authz.
// Upstream: "The recommended approach is to use cert-manager for certificate management. For details, refer to Deploy cert-manager for mTLS."
// Refer to https://www.cockroachlabs.com/docs/stable/secure-cockroachdb-kubernetes#deploy-cert-manager-for-mtls

#InputKeys: component: "crdb"

#KubernetesObjects & {
	apiObjects: {
		Issuer: {
			// https://github.com/cockroachdb/helm-charts/blob/3dcf96726ebcfe3784afb526ddcf4095a1684aea/README.md?plain=1#L196-L201
			crdb: #Issuer & {
				_description: "Issues the self signed root ca cert for cockroach db"
				metadata: name:      #ComponentName
				metadata: namespace: #TargetNamespace
				spec: selfSigned: {}
			}
			"crdb-ca-issuer": #Issuer & {
				_description: "Issues mtls certs for cockroach db"
				metadata: name:      "crdb-ca-issuer"
				metadata: namespace: #TargetNamespace
				spec: ca: secretName: "cockroach-ca"
			}
		}
		Certificate: {
			"crdb-ca-cert": #Certificate & {
				_description: "Root CA cert for cockroach db"
				metadata: name:      "crdb-ca-cert"
				metadata: namespace: #TargetNamespace
				spec: {
					commonName: "root"
					isCA:       true
					issuerRef: group:      "cert-manager.io"
					issuerRef: kind:       "Issuer"
					issuerRef: name:       "crdb"
					privateKey: algorithm: "ECDSA"
					privateKey: size:      256
					secretName: "cockroach-ca"
					subject: organizations: ["Cockroach"]
				}
			}
			"crdb-node": #Certificate & {
				metadata: name:      "crdb-node"
				metadata: namespace: #TargetNamespace
				spec: {
					commonName: "node"
					dnsNames: [
						"localhost",
						"127.0.0.1",
						"crdb-public",
						"crdb-public.\(#TargetNamespace)",
						"crdb-public.\(#TargetNamespace).svc.cluster.local",
						"*.crdb",
						"*.crdb.\(#TargetNamespace)",
						"*.crdb.\(#TargetNamespace).svc.cluster.local",
					]
					duration: "876h"
					issuerRef: group:      "cert-manager.io"
					issuerRef: kind:       "Issuer"
					issuerRef: name:       "crdb-ca-issuer"
					privateKey: algorithm: "RSA"
					privateKey: size:      2048
					renewBefore: "168h"
					secretName:  "cockroachdb-node"
					subject: organizations: ["Cockroach"]
					usages: ["digital signature", "key encipherment", "server auth", "client auth"]
				}
			}
			"crdb-root-client": #Certificate & {
				metadata: name:      "crdb-root-client"
				metadata: namespace: #TargetNamespace
				spec: {
					commonName: "root"
					duration:   "672h"
					issuerRef: group:      "cert-manager.io"
					issuerRef: kind:       "Issuer"
					issuerRef: name:       "crdb-ca-issuer"
					privateKey: algorithm: "RSA"
					privateKey: size:      2048
					renewBefore: "48h"
					secretName:  "cockroachdb-root"
					subject: organizations: ["Cockroach"]
					usages: ["digital signature", "key encipherment", "client auth"]
				}
			}
		}
		Certificate: zitadel: #Certificate & {
			metadata: name:      "crdb-zitadel-client"
			metadata: namespace: #TargetNamespace
			spec: {
				commonName: "zitadel"
				issuerRef: {
					group: "cert-manager.io"
					kind:  "Issuer"
					name:  "crdb-ca-issuer"
				}
				privateKey: algorithm: "RSA"
				privateKey: size:      2048
				renewBefore: "48h0m0s"
				secretName:  "cockroachdb-zitadel"
				subject: organizations: ["Cockroach"]
				usages: ["digital signature", "key encipherment", "client auth"]
			}
		}
	}
}
