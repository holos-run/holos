# ESO Credentials Refresher

This component manages the `eso-creds-refresher` job which uses workload identity to impersonate an iam service account to access the provisioner cluster.  The job creates tokens for the `eso-reader` provisioner cluster service account and stores them in a secret in the same namespace in the workload cluster for use by the default SecretStore.
