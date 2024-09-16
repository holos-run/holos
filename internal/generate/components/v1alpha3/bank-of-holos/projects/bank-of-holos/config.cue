package holos

import (
	core "k8s.io/api/core/v1"
	es "external-secrets.io/externalsecret/v1beta1"
	ss "external-secrets.io/secretstore/v1beta1"
)

let BankName = #BankOfHolos.Name

#BankOfHolos: {
	// Resources to make available in each of the project namespaces.
	Resources: {
		ServiceAccount: (BankName): core.#ServiceAccount & {
			apiVersion: "v1"
			kind:       "ServiceAccount"
			metadata: name: BankName
		}

		// SecretStore to fetch secrets owned by the security team
		SecretStore: (BankName): ss.#SecretStore & {
			metadata: name: #BankOfHolos.Security.Namespace
			spec: provider: {
				kubernetes: {
					remoteNamespace: #BankOfHolos.Security.Namespace
					auth: serviceAccount: name: ServiceAccount[BankName].metadata.name
					server: {
						url: "https://kubernetes.default.svc"
						caProvider: {
							type: "ConfigMap"
							name: "kube-root-ca.crt"
							key:  "ca.crt"
						}
					}
				}
			}
		}

		// We do not check the private key into version control.
		// https://github.com/GoogleCloudPlatform/bank-of-anthos/tree/v0.6.5/extras/jwt
		ExternalSecret: "jwt-key": es.#ExternalSecret & {
			metadata: name: "jwt-key"
			spec: {
				target: name: metadata.name
				dataFrom: [{extract: {key: metadata.name}}]
				refreshInterval: "5s"
				secretStoreRef: kind: "SecretStore"
				secretStoreRef: name: SecretStore[BankName].metadata.name
			}
		}

		// https://github.com/GoogleCloudPlatform/bank-of-anthos/blob/release/v0.6.5/kubernetes-manifests/config.yaml
		ConfigMap: "environment-config": core.#ConfigMap & {
			apiVersion: "v1"
			kind:       "ConfigMap"
			metadata: name: "environment-config"
			data: {
				LOCAL_ROUTING_NUM: "883745000"
				PUB_KEY_PATH:      "/tmp/.ssh/publickey"
			}
		}

		ConfigMap: "service-api-config": core.#ConfigMap & {
			apiVersion: "v1"
			kind:       "ConfigMap"
			metadata: name: "service-api-config"
			data: {
				TRANSACTIONS_API_ADDR: "ledgerwriter.\(#BankOfHolos.Backend.Namespace).svc:8080"
				BALANCES_API_ADDR:     "balancereader.\(#BankOfHolos.Backend.Namespace).svc:8080"
				HISTORY_API_ADDR:      "transactionhistory.\(#BankOfHolos.Backend.Namespace).svc:8080"
				CONTACTS_API_ADDR:     "contacts.\(#BankOfHolos.Backend.Namespace).svc:8080"
				USERSERVICE_API_ADDR:  "userservice.\(#BankOfHolos.Backend.Namespace).svc:8080"
			}
		}

		ConfigMap: "demo-data-config": core.#ConfigMap & {
			apiVersion: "v1"
			kind:       "ConfigMap"
			metadata: name: "demo-data-config"
			data: {
				USE_DEMO_DATA:       "True"
				DEMO_LOGIN_USERNAME: "testuser"
				// All demo user accounts are hardcoded to use the login password 'bankofanthos'
				DEMO_LOGIN_PASSWORD: "bankofanthos"
			}
		}
	}
}
