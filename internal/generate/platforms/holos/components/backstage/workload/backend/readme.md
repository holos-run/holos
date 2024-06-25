# Backstage Secrets

Backstage needs secrets in place on the management cluster to operate.

ExternalSecret `github-app-credentials` generated using:

Refer to the [portal](https://github.com/holos-run/portal) repo.

```sh
# portal is your backstage repository created with `npx @backstage/create-app`
cd portal

# my-org is your github organization.
yarn backstage-cli create-github-app my-org

# Create the secret in your management cluster.
mv github-app-backstage-*-credentials.yaml github-app-credentials.yaml
holos create secret -n backstage --append-hash=false --from-file=github-app-credentials.yaml github-app-credentials

# Remove the secret from the local host.
rm -f github-app-credentials.yaml
```
