package holos

import (
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"

	ci "cert-manager.io/clusterissuer/v1"
	rgv1 "gateway.networking.k8s.io/referencegrant/v1beta1"
	certv1 "cert-manager.io/certificate/v1"
	hrv1 "gateway.networking.k8s.io/httproute/v1"
	gwv1 "gateway.networking.k8s.io/gateway/v1"
)

#Resources: {
	[Kind=string]: [InternalLabel=string]: {
		kind: Kind
		metadata: name: string | *InternalLabel
	}

	Certificate: [_]:        certv1.#Certificate
	ClusterIssuer: [_]:      ci.#ClusterIssuer
	ClusterRole: [_]:        rbacv1.#ClusterRole
	ClusterRoleBinding: [_]: rbacv1.#ClusterRoleBinding
	ConfigMap: [_]:          corev1.#ConfigMap
	CronJob: [_]:            batchv1.#CronJob
	Deployment: [_]:         appsv1.#Deployment
	HTTPRoute: [_]:          hrv1.#HTTPRoute
	Job: [_]:                batchv1.#Job
	Namespace: [_]:          corev1.#Namespace
	ReferenceGrant: [_]:     rgv1.#ReferenceGrant
	Role: [_]:               rbacv1.#Role
	RoleBinding: [_]:        rbacv1.#RoleBinding
	Service: [_]:            corev1.#Service
	ServiceAccount: [_]:     corev1.#ServiceAccount
	StatefulSet: [_]:        appsv1.#StatefulSet

	Gateway: [_]: gwv1.#Gateway & {
		spec: gatewayClassName: string | *"istio"
	}
}
