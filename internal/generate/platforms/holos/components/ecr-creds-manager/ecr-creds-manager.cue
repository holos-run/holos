package holos

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"
)

let NAME = "ecr-creds-manager"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

// The path Pod Identity uses.
let MOUNT = "/var/run/secrets/eks.amazonaws.com/serviceaccount/"
let AWS_ACCOUNT = _Platform.Model.aws.accountNumber
let AWS_REGION = _Platform.Model.aws.primaryRegion
let AWS_ROLE_ARN = "arn:aws:iam::\(AWS_ACCOUNT):role/\(NAME)"

let Objects = {
	Name:      NAME
	Namespace: "holos-system"

	Resources: {
		// Kubernetes ServiceAccount used by the Job.
		ServiceAccount: "\(Name)": corev1.#ServiceAccount & {
			metadata: {
				name:      Name
				namespace: Namespace
				annotations: "holos.run/description": "Refreshes image pull credentials for use with AWS ECR."
				// annotations: "eks.amazonaws.com/role-arn": AWS_ROLE_ARN
			}
		}

		// Job needs to read and write secrets across the cluster.
		ClusterRole: "\(Name)": rbacv1.#ClusterRole & {
			metadata: name: Name
			rules: [
				{
					apiGroups: [""]
					resources: ["secrets"]
					verbs: ["*"]
				},
				{
					apiGroups: [""]
					resources: ["namespaces"]
					verbs: ["list"]
				},
			]
		}

		// Bind the Role to the ServiceAccount for the Job.
		ClusterRoleBinding: "\(Name)": rbacv1.#ClusterRoleBinding & {
			metadata: name: Name
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "ClusterRole"
				name:     Name
			}
			subjects: [
				{
					kind:      "ServiceAccount"
					name:      Name
					namespace: Namespace
				},
			]
		}

		// Make the CronJob and Job identical.
		let JobSpec = {
			serviceAccountName: Name
			restartPolicy:      "OnFailure"
			securityContext: {
				seccompProfile: type: "RuntimeDefault"
				runAsNonRoot: true
				runAsUser:    8192 // app
			}
			nodeSelector: {
				"cloud.google.com/gke-spot": "true"
				"kubernetes.io/os":          "linux"
			}
			containers: [
				{
					name:  "toolkit"
					image: "quay.io/holos/toolkit:latest"
					securityContext: {
						capabilities: drop: ["ALL"]
						allowPrivilegeEscalation: false
					}
					command: ["/bin/bash"]
					args: ["/config/entrypoint"]
					env: [
						{
							name:  "HOME"
							value: "/tmp"
						},
						{
							name:  "AWS_DEFAULT_REGION"
							value: AWS_REGION
						},
						{
							name:  "AWS_REGION"
							value: AWS_REGION
						},
						{
							name:  "AWS_ROLE_ARN"
							value: AWS_ROLE_ARN
						},
						{
							name:  "AWS_WEB_IDENTITY_TOKEN_FILE"
							value: MOUNT + "token"
						},
						{
							name:  "AWS_STS_REGIONAL_ENDPOINTS"
							value: "regional"
						},
					]
					volumeMounts: [
						{
							name:      "config"
							mountPath: "/config"
							readOnly:  true
						},
						{
							name:      "aws-token"
							mountPath: MOUNT
							readOnly:  true
						},
					]
				},
			]
			volumes: [
				{
					name: "config"
					configMap: name: Name
				},
				{
					name: "aws-token"
					projected: sources: [{
						serviceAccountToken: {
							path:              "token"
							expirationSeconds: 3600
							audience:          "sts.amazonaws.com"
						}
					}]
				},
			]
		}

		Job: "\(Name)": batchv1.#Job & {
			metadata: {
				name:      Name
				namespace: Namespace
			}
			spec: template: spec: JobSpec
		}

		CronJob: "\(Name)": batchv1.#CronJob & {
			metadata: name:      Name
			metadata: namespace: Namespace
			spec: {
				schedule: "0 */8 * * *"
				jobTemplate: spec: {
					template: spec: JobSpec
					backoffLimit: 20
				}
			}
		}

		ConfigMap: "\(Name)": corev1.#ConfigMap & {
			metadata: name:      Name
			metadata: namespace: Namespace
			data: entrypoint:    ENTRYPOINT
		}
	}
}

let ENTRYPOINT = """
	#! /bin/bash
	#

	tmpdir="$(mktemp -d)"
	finish() {
		rm -rf "${tmpdir}"
	}
	trap finish EXIT

	set -xeuo pipefail

	aws sts get-caller-identity

	aws ecr get-login-password --region \(AWS_REGION) \\
		| docker login --username AWS --password-stdin \(AWS_ACCOUNT).dkr.ecr.\(AWS_REGION).amazonaws.com

	kubectl create secret docker-registry ecr-creds-\(AWS_ACCOUNT) \\
		--from-file=.dockerconfigjson=${HOME}/.docker/config.json \\
		--dry-run=client -o yaml \\
			> "${tmpdir}/secret.yaml"

	# Copy the secret to all namespaces
	for ns in $(kubectl -o=jsonpath='{.items[*].metadata.name}' get namespaces); do
		echo -n "Copying secret to namespace ${ns}: "
		kubectl -n $ns apply --server-side=true -f "${tmpdir}/secret.yaml" || continue
		echo "Usage: "kubectl -n $ns patch serviceaccount default -p "'"'{"imagePullSecrets": [{"name": "ecr-creds-\(AWS_ACCOUNT)"}]}'"'"
	done
	"""
