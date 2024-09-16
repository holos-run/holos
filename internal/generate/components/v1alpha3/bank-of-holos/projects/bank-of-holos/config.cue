package holos

import core "k8s.io/api/core/v1"

#BankOfHolos: {
	// Resources to make available in each of the project namespaces.
	Resources: {
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
				TRANSACTIONS_API_ADDR: "ledgerwriter:8080"
				BALANCES_API_ADDR:     "balancereader:8080"
				HISTORY_API_ADDR:      "transactionhistory:8080"
				CONTACTS_API_ADDR:     "contacts:8080"
				USERSERVICE_API_ADDR:  "userservice:8080"
			}
		}

		ConfigMap: "demo-data-config": core.#ConfigMap & {
			apiVersion: "v1"
			kind:       "ConfigMap"
			metadata: name: "demo-data-config"
			data: {
				USE_DEMO_DATA:       "True"
				DEMO_LOGIN_USERNAME: "testuser"
				// All demo user accounts are hardcoded to use the login password 'bankofholos'
				DEMO_LOGIN_PASSWORD: "bankofholos"
			}
		}

		ServiceAccount: "bank-of-holos": core.#ServiceAccount & {
			apiVersion: "v1"
			kind:       "ServiceAccount"
			metadata: name: "bank-of-holos"
		}
	}
}
