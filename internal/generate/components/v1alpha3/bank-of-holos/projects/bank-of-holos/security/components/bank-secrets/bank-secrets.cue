package holos

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"
)

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Objects = {
	Name:      "bank-secrets"
	Namespace: #BankOfHolos.Security.Namespace

	Resources: [_]: [_]: metadata: namespace:    Namespace
	Resources: [_]: [ID=string]: metadata: name: string | *ID

	Resources: {
		// Kubernetes ServiceAccount used by the Job.
		ServiceAccount: (Name): corev1.#ServiceAccount
		// Role to allow the ServiceAccount to update secrets.
		Role: (Name): rbacv1.#Role & {
			rules: [{
				apiGroups: [""]
				resources: ["secrets"]
				verbs: ["create", "update", "patch"]
			}]
		}

		RoleBinding: (Name): rbacv1.#RoleBinding & {
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "Role"
				name:     Role[Name].metadata.name
			}
			subjects: [{
				kind:      "ServiceAccount"
				name:      ServiceAccount[Name].metadata.name
				namespace: Namespace
			}]
		}

		let JobSpec = {
			serviceAccountName: Name
			restartPolicy:      "OnFailure"
			securityContext: {
				seccompProfile: type: "RuntimeDefault"
				runAsNonRoot: true
				runAsUser:    8192 // app
			}
			containers: [
				{
					name:  "toolkit"
					image: "quay.io/holos-run/toolkit:2024-09-16"
					securityContext: {
						capabilities: drop: ["ALL"]
						allowPrivilegeEscalation: false
					}
					command: ["/bin/bash"]
					args: ["/config/entrypoint"]
					env: [{
						name:  "HOME"
						value: "/tmp"
					}]
					volumeMounts: [{
						name:      "config"
						mountPath: "/config"
						readOnly:  true
					}]
				},
			]
			volumes: [
				{
					name: "config"
					configMap: name: Name
				},
			]
		}

		Job: (Name): batchv1.#Job & {
			spec: template: spec: JobSpec
		}

		ConfigMap: (Name): corev1.#ConfigMap & {
			data: entrypoint: ENTRYPOINT
		}
	}
}

let ENTRYPOINT = """
	#! /bin/bash
	#

	tmpdir="$(mktemp -d)"
	finish() {
	  status=$?
	  rm -rf "${tmpdir}"
	  return $status
	}
	trap finish EXIT

	set -euo pipefail

	cd "$tmpdir"
	mkdir jwt-key
	cd jwt-key
	echo "generating private key" >&2
	ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key -q -N "" -C bank-of-holos
	echo "generating public key" >&2
	ssh-keygen -e -m PKCS8 -f jwtRS256.key > jwtRS256.key.pub
	cd ..

	echo "copying keys into manifest secret.yaml" >&2
	kubectl create secret generic jwt-key --from-file=jwt-key --dry-run=client -o yaml > secret.yaml

	echo "applying secret" >&2
	kubectl apply --server-side=true -f secret.yaml

	echo "cleaning up" >&2
	rm -rf jwt-key
	rm -f secret.yaml

	echo "ok done" >&2
	"""
