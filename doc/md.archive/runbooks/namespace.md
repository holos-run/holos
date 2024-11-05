# Namespaces

Holos follows the [Namespace Sameness - Sig Multicluster Position][1].  A
namespace is the same on all clusters within the scope of a platform.

Namespaces are also security boundaries for role based access control.  As such,
permission to read a secret in a namespace means the secret is readable on all
clusters in the platform.

When adding a component to a platform, create a namespace using the following
process.  This ensures a namespace scoped `SecretStore` is created to sync
`ExternalSecret` resources from the management cluster.

 1. Add a new project to the `_Projects` struct in `platform.cue`.
 2. Add the namespace to the `spec.namespaces` field of the project.
 3. Render the platform
 4. Apply the `namespaces` component to the management cluster
 5. Apply the `eso-creds-manager` component to the management cluster to create the `eso-reader` ksa for the namespace `SecretStore`
 6. Get a timestamp: `STAMP="$(date +%s)"`
 7. Run the job to populate ecr creds: `kubectl create job -n holos-system --from=cronjob/ecr-creds-manager ecr-creds-manager-$STAMP`
 8. Wait for the job to complete: `kubectl -n holos-system logs -l job-name=ecr-creds-manager-$STAMP -f`
 9. Apply the `namespaces` component to the workload clusters
 10. On the workload cluster, run the job to fetch the eso-reader creds: `kubectl create job -n holos-system --from=cronjob/eso-creds-refresher eso-creds-refresher-${STAMP}`
 11. Wait for the job to complete: `kubectl -n holos-system logs -l job-name=eso-creds-refresher-${STAMP}`
 12. Apply the secretstores component to the workload cluster.
 13. Apply any other cluster specific components which were modified by the `holos render platform ./platform` command.

Your namespace is created and you have the ability to create secrets in the management cluster and pull them using ExternalSecret resources. (edited)  

[1]: https://github.com/kubernetes/community/blob/dd4c8b704ef1c9c3bfd928c6fa9234276d61ad18/sig-multicluster/namespace-sameness-position-statement.md
