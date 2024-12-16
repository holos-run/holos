# Kargo

https://docs.kargo.io/how-to-guides/installing-kargo

## Holos Integration

1. We generate the admin password and JWT signing key using an ESO generator.
2. Credentials are stored in Secret `admin-credentials`
3. Helm values are imported to CUE.

## Admin Password

Get the generated admin password and access Kargo at
https://kargo.holos.localhost/

```bash
kubectl get secret -n kargo admin-credentials -o json \
  | jq --exit-status -r '.data.password | @base64d' \
  | pbcopy
```

## Helm Values

```bash
helm inspect values oci://ghcr.io/akuity/kargo-charts/kargo > values.yaml
holos cue import --package holos --path 'Kargo: Values:' --outfile values.cue ./values.yaml
```
