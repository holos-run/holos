# Install Holos

Holos is distributed as a single file executable.

## Releases

Download `holos` from the [releases](https://github.com/holos-run/holos/releases) page and place the executable into your shell path.

## Go install

Alternatively, install directly into your go bin path using:

```shell
go install github.com/holos-run/holos/cmd/holos@latest
```

### What you'll need

- [helm](https://github.com/helm/helm/releases) to fetch and render Helm chart components.
- [kubectl](https://kubernetes.io/docs/tasks/tools/) to [kustomize](https://kustomize.io/) components.

