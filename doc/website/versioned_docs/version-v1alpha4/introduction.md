---
description: Holos Documentation
slug: /
---

# Introduction

Welcome to Holos.  Holos is an open source tool to manage software development
platforms safely, easily, and consistently.  We built Holos to help engineering
teams work more efficiently together by empowering them to build golden paths
and paved roads for other teams to leverage for quicker delivery.

## Documentation

- [Guides] are organized into example use-cases of how Holos helps engineering
teams at the fictional Bank of Holos deliver business value on the bank's
platform.
- The [API Reference] is a technical reference for when you're writing CUE code to define your platform.

## Backstory

At [Open Infrastructure Services], we've each helped dozens of companies build and operate their software development platforms. During the U.S. presidential election just before the pandemic, our second-largest client, Twitter, experienced a global outage that lasted nearly a full day. We were managing their production configuration system, allowing the core infrastructure team to focus on business-critical objectives. This gave us a front-row seat to the incident.

A close friend and engineer on the team made a trivial one-line change to the firewall configuration. Less than 30 minutes later, everything was down. That change, which passed code review, caused the host firewall to revert to its default state on hundreds of thousands of servers, blocking all connections globally—except for SSH, thankfully. Even a Presidential candidate complained loudly.

This incident forced us to reconsider key issues with Twitter's platform:

1. **Lack of Visibility** - Engineers couldn't foresee the impact of even a small change, making it difficult to assess risks.
2. **Large Blast Radius** - Small changes affected the entire global fleet in under 30 minutes. There was no way to limit the impact of a single change.
3. **Incomplete Tooling** - The right processes were in place, but the tooling didn't fully support them. The change was tested and reviewed, but critical information wasn't surfaced in time.

Over the next few years, we built features to address these issues. Meanwhile, I began exploring how these solutions could work in the Kubernetes and cloud-native space.

As Google Cloud partners, we worked with large customers to understand how they built their platforms on Kubernetes. During the pandemic, we built a platform using CNCF projects like ArgoCD, Prometheus Stack, Istio, Cert Manager, and External Secrets Operator, integrating them into a cohesive platform. We started with upstream recommendations—primarily Helm charts—and wrote scripts to integrate each piece into the platform. For example, we passed Helm outputs to Kustomize to add labels or fix bugs, and wrote umbrella charts to add Ingress, HTTPRoute, and ExternalSecret resources.

These scripts served as necessary glue to hold everything together but became difficult to manage across multiple environments, regions, and cloud providers. YAML templates and nested loops created friction, making them hard to troubleshoot. The scripts themselves made it difficult to see what was happening and to fix issues affecting the entire platform.

Still, the scripts had a key advantage: they produced fully rendered manifests in plain text, committed to version control, and applied via ArgoCD. This clarity made troubleshooting easier and reduced errors in production.

Despite the makeshift nature of the scripts, I kept thinking about the "[Why are we templating YAML]?" post on Hacker News. I wanted to replace our scripts and charts with something more robust and easier to maintain—something that addressed Twitter's issues head-on.

I rewrote our scripts and charts using CUE and Go, replacing the glue layer. The result is **Holos**—a tool designed to complement Helm, Kustomize, and Jsonnet, making it easier and safer to define golden paths and paved roads without bespoke scripts or templates.

Thanks for reading. Take Holos for a spin on your local machine with our [Quickstart] guide.

[Guides]: /docs/guides/
[API Reference]: /docs/api/
[Quickstart]: /docs/quickstart/
[CUE]: https://cuelang.org/
[Author API]: /docs/api/author/
[Core API]: /docs/api/core/
[Open Infrastructure Services]: https://openinfrastructure.co/
[Why are we templating YAML]: https://hn.algolia.com/?dateRange=all&page=0&prefix=false&query=https%3A%2F%2Fleebriggs.co.uk%2Fblog%2F2019%2F02%2F07%2Fwhy-are-we-templating-yaml&sort=byDate&type=story
