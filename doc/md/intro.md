# Introduction

⚡️ Holos will help you build your **software development platform in no time.**

💸 Building a software development platform is time consuming and expensive. Instead, focus on your application and let Holos manage the platform for you.

💥 Already delivering? Use the advanced functionality of Holos to manage your existing platform safer and easier.

🧐 Holos is a platform builder. It builds a hollistic software development platform composed of best-of-breed cloud native open source projects.  Holos is also a tool to make it easier to manage cloud infrastructure by providing a typed alternative to yaml templates.

## Features

Holos was built to solve two main problems:

 1. Building a platform is **time consuming and expensive**, usually taking 3 engineers 6-9 months of effort.  Holos provides a reference platform that enables an individual to build and customize their platform quickly.
 2. Configuration changes often lead to downtime.  Existing tools like Helm make it difficult to reason about the impact a configuration change will have.  Holos provides a unique, unified configuration model powered by CUE which makes it safer and easier to manage the platform.

In addition, a core principle of Holos is that Organizations gain value from owning the the platform they build on.  Avoid vendor lock-in, future price hikes, and licensing changes by building on a solid foundation of open source, cloud native computing foundation backed projects.

The following features are built into the Holos reference platform.

:::tip

Prefer a different implementation for a feature, or want a feature not listed?  No problem, Holos is highly customizable and composable.

:::

- **Continuous Delivery**
  - Holos builds a GitOps workflow for each application running in the platform.
  - Developers push changes which are automatically deployed.
  - Powered by [ArgoCD](https://argo-cd.readthedocs.io/en/stable/)
- **Identity and Access Management** (IAM)
  - Holos builds a standard OIDC identity provider for you.
  - Integrates with your exisitng IAM and SSO system, or works independently.
  - Powerful customer identity and access management features.
  - Role based access control.
  - Powered by [ZITADEL](https://zitadel.com/)
- **Zero Trust**
  - Authenticate and Authorize users at the platform layer instead of or in addition to the application layer.
  - Integrated with observability to measure and alert about problems before customers complain.
  - Powered by [Istio](https://istio.io/)
- **Observability**
  - Holos collects performance and availability metrics automatically, without requiring application changes.
  - Optional, deeper integration into the application layer.
  - Distributed Tracing
  - Logging
  - Powered by Prometheus, Grafana, Loki, and OpenTelemetry.
- **Data Platform**
  - Integrated management of PostgreSQL
  - Automatic backups
  - Automatic restore from backup
  - Quickly fail over across multiple regions
- **Multi-Region**
  - Holos is designed to operate in multiple regions and multiple clouds.
  - Keep customer data in the region that makes the most sense for your business.
  - Easily cut over from one region to another for redundancy and business continuity.

## Development Status

Holos is being actively developed by [Open Infrastructure Services](https://openinfrastructure.co).  Release can be found [here](https://github.com/holos-run/holos/releases).

## Adoption

Organizations who have officially adopted Holos can be found [here](https://github.com/holos-run/holos/blob/main/ADOPTERS.md).
