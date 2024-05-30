# IAM

The IAM service provides identity and access management for a holos managed
platform.  Zitadel is the identity provider which integrates tightly with:

 1. AuthorizationPolicy at the level of the service mesh.
 2. Application level oidc login (ArgoCD, Grafana, etc...)
 3. Cloud provider IAM via oidc.

## Preflight

The zitadel master key needs to have a data key named `masterkey` with a Secret
name of `zitadel-masterkey`.

```bash
holos create secret zitadel-masterkey --namespace zitadel --append-hash=false --data-stdin <<EOF
{"masterkey":"$(tr -dc A-Za-z0-9 </dev/urandom | head -c 32)"}
EOF
```
