---
description: Compare Holos with other tools in the ecosystem.
slug: /comparison
sidebar_position: 300
---

# Comparison

:::tip
Holos is designed to complement and improve, not replace, existing tools in the
cloud native ecosystem.
:::

## Helm

### Chart Users

Describe how things are different when using an upstream helm chart.

### Chart Authors

Describe how things are different when writing a new helm chart.

## Kustomize

TODO

## ArgoCD

TODO

## Flux

TODO

## Timoni

| Aspect     | Timoni               | Holos                | Comment                                                                                  |
| ---------- | -------------------- | -------------------- | ---------------------------------------------------------------------------------------- |
| Language   | CUE                  | CUE                  | Like Holos, Timoni is also built on CUE.                                                 |
| Artifact   | OCI Image            | Plain YAML Files     | The Holos Authors find plain files easier to work with and reason about than OCI images. |
| Outputs to | OCI Image Repository | Local Git repository | Holos is designed for use with existing GitOps tools.                                    |
| Concept    | Module               | Component            | A Timoni Module is analogous to a Holos Component.                                       |
| Concept    | Bundle               | Platform             | A Timoni Bundle is somewhat similar, but smaller in scope to a Holos Platform.           |

:::important

The Holos Authors are deeply grateful to Stefan and Timoni for the capability of
importing Kubernetes custom resource definitions into CUE.  Without this
functionality, much of the Kubernetes ecosystem would be more difficult to
manage in CUE and therefore in Holos.

:::


## KubeVela

1. Also built on CUE.
2. Intended to create an Application abstraction.
3. Holos prioritizes composition over abstraction.
4. An abstraction of an Application acts as a filter that removes all but the lowest common denominator functionality.  The Holos Authors have found this filtering effect to create excessive friction for software developers.
5. Holos focuses instead on composition to empower developers and platform engineers to leverage the unique features and functionality of their software and platform.

## Pulumi

TODO

## Jsonnet

TODO
