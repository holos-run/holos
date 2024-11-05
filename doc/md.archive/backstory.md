# Backstory

Holos is a tool intended to lighten the burden of managing Kubernetes resources.  In 2020 we set out to develop a holistic platform composed from open source cloud native components.  We quickly became frustrated with how each of the major components packaged and distributed their software in a different way.  Many projects choose to distribute their software with Helm charts, while others provide plain yaml files and Kustomize bases.  The popular Kube Prometheus Stack project provides Jsonnet to render and update Kubernetes yaml manifests.
