# Local Development

This document captures notes on locally developing Holos.

Follow the steps in [Try Holos Locally](../guides/try-holos), but take care
to select `Develop` tabs when creating the k3d cluster so you have a local
registry to push to.

## Apply Resources

Work will be done in the `dev-holos` namespace.

Apply the infrastructure, which should persist when tilt is started / stopped.

```bash
kubectl apply --server-side=true -f ./hack/tilt/k8s/dev-holos-infra
```

This creates the PostgresCluster, service account, etc...

## Start tilt

Tilt will build the go executable, build the container, then push it to the
local repository associated with k3d.

```bash
./hack/tilt/bin/tilt up
```
