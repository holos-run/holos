# IAM

The IAM service provides identity and access management for a holos managed platform.  Zitadel is the identity provider which integrates tightly with:

 1. AuthorizationPolicy at the level of the service mesh.
 2. Application level oidc login (ArgoCD, Grafana, etc...)
 3. Cloud provider IAM via oidc.
