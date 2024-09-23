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
- The [API Reference] is a technical reference for when you're writing CUE code to define your your platform.

## Backstory

At [Open Infrastructure Services], each of us has helped dozens of companies
build and operate their software development platform over the course of our
career.  During the U.S. presidential election just before the pandemic our
second largest client, Twitter, had a global outage lasting the better part of a
day.  At the time, we were helping the core infrastructure team by managing
their production configuration management system so the team could focus on
delivering key objectives for the business.  In this position, we had a front
row seat into what happened that day.

One of Twitter's employees, a close friend of ours and engineer on the team,
landed a trivial one line change to the firewall configuration.  Less than 30
minutes later literally everything was down.  The one line change, which passed
code review and seemed harmless, resulted in the host firewall reverting to the
default state on hundreds of thousands of servers.  All connections to all
servers globally were blocked and dropped.  Except SSH, thankfully.  At least
one Presidential candidate complained loudly.

This incident led us to deeply consider a few key problems about Twitter's
software development platform, problems which made their way all the way up to
the board of directors to solve.

1. **Lack of Visibility** - There was no way to see clearly the impact a simple one-line
change could have on the platform.  This lack of visibility made it difficult
for engineers to reason about changes they were writing and reviewing.
2. **Large Blast Radius** - All changes, no matter how small or large, affected the
entire global fleet within 30 minutes.  Twitter needed a way to cap the
potential blast radius to prevent global outages.
3. **Incomplete Tooling** - Twitter had the correct processes in place, but
their tooling didn't support their process.  The one line change was tested and
peer reviewed prior to landing, but the tooling didn't surface the information
they needed when they needed it.

Over the next few years we built features for Twitter's configuration management
system that solved each of these problems.  At the same time, I started
exploring my curiosity of what these solutions would look like in the context of
Kubernetes and cloud native software instead of a traditional configuration
management context.

As Google Cloud partners, we had the opportunity to work with Google's largest
customers and learn how they built their software development platforms on
Kubernetes.  Over the course of the pandemic, we built a software development
platform made largely in the same way, taking off the shelf CNCF projects like
ArgoCD, the Kubernetes Prometheus Stack, Istio, Cert Manager, External Secrets
Operator, etc... and integrating them into a holistic software development
platform.  We started with the packaging recommended by the upstream project.
Helm was and still is the most common distribution method, but many projects
also provided plain yaml, kustomize bases, or sometimes even jsonnet in the case
of the prometheus community.  We then wrote scripts to mix-in the resources we
needed to integrate each piece of software with the platform as a whole. For
example, we often passed Helm's output to Kustomize to add common labels or fix
bugs in the upstream chart, like missing namespace fields.  We wrote umbrella
charts to mix in Ingress, HTTPRoute, and ExternalSecret resources to the vendor
provided chart.

We viewed these scripts as necessary glue to assemble and fasten the components
together into a platform, but we were never fully satisfied with them.  Umbrella
charts became difficult to maintain once there were multiple environments,
regions, and cloud providers in the platform.  Nested for loops in yaml
templates created significant friction and were a challenge to troubleshoot
because they obfuscated what was happening. The scripts, too, made it difficult
to see what was happening, when, and fix issues in them that affected all
components in the platform.

Despite the makeshift scripts and umbrella charts, the overall approach had a
significant advantage. The scripts always produced fully rendered manifests
stored in plain text files.  We committed these files to version control and
used ArgoCD to apply them. The ability to make a one-line change, render the
whole platform, then see clearly what changed platform-wide resulted in less
time spent troubleshooting and fewer errors making their way to production.

For awhile, we lived with the scripts and charts.  I couldn't stop thinking
about the [Why are we templating
YAML?](https://hn.algolia.com/?dateRange=all&page=0&prefix=false&query=https%3A%2F%2Fleebriggs.co.uk%2Fblog%2F2019%2F02%2F07%2Fwhy-are-we-templating-yaml&sort=byDate&type=story)
post on Hacker News though.  I was curious what it would look like to replace
our scripts and umbrella charts with something that helped address the 3 main
problems Twitter experienced.

After doing quite a bit of digging and talking with folks, I found
[CUE](https://cuelang.org).  I took our scripts and charts and re-wrote the same
functionality we needed from the glue-layer in Go and CUE. The result is a new
tool, `holos`, built to complement Helm, Kustomize, and JSonnet, but not replace
them.  Holos leverages CUE to make it easier and safer for teams to define
golden paths and paved roads without having to write bespoke, makeshift scripts
or template text.

Thanks for reading this far, give Holos a try locally with out [Quickstart]
guide.

[Guides]: /docs/guides/
[API Reference]: /docs/api/
[Quickstart]: /docs/quickstart/
[CUE]: https://cuelang.org/
[Author API]: /docs/api/author/
[Core API]: /docs/api/core/
[Open Infrastructure Services]: https://openinfrastructure.co/
