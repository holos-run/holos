package holos

import (
	corev1 "k8s.io/api/core/v1"
)

_NamespaceObject: {
	metadata: name:      string
	metadata: namespace: string
	metadata: labels: "app.holos.run/managed": "true"
}

#ConfigMap: _NamespaceObject & corev1.#ConfigMap
