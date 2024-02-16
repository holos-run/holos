# Provisioner Secrets

This document describes how to deploy a service to store and provide secrets for workload clusters in a Holos platform.  The document walks through setting up a Provisioner Cluster, which is an api server providing a limited set of resources for the purpose of provisioning multiple Workload Clusters.  The limited set of resources are:

1. Secrets to store sensitive data workload clusters need.
2. Namespaces for security boundaries.
3. ServiceAccounts for automated and secure access control.
4. RoleBindings for role based access control.

Requirements:

The provisioner service must:

1. Store and provide Secrets for workload clusters.
2. Enable automated provisioning of workload clusters.
3. Globally accessible.
4. Highly available over multiple availability zones in one region.
5. Provide accessible audit logs for read and write operates over the Secrets.
6. Provide granular role based access control to at least the Namespace boundary.
7. Enable automated access from workload clusters running anywhere.
8. Store at least 4,000 secrets up to 1MiB each.

A control plane only gke autopilot cluster meets all of these requirements in a cost-efficient way.

This document is broken down into the steps required each time a new Platform, Cluster, or Namespace is added.

## References

Refer to [Configure workload identity federation with Kubernetes](https://cloud.google.com/iam/docs/workload-identity-federation-with-kubernetes#kubernetes) which provides a step-by-step guide this document reflects.

For additional information refer to:

- [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation) enables a kubernetes service account outside of Google Cloud to authenticate and access secrets in the provisioner cluster.
- [On-premises or another cloud provider](https://cloud.google.com/docs/authentication/provide-credentials-adc#on-prem)

The preferred way to authenticate with Google Cloud using credentials from a different IdP is to use workload identity federation; you create a credential configuration file and set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to it.

For help with setting up workload identity federation for ADC, see [Configuring workload identity federation](https://cloud.google.com/iam/docs/configuring-workload-identity-federation) and [Obtaining short-lived credentials with identity federation](https://cloud.google.com/iam/docs/workload-identity-federation-with-other-clouds#generate-automatic).

## Project

Create a project for Holos IAM and Secrets management.  This project will have no workloads, one cluster api server designated the Provisioner Cluster for secrets, rbac, and audit logging.  THe project will also have multiple IAM service accounts for workload identity federation using the security token service.

Billing is necessary, expected costs for a production setup are expected to be less than one usd per month.

```shell
PROJECT_ID="$(gcloud config get-value project)"
PROJECT_NUMBER="$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')"
ORG_DOMAIN="example.com"
```

## Seed Cluster

```shell
gcloud container clusters create-auto provisioner \
  --release-channel=stable \
  --monitoring=SYSTEM \
  --logging=SYSTEM \
  --enable-master-global-access \
  --enable-master-authorized-networks \
  --master-authorized-networks=0.0.0.0/0 \
  --region=us-central1
```

Optional [Groups for RBAC](https://cloud.google.com/kubernetes-engine/docs/how-to/google-groups-rbac).

```shell
--security-group=gke-security-groups@$ORG_DOMAIN
```

## Access

```shell
mkdir -p ~/.holos
KUBECONFIG=${HOME}/.holos/kubeconfig.provisioner gcloud container clusters get-credentials provisioner --region=us-central1
```

## Platform Setup - Workload Identity

This section enables workload clusters to automatically fetch secrets when the workload cluster is provisioned.  The `eso-creds-refresher` Job in the workload cluster provides credentials for the External Secrets Operator to fetch secrets from the Provisioner Cluster.

Execute these steps once for each new platform.

### Identity Pool

```shell
gcloud iam workload-identity-pools create holos \
  --location="global" \
  --display-name="Holos Clusters" \
  --description="Authenticates Holos workload cluster service account id tokens"
```

### IAM Service Account

Create a service account to refresh eso credentials.  A job running on each cluster uses this iam service account to refresh credentials ESO uses to sync secrets from the provisioner cluster.

NOTE: IAM permissions are namespace scoped, not cluster scoped.  They're necessary once for all clusters in the identity pool.

```shell
gcloud iam service-accounts create eso-creds-refresher \
  --display-name="ESO Credential Refresher" \
  --description="Used by eso-creds-refresher job in each workload cluster to refresh the creds ESO needs to sync secrets from the provisioner cluster."
```

Allow the `eso-creds-refresher` Kubernetes service account in the `holos-system` namespace of workload clusters to impersonate the IAM service account.

```shell
gcloud iam service-accounts add-iam-policy-binding \
  eso-creds-refresher@$PROJECT_ID.iam.gserviceaccount.com \
  --member="principal://iam.googleapis.com/projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/holos/subject/system:serviceaccount:holos-system:eso-creds-refresher" \
  --role=roles/iam.workloadIdentityUser
```

Allow the IAM service account to authenticate to the provisioner cluster.

```shell
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:eso-creds-refresher@$PROJECT_ID.iam.gserviceaccount.com" \
  --role=roles/container.clusterViewer
```

### Role Binding

Grant eso-creds-refresher access to the provisioner cluster.

```shell
KUBECONFIG=${HOME}/.holos/kubeconfig.provisioner kubectl 
```

### Summary

This section configured:

1. Provisioner Cluster to provide secrets to workload clusters.
2. IAM service account `eso-creds-refresher` to identify the credential refresher job.
3. Workload identity pool to authenticate the `eso-creds-refresher` Kubernetes service account in an external cluster.
4. IAM policy to allow `eso-creds-refresher` to authenticate to the Provisioner Cluster.
5. RoleBinding to allow `eso-creds-refresher` to create kubernetes service account tokens representing the credentials for use by SecretStore resources in workload clusters.

## Cluster Setup

The goal of this section is to deploy the `eso-creds-refresher` Job into a workload cluster which will fetch credentials for SecretStore resources to access secrets in the Provisioner Cluster.

Execute once for each cluster in the platform whenever the cluster json web key set (jwks) changes.

GKE and EKS clusters always get new key on creation.  A Talos cluster can preserve the issuer jwks if the same initial secret is applied after a reset.

```shell
HOLOS_CLUSTER_NAME=west1
ISSUER_URL="https://example.com/clusters/${HOLOS_CLUSTER_NAME}"
```

```shell
gcloud iam workload-identity-pools providers create-oidc \
  k8s-$HOLOS_CLUSTER_NAME \
  --workload-identity-pool=holos \
  --location=global \
  --display-name="$HOLOS_CLUSTER_NAME cluster" \
  --description="Trust service account id tokens issued by the cluster api server" \
  --attribute-mapping="google.subject=assertion.sub,attribute.pod=assertion['kubernetes.io']['pod']['name'],attribute.service_account_name=assertion['kubernetes.io']['serviceaccount']['name'],attribute.uid=assertion['kubernetes.io']['serviceaccount']['uid']" \
  --issuer-uri="$ISSUER_URL"
```

## Credential configuration

Refer to [Deploy the kubernetes workload](https://cloud.google.com/iam/docs/workload-identity-federation-with-kubernetes#kubernetes).  See also [AIP-4117](https://google.aip.dev/auth/4117).

The Cloud Client Libraries, the gcloud CLI, and Terraform, can automatically obtain external credentials, and use these credentials to impersonate a service account. To let libraries and tools complete this process, you have to provide a credential configuration file. This file defines the following:

- Where to obtain external credentials from
- Which workload identity pool and provider to use
- Which service account to impersonate

The pod running in the holos manages cluster will authenticate to Google using the following credential config. This
needs to be executed in the external workload.

```shell
gcloud iam workload-identity-pools create-cred-config \
    projects/$PROJECT_NUMBER/locations/global/workloadIdentityPools/holos/providers/k8s-$HOLOS_CLUSTER_NAME \
    --service-account=eso-creds-refresher@$PROJECT_ID.iam.gserviceaccount.com \
    --credential-source-file=/var/run/service-account/token \
    --credential-source-type=text \
    --output-file=gcloud/credential-configuration.json
```

## Additional config

The `eso-creds-refresher` Pod needs some configuration.

```shell
env PROJECT_ID=$PROJECT_ID PROJECT_NUMBER=$PROJECT_NUMBER HOLOS_CLUSTER_NAME=$HOLOS_CLUSTER_NAME \
  envsubst < gcloud/pod.template.yaml > gcloud/pod.yaml
```

## Test authentication

```shell
kubectl apply -k gcloud
kubectl exec gcloud -it -- bash
```

Then in the pod:

```shell
gcloud auth print-access-token
```

## Cluster Credentials



## Audit Logs

Log admin read operations on the `container.googleapis.com` service to log when secrets are read from the cluster api control plane.  This logs when an admin uses `kubectl get secret` to read a secret.

```shell
PROJECT_ID=$(gcloud config get-value project)
gcloud projects get-iam-policy $PROJECT_ID --format=json \
  | jq '.auditConfigs |= map(select(.service | contains("container.googleapis.com") | not)) | .auditConfigs += [{"auditLogConfigs":[{"logType":"ADMIN_READ"}],"service":"container.googleapis.com"}]' \
  | gcloud projects set-iam-policy $PROJECT_ID
```

## Principal Identifiers

Refer to [Service account impersonation](https://cloud.google.com/iam/docs/workload-identity-federation#impersonation)

```shell
# Single identity
principal://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/subject/SUBJECT_ATTRIBUTE_VALUE 
# All identities in a group
principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/group/GROUP_ID 
# All identities with a specific attribute value
principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/attribute.ATTRIBUTE_NAME/ATTRIBUTE_VALUE 
```
