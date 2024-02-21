package holos

import "encoding/json"

// objects are kubernetes api objects to apply
objects: #CredsRefresherService.objects

// output kubernetes api objects for holos
{} & #KubernetesObjects

#InputKeys: {
	project:   "secrets"
	component: "eso-creds-refresher"
}

#TargetNamespace: #CredsRefresher.namespace

let SA_NAME = #CredsRefresher.name
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

let ENTRYPOINT = """
#! /bin/bash
#
set -xeuo pipefail

cd "$HOME"

# Try and get secrets from the local cluster.
kubectl get secrets -A

gcloud config set disable_usage_reporting true
gcloud auth login --cred-file $GOOGLE_APPLICATION_CREDENTIALS

export KUBECONFIG="${HOME}/kubeconfig.provisioner"

# Log into k8s.
gcloud container clusters get-credentials provisioner --region=\(REGION)

# Check access.
kubectl version --output=yaml

#Try and get service accounts from the provisioner cluster
kubectl get serviceaccounts -A

sleep 3600
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
					name:      SA_NAME
					namespace: #CredsRefresher.namespace
				},
			]
		},
		#Pod & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
			}
			spec: #PodSpec
		},
		#ConfigMap & {
			metadata: {
				name:      #CredsRefresher.name
				namespace: #CredsRefresher.namespace
			}
			data: {
				entrypoint:                      ENTRYPOINT
				"credential-configuration.json": json.Marshal(CREDCONFIG)
			}
		},
	]
}

// #PodSpec is the pod spec field of the eso-creds-refresher job
#PodSpec: {
	serviceAccountName: SA_NAME
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
