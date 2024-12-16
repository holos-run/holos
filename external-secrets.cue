@if(!NoExternalSecrets)
package holos

import (
	es "external-secrets.io/externalsecret/v1beta1"
	ss "external-secrets.io/secretstore/v1beta1"
	pw "generators.external-secrets.io/password/v1alpha1"
)

ExternalSecrets: {
	Version: "0.10.7"
}

#Resources: {
	ExternalSecret?: [_]: es.#ExternalSecret
	Password?: [_]:       pw.#Password
	SecretStore?: [_]:    ss.#SecretStore
}
