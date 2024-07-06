# Overview

<!-- https://kubernetes.io/docs/contribute/style/diagram-guide/ -->

This tutorial covers the following process of getting started with Holos.

```mermaid
graph LR
    A[1. Install <br>holos] -->
    B[2. Register <br>account] -->
    C[3. Generate <br>platform] -->
    D[4. Render <br>platform] -->
    E[5. Apply <br>config]
 
    classDef box fill:#fff,stroke:#000,stroke-width:1px,color:#000;
    class A,B,C,D,E box
```

## Rendering Pipeline

Holos uses the kubernetes resource model to manage configuration.  The `holos` command line interface (cli) is the primary method you'll use to manage your platform.  Holos uses CUE to provide a unified configuration model of the platform which is built from components packaged with Helm, Kustomize, CUE, or any tool that can produce kubernetes resources as output.  This process can be thought of as a yaml **rendering pipeline**.

Each component in a platform defines a rendering pipeline shown in Figure 2 to produce kubernetes api resources

```mermaid
---
title: Figure 2 - Render Pipeline
---
graph LR
    PS[<a href="/docs/api/core/v1alpha2#PlatformSpec">PlatformSpec</a>]
    BP[<a href="/docs/api/core/v1alpha2#BuildPlan">BuildPlan</a>]
    HC[<a href="/docs/api/core/v1alpha2#HolosComponent">HolosComponent</a>]

    H[<a href="/docs/api/core/v1alpha2#HelmChart">HelmChart</a>]
    K[<a href="/docs/api/core/v1alpha2#KustomizeBuild">KustomizeBuild</a>]
    O[<a href="/docs/api/core/v1alpha2#KubernetesObjects">KubernetesObjects</a>]

    P[<a href="/docs/api/core/v1alpha2#Kustomize">Kustomize</a>]
    Y[Kubernetes <br>Resources]
    G[GitOps <br>Resource]

    C[Kube API Server]

    PS --> BP --> HC
    HC --> H --> P
    HC --> K --> P
    HC --> O --> P

    P --> Y --> C
    P --> G --> C
 
    classDef box fill:#fff,stroke:#000,stroke-width:1px,color:#000;
    class PS,BP,HC,H,K,O,P,Y,G,C box
```

The `holos` cli can be thought of as executing a data pipeline.  The Platform Model is the top level input to the pipeline and specifies the ways your platform varies from other organizations.  The `holos` cli takes the Platform Model as input and executes a series of steps to produce the platform configuration.  The platform configuration output of `holos` are full Kubernetes API resources, suitable for application to a cluster with `kubectl apply -f`, or GitOps tools such as ArgoCD or Flux.
