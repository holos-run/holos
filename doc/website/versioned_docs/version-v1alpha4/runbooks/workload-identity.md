# Workload Identity

When a new workload cluster is provisioned, allow it to access the Management
Cluster using workload identity.  This is necessary for the
`eso-creds-refresher` component and `Job` that executes in each workload
cluster, which in turn enables the `SecretStore` in each namespace to sync
secrets.

Build the cluster with Cluster API.
See https://github.com/holos-run/holos-infra/blob/main/hack/capi/eks/aws2/aws2-managedmachinepool.yaml#L81-L84

## Workload Identity Provider
Add the Cluster as a workload identity provider to the `holos-ops` gcp project.

Pool: [holos](https://console.cloud.google.com/iam-admin/workload-identity-pools/pool/holos?organizationId=358674006047&project=holos-ops)
Name: `k8s-aws1`, `k8s-aws2`, etc...
### Issuer URL:
```
kubectl create -n default token default | cut -d. -f2 | base64 -d | jq -r .iss
```

### Audience
Use the default audience.
### Attribute Mapping

| Google                           | OIDC                                                   |
| -------------------------------- | ------------------------------------------------------ |
| `google.subject`                 | `assertion.sub`                                        |
| `attribute.service_account_name` | `assertion['kubernetes.io']['serviceaccount']['name']` |
| `attribute.uid`                  | `assertion['kubernetes.io']['serviceaccount']['uid']`  |
| `attribute.pod`                  | `assertion['kubernetes.io']['pod']['name']`            |
