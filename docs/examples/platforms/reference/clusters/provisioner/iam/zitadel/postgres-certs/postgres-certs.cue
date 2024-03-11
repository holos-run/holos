package holos

// Manage an Issuer for the database.

// Both cockroach and postgres handle tls database connections with cert manager
// PGO: https://github.com/CrunchyData/postgres-operator-examples/tree/main/kustomize/certmanager/certman
// CRDB: https://github.com/cockroachdb/helm-charts/blob/3dcf96726ebcfe3784afb526ddcf4095a1684aea/README.md?plain=1#L196-L201

// Refer to [Using Cert Manager to Deploy TLS for Postgres on Kubernetes](https://www.crunchydata.com/blog/using-cert-manager-to-deploy-tls-for-postgres-on-kubernetes)

#InputKeys: component: "postgres-certs"

let SelfSigned = "\(_DBName)-selfsigned"
let RootCA = "\(_DBName)-root-ca"
let Orgs = ["Database"]

#KubernetesObjects & {
	apiObjects: {
		// Put everything in the target namespace.
		[_]: {
			[Name=_]: {
				metadata: name:      Name
				metadata: namespace: #TargetNamespace
			}
		}
		Issuer: {
			"\(SelfSigned)": #Issuer & {
				_description: "Self signed issuer to issue ca certs"
				metadata: name: SelfSigned
				spec: selfSigned: {}
			}
			"\(RootCA)": #Issuer & {
				_description: "Root signed intermediate ca to issue mtls database certs"
				metadata: name: RootCA
				spec: ca: secretName: RootCA
			}
		}
		Certificate: {
			"\(RootCA)": #Certificate & {
				_description: "Root CA cert for database"
				metadata: name: RootCA
				spec: {
					commonName: RootCA
					isCA:       true
					issuerRef: group:      "cert-manager.io"
					issuerRef: kind:       "Issuer"
					issuerRef: name:       SelfSigned
					privateKey: algorithm: "ECDSA"
					privateKey: size:      256
					secretName: RootCA
					subject: organizations: Orgs
				}
			}
			"\(_DBName)-primary-tls": #DatabaseCert & {
				// PGO managed name is "<cluster name>-cluster-cert" e.g. zitadel-cluster-cert
				spec: {
					commonName: "\(_DBName)-primary"
					dnsNames: [
						commonName,
						"\(commonName).\(#TargetNamespace)",
						"\(commonName).\(#TargetNamespace).svc",
						"\(commonName).\(#TargetNamespace).svc.cluster.local",
						"localhost",
						"127.0.0.1",
					]
					usages: ["digital signature", "key encipherment"]
				}
			}
			"\(_DBName)-repl-tls": #DatabaseCert & {
				spec: {
					commonName: "_crunchyrepl"
					dnsNames: [commonName]
					usages: ["digital signature", "key encipherment"]
				}
			}
			"\(_DBName)-client-tls": #DatabaseCert & {
				spec: {
					commonName: "\(_DBName)-client"
					dnsNames: [commonName]
					usages: ["digital signature", "key encipherment"]
				}
			}
		}
	}
}

#DatabaseCert: #Certificate & {
	metadata: name:      string
	metadata: namespace: #TargetNamespace
	spec: {
		duration:    "2160h" // 90d
		renewBefore: "360h"  // 15d
		issuerRef: group:      "cert-manager.io"
		issuerRef: kind:       "Issuer"
		issuerRef: name:       RootCA
		privateKey: algorithm: "ECDSA"
		privateKey: size:      256
		secretName: metadata.name
		subject: organizations: Orgs
	}
}
