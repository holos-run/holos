package v1alpha1

import corev1 "k8s.io/api/core/v1"

#ConfigMap: corev1.#ConfigMap & {
	apiVersion: "v1"
	kind:       "ConfigMap"
	...
}
