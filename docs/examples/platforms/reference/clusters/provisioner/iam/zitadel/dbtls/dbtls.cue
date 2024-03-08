package holos

// Manage an Issuer for the database.

// Both cockroach and postgres handle tls database connections with cert manager
// PGO: https://github.com/CrunchyData/postgres-operator-examples/tree/main/kustomize/certmanager/certman
// CRDB: https://github.com/cockroachdb/helm-charts/blob/3dcf96726ebcfe3784afb526ddcf4095a1684aea/README.md?plain=1#L196-L201

#InputKeys: component: "dbtls"

let CA = {
	Root: "root-issuer"
	MTLS: "mtls-issuer"
}
let Secret = {
	RootCA: "\(_DBName)-root-ca"
}

#KubernetesObjects & {
	apiObjects: {
		Issuer: {
			caIssuer: #Issuer & {
				_description: "Self signed root ca to issue intermediate ca certs"
				metadata: name:      CA.Root
				metadata: namespace: #TargetNamespace
				spec: selfSigned: {}
			}
			dbIssuer: #Issuer & {
				_description: "Root signed intermediate ca to issue mtls database certs"
				metadata: name:      CA.MTLS
				metadata: namespace: #TargetNamespace
				spec: ca: secretName: Secret.RootCA
			}
		}
		Certificate: {
			caCert: #Certificate & {
				_description: "Root CA cert for database"
				metadata: name:      spec.secretName
				metadata: namespace: #TargetNamespace
				spec: {
					commonName: "root"
					isCA:       true
					issuerRef: group:      "cert-manager.io"
					issuerRef: kind:       "Issuer"
					issuerRef: name:       CA.Root
					privateKey: algorithm: "ECDSA"
					privateKey: size:      256
					secretName: Secret.RootCA
					subject: organizations: ["Database"]
				}
			}
			primary: #DatabaseCert & {
				metadata: name: "\(_DBName)-primary"
				spec: {
					dnsNames: [
						"localhost",
						"127.0.0.1",
						metadata.name,
						"\(metadata.name).\(#TargetNamespace)",
						"\(metadata.name).\(#TargetNamespace).svc",
						"\(metadata.name).\(#TargetNamespace).svc.cluster.local",
						"*.\(metadata.name)",
						"*.\(metadata.name).\(#TargetNamespace)",
						"*.\(metadata.name).\(#TargetNamespace).svc",
						"*.\(metadata.name).\(#TargetNamespace).svc.cluster.local",
						"*.\(_DBName)-pods",
						"*.\(_DBName)-pods.\(#TargetNamespace)",
						"*.\(_DBName)-pods.\(#TargetNamespace).svc",
						"*.\(_DBName)-pods.\(#TargetNamespace).svc.cluster.local",
					]
					usages: ["digital signature", "key encipherment", "server auth", "client auth"]
				}
			}
			repl: #ClientCert & {
				metadata: name: "\(_DBName)-repl"
				// NOTE: This must be the literal value _crunchyrepl for replication to work.
				spec: commonName: "_crunchyrepl"
			}
			client: #ClientCert & {
				metadata: name: "\(_DBName)-client"
			}
		}
	}
}

#DatabaseCert: #Certificate & {
	metadata: name:      string
	metadata: namespace: #TargetNamespace
	spec: {
		commonName:  string | *metadata.name
		duration:    "2160h" // 90d
		renewBefore: "360h"  // 15d
		issuerRef: group:      "cert-manager.io"
		issuerRef: kind:       "Issuer"
		issuerRef: name:       CA.MTLS
		privateKey: algorithm: "RSA"
		privateKey: size:      2048
		secretName: metadata.name
		subject: organizations: ["Database"]
	}
}

#ClientCert: {
	#DatabaseCert
	metadata: name: string
	spec: {
		commonName: _
		dnsNames: [commonName]
		usages: ["digital signature", "key encipherment", "client auth"]
	}
}
