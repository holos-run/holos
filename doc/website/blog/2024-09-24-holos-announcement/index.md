---
slug: holos-platform-manager
title: Holos Platform Manager
authors: [jeff]
tags: [holos]
---

## Introducing Holos

I’m excited to announce Holos, a tool designed to help engineering teams
manage their software development platforms built on the Kubernetes resource
model.

:::tip
For a hands-on introduction, check out our [Quickstart] Guide.
:::

### The Backstory

In our roles at [Open Infrastructure Services], and earlier at Puppet, we helped
many companies automate infrastructure management. In 2017, we had the
opportunity to work with Twitter to improve their configuration management
system.  This opportunity gave us insight into the challenges of managing a
large-scale platform with multiple engineering teams. Our work involved
everything from observability systems to application deployment workflows and of
course, managing the core infrastructure.

This experience demonstrated the value of platform engineering. As the pandemic
hit, I began thinking about what a fully cloud-native platform might look like
using the Kubernetes resource model. Around the same time, I came across the
Hacker News post, “[Why Are We Templating YAML]?”, which sparked a good
discussion. It was clear I wasn’t alone in my frustration with managing YAML
files and ensuring clear, predictable changes before merging them into
production.

A common pain point and theme is the complexity of working with nested YAML
configurations, especially with tools like ArgoCD and Helm. The lack of a
standard for rendering YAML templates makes it difficult to see what changes are
actually being applied to the Kubernetes API. This often results in trial and
error, costly blue-green deployments, and hours of debugging.

During the pandemic, I began experimenting with a tool to address this issue,
drawing on lessons from our work at Twitter. The key problems we aimed to solve
are:

- **Lack of visibility**: Engineers struggled to foresee the impact of small changes.
- **Large blast radius**: Small changes affected global systems, with no way to limit the impact.
- **Incomplete tooling**: While processes were in place, the right information wasn’t surfaced at the right time.

We built several iterations of a reference platform based on Kubernetes,
initially focusing on fully rendering manifests into plain files—a pattern now
called the [rendered manifests pattern]. Over time, we realized we were spending
most of our time maintaining bash scripts and YAML templates. This led back to
the question: Why are we templating YAML? What _should_ replace templates?

We'd previously seen a colleague use CUE effectively to generate large scale
configurations for Envoy, and ran into CUE again when we worked on a project
involving Dagger, but I still hadn't taken a deep look at CUE.

At the end of 2023, I decided to dive deep with [CUE].  I quickly came to
appreciate CUE’s unified approach where **order is irrelevant**.  Before CUE, we
handled configuration data in a hierarchy with a precedence ordering, similar to
how we handled data in Puppet with Hiera.  CUE's promise of no longer needing to
think about ordering and precedence rules held, alleviating a large cognitive
burden when dealing with complex configurations. CUE quickly allowed me to
replace the unmaintainable bash scripts and complex Helm templates, simplifying
our workflow.

### Enter Holos

Holos adds CUE as a well-specified integration layer over tools like Helm,
Kustomize, ArgoCD, and Crossplane. With Holos, we can now efficiently integrate
upstream Helm charts and Kustomize bases into our platform without the
complexity of templates and scripts. This has also made it easy for one team to
define "golden paths" that other teams can follow—like automatically configuring
namespaces and security policies when dev teams start new projects.

We've found Holos incredibly useful and hope you do too. Let us know your
thoughts!

[Guides]: /docs/guides/
[API Reference]: /docs/api/
[Quickstart]: /docs/quickstart/
[CUE]: https://cuelang.org/
[Author API]: /docs/api/author/
[Core API]: /docs/api/core/
[Open Infrastructure Services]: https://openinfrastructure.co/
[Why are we templating YAML]: https://hn.algolia.com/?dateRange=all&page=0&prefix=false&query=https%3A%2F%2Fleebriggs.co.uk%2Fblog%2F2019%2F02%2F07%2Fwhy-are-we-templating-yaml&sort=byDate&type=story

[Holos]: https://holos.run/
[Quickstart]: /docs/quickstart/
[rendered manifests pattern]: https://akuity.io/blog/the-rendered-manifests-pattern/
