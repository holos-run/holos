package holos

import "encoding/json"

#DependsOn: _ESO

#InputKeys: {
	project:   "secrets"
	component: "eso-creds-refresher"
}

#TargetNamespace: #CredsRefresher.namespace

// output kubernetes api objects for holos
#KubernetesObjects & {
	apiObjects: {
		for obj in #CredsRefresherService.objects {
			let Kind = obj.kind
			let Name = obj.metadata.name
			"\(Kind)": "\(Name)": obj
		}
	}
}

let NAME = #CredsRefresher.name
let AUD = "//iam.googleapis.com/projects/\(#InputKeys.gcpProjectNumber)/locations/global/workloadIdentityPools/holos/providers/k8s-\(#InputKeys.cluster)"
let MOUNT = "/var/run/service-account"
let EMAIL = #CredsRefresher.iamServiceAccount
let REGION = #CredsRefresher.region

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
KUBECONFIG=${HOME}/kubeconfig.provisioner gcloud container clusters get-credentials provisioner --region=\(REGION)

# Get a list of the service accounts to issue tokens for.
provisioner get serviceaccount -A --selector=holos.run/job.name=\(NAME) --output=json > serviceaccounts.json

# Create the tokens
mkdir tokens
jq -r '.items[].metadata | "provisioner -n \\(.namespace) create token --duration=12h \\(.name) > tokens/\\(.namespace).\\(.name).jwt"' serviceaccounts.json | bash -x

# Create the secrets
mksecret tokens/*.jwt

# Apply the secrets to the local cluster
kubectl apply --server-side=true -f secrets.yaml
"""

// #CredsRefresherService defines the job that refreshes credentials used by eso SecretStore resources.
#CredsRefresherService: {
	objects: [
		// Kubernetes ServiceAccount used by the Job.
		#ServiceAccount & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
				annotations: "holos.run/description": "Refreshes provisioner cluster eso-reader and eso-writer service account credentials."
			}
		},
		// Job needs to read and write secrets across the cluster.
		#ClusterRole & {
			metadata: {
				name: #CredsRefresher.name
			}
			rules: [
				{
					apiGroups: [""]
					resources: ["secrets"]
					verbs: ["*"]
				},
			]
		},
		// Bind the Role to the ServiceAccount for the Job.
		#ClusterRoleBinding & {
			metadata: {
				name: #CredsRefresher.name
			}
			roleRef: {
				apiGroup: "rbac.authorization.k8s.io"
				kind:     "ClusterRole"
				name:     #CredsRefresher.name
			}
			subjects: [
				{
					kind:      "ServiceAccount"
					name:      NAME
					namespace: #CredsRefresher.namespace
				},
			]
		},
		#Job & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
			}
			spec: template: spec: #PodSpec
		},
		#CronJob & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
			}
			spec: {
				schedule: "0 */8 * * *"
				jobTemplate: spec: {
					template: spec: #PodSpec
					backoffLimit: 20
				}
			}
		},
		#ConfigMap & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
			}
			data: {
				entrypoint:                      ENTRYPOINT
				provisioner:                     PROVISIONER_SCRIPT
				mksecret:                        MKSECRET
				"credential-configuration.json": json.Marshal(CREDCONFIG)
			}
		},
	]
}

// #PodSpec is the pod spec field of the eso-creds-refresher job
#PodSpec: {
	serviceAccountName: NAME
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
					value: #CredsRefresher.project
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
			configMap: name: "eso-creds-refresher"
		},
	]
}
