# ESO Credentials Manager

The `eso-creds-manager` component manages namespaces and service accounts in the Management Cluster for the purpose of making secrets available to workload clusters via the external secrets operator.

The `eso-creds-refresher` IAM Service Account accesses the Management Cluster api server as a kubernetes `User`.  In workload clusters, a `Job` named `eso-creds-refresher` runs in the `holos-sytem` namespace.  Refer to the `eso-creds-refresher` Holos component for the configuration of this job.

The `Job` exchanges a local kubernetes service account token for an access token of the `eso-creds-refresher` IAM Google Service Account.  The job then identifies every `eso-reader` and `eso-writer` service account and obtains an access token for each, storing it in a secret in the cluster.

Finally, a SecretStore in the external cluster fetches ExternalSecrets from the provisioner cluster using the token credential stored by the `eso-creds-refresher` Job.
