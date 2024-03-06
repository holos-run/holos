# Namespaces

This component manages namespaces and service accounts in the provisioner cluster for the purpose of providing secrets to workload clusters.

The `eso-creds-refresher` IAM Service Account accesses the provisioner api server as a User.  A job of the same name runs in the `holos-sytem` namespace of each non-provisioner cluster in the platform.  These clusters are external to the provisioner and so called external clusters.  The job exchanges a kubernetes service account bearer token for an access token of the `eso-creds-refresher` IAM Service Account.  The job then identifies every `eso-reader` and `eso-writer` service account and obtains an access token for each, storing it in a secret in the cluster.

Finally, a SecretStore in the external cluster fetches ExternalSecrets from the provisioner cluster using the token credential stored by the `eso-creds-refresher` Job.
