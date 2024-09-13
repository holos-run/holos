package holos

import (
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	batchv1 "k8s.io/api/batch/v1"
)

let NAME = "eso-creds-refresher"

// Produce a kubernetes objects build plan.
(#Kubernetes & Objects).BuildPlan

let Objects = {
	Name:      NAME
	Namespace: "holos-system"

	Resources: {
		// Kubernetes ServiceAccount used by the Job.
		ServiceAccount: "\(Name)": corev1.#ServiceAccount & {
			metadata: {
				name:      Name
				namespace: Namespace
				annotations: "holos.run/description": "Refreshes provisioner cluster eso-reader and eso-writer service account credentials."
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
				runAsUser:    1000
			}
			containers: [
				{
					name:  "gcloud"
					image: "gcr.io/google.com/cloudsdktool/google-cloud-cli:latest"
					securityContext: {
						capabilities: drop: ["ALL"]
						allowPrivilegeEscalation: false
					}
					command: ["/bin/bash"]
					args: ["/config/entrypoint"]
					env: [
						{
							name:  "GOOGLE_APPLICATION_CREDENTIALS"
							value: "/config/credential-configuration.json"
						},
						{
							name:  "CLOUDSDK_CORE_PROJECT"
							value: _Platform.Model.eso.gcpProjectID
						},
					]
					volumeMounts: [
						{
							name:      "config"
							mountPath: "/config"
							readOnly:  true
						},
						{
							name:      "token"
							mountPath: MOUNT
							readOnly:  true
						},
					]
				},
			]
			volumes: [
				{
					name: "token"
					projected: {
						sources: [
							{
								serviceAccountToken: {
									path:              "token"
									expirationSeconds: 3600
									audience:          AUD
								}
							},
						]
					}
				},
				{
					name: "config"
					configMap: name: Name
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
			data: {
				entrypoint:                      ENTRYPOINT
				provisioner:                     PROVISIONER_SCRIPT
				mksecret:                        MKSECRET
				"credential-configuration.json": json.Marshal(CREDCONFIG)
			}
		}
	}
}

// This is the default audience when configuring a google workload identity provider
let AUD = "//iam.googleapis.com/projects/\(ProjectNumber)/locations/global/workloadIdentityPools/holos/providers/k8s-\(_ClusterName)"
let ProjectNumber = _Platform.Model.eso.gcpProjectNumber
let SecretStoreCluster = _Platform.Model.eso.gkeClusterName
let MOUNT = "/var/run/service-account"
let EMAIL = _Platform.Model.eso.gcpServiceAccount

let REGION = _Platform.Model.eso.gkeRegion

// Refer to https://google.aip.dev/auth/4117
let CREDCONFIG = {
	type:               "external_account"
	audience:           AUD
	subject_token_type: "urn:ietf:params:oauth:token-type:jwt"
	token_url:          "https://sts.googleapis.com/v1/token"
	credential_source: {
		file: MOUNT + "/token"
		format: type: "text"
	}
	service_account_impersonation_url: "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/\(EMAIL):generateAccessToken"
}

let PROVISIONER_SCRIPT = """
	#! /bin/bash
	export KUBECONFIG="${HOME}/kubeconfig.provisioner"
	exec kubectl "$@"
	"""

let MKSECRET = """
	#! /bin/bash
	#
	# Make a secret from a jwt file
	set -euo pipefail
	echo -n > secrets.yaml
	while [[ $# -gt 0 ]]; do
	  file="$1"
	  shift
	  base="$(basename $file)"
	  pair="${base%.jwt}"
	  namespace="${pair%.*}"
	  name="${pair#*.}"
	  # Refer to https://external-secrets.io/latest/api/spec/#external-secrets.io/v1beta1.KubernetesAuth
	  echo "---" >> secrets.yaml
	  kubectl create secret generic -n $namespace $name --from-file=token=${file} --dry-run=client -o yaml >> secrets.yaml
	done
	"""

let ENTRYPOINT = """
#! /bin/bash
#
set -euo pipefail

cd "$HOME"
mkdir -p "${HOME}/bin"
curl -L -o "${HOME}/bin/jq" https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64
chmod 755 "${HOME}/bin/jq"
export PATH="${HOME}/bin:${PATH}"

# script to talk to the provisioner cluster, kubectl is for the local cluster.
install -m 0755 /config/provisioner $HOME/bin/provisioner
install -m 0755 /config/mksecret $HOME/bin/mksecret

# Log into gcp using workload identity iam service account impersonation
gcloud config set disable_usage_reporting true
gcloud auth login --cred-file $GOOGLE_APPLICATION_CREDENTIALS

# Log into the provisioner cluster using the iam service account
KUBECONFIG=${HOME}/kubeconfig.provisioner gcloud container clusters get-credentials \(SecretStoreCluster) --region=\(REGION)

# Get a list of the reader service accounts to issue tokens for.
provisioner get serviceaccount -A --selector=holos.run/job.name=\(NAME) --selector=holos.run/iam.role=reader --output=json > serviceaccounts.json

# Create the tokens
mkdir tokens

kubectl get namespaces -o name > namespaces.txt

echo "Provisioning eso-reader tokens ..." >&2

# Iterate over local namespaces
while IFS= read -r NAMESPACE; do
  echo "Checking local cluster $NAMESPACE ..." >&2
  jq -r '.items[] | select("namespace/"+.metadata.namespace == "'${NAMESPACE}'") | .metadata | "provisioner -n \\(.namespace) create token --duration=12h \\(.name) > tokens/\\(.namespace).\\(.name).jwt"' serviceaccounts.json | bash -x
done < namespaces.txt

# Create the secrets
mksecret tokens/*.jwt

# Apply the secrets to the local cluster
kubectl apply --server-side=true -f secrets.yaml

## Do it again for namespaces which have the eso writer label.
echo "Provisioning eso-writer tokens for local namespaces with selector holos.run/eso.writer=true ..." >&2
kubectl get namespaces --selector=holos.run/eso.writer=true -o name > namespaces.txt
if [[ ! -s namespaces.txt ]]; then
  echo "No namespaces have the holos.run/eso.writer=true label, exiting." >&2
  exit 0
fi

rm -rf tokens
mkdir tokens

provisioner get serviceaccount -A --selector=holos.run/job.name=\(NAME) --selector=holos.run/iam.role=writer --output=json > serviceaccounts.json

while IFS= read -r NAMESPACE; do
  echo "Checking local cluster $NAMESPACE ..." >&2
  jq -r '.items[] | select("namespace/"+.metadata.namespace == "'${NAMESPACE}'") | .metadata | "provisioner -n \\(.namespace) create token --duration=12h \\(.name) > tokens/\\(.namespace).\\(.name).jwt"' serviceaccounts.json | bash -x
done < namespaces.txt

# Create the secrets
mksecret tokens/*.jwt

# Apply the secrets to the local cluster
kubectl apply --server-side=true -f secrets.yaml
"""
