---
description: Simplified abstraction to generate core v1alpha4 build plans.
sidebar_position: 996
---
<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# v1alpha4

```go
import "github.com/holos-run/holos/api/author/v1alpha4"
```

### Author API

Package v1alpha4 contains ergonomic CUE definitions for Holos component authors. These definitions serve as adapters to produce [Core API](<https://holos.run/docs/api/core/v1alpha4/>) resources for the holos command line tool.

## Index

- [type ArgoConfig](<#ArgoConfig>)
- [type Cluster](<#Cluster>)
- [type ComponentConfig](<#ComponentConfig>)
- [type Fleet](<#Fleet>)
- [type Helm](<#Helm>)
- [type Hostname](<#Hostname>)
- [type Kubernetes](<#Kubernetes>)
- [type Kustomize](<#Kustomize>)
- [type KustomizeConfig](<#KustomizeConfig>)
- [type NameLabel](<#NameLabel>)
- [type Namespace](<#Namespace>)
- [type Organization](<#Organization>)
- [type OrganizationStrict](<#OrganizationStrict>)
- [type Owner](<#Owner>)
- [type Platform](<#Platform>)
- [type Project](<#Project>)
- [type Projects](<#Projects>)
- [type StandardFleets](<#StandardFleets>)


<a name="ArgoConfig"></a>
## type ArgoConfig {#ArgoConfig}

ArgoConfig represents the ArgoCD GitOps configuration associated with a [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#buildplan>). Useful to define once at the root of the Platform configuration and reuse across all components.

```go
type ArgoConfig struct {
    // Enabled causes holos to render an Application resource when true.
    Enabled bool `cue:"true | *false"`
    // RepoURL represents the value passed to the Application.spec.source.repoURL
    // field.
    RepoURL string
    // Root represents the path from the git repository root to the WriteTo output
    // directory, the behavior of the holos render component --write-to flag and
    // the Core API Component WriteTo field.  Used as a prefix for the
    // Application.spec.source.path field.
    Root string `cue:"string | *\"deploy\""`
    // TargetRevision represents the value passed to the
    // Application.spec.source.targetRevision field.  Defaults to the branch named
    // main.
    TargetRevision string `cue:"string | *\"main\""`
    // AppProject represents the ArgoCD Project to associate the Application with.
    AppProject string `cue:"string | *\"default\""`
}
```

<a name="Cluster"></a>
## type Cluster {#Cluster}

Cluster represents a cluster managed by the Platform.

```go
type Cluster struct {
    // Name represents the cluster name, for example "east1", "west1", or
    // "management".
    Name string `json:"name"`
    // Primary represents if the cluster is marked as the primary among a set of
    // candidate clusters.  Useful for promotion of database leaders.
    Primary bool `json:"primary" cue:"true | *false"`
}
```

<a name="ComponentConfig"></a>
## type ComponentConfig {#ComponentConfig}

ComponentConfig represents the configuration common to all kinds of component.

- [Helm](<#Helm>) charts.
- [Kubernetes](<#Kubernetes>) resources generated from CUE.
- [Kustomize](<#Kustomize>) bases.

See the following resources for additional details:

- [Resources](<https://holos.run/docs/api/core/v1alpha4/#Resources>)
- [ArgoConfig](<#ArgoConfig>)
- [KustomizeConfig](<#KustomizeConfig>)
- [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#BuildPlan>)

```go
type ComponentConfig struct {
    // Name represents the BuildPlan metadata.name field.  Used to construct the
    // fully rendered manifest file path.
    Name string
    // Component represents the path to the component producing the BuildPlan.
    Component string
    // Cluster represents the name of the cluster this BuildPlan is for.
    Cluster string
    // Resources represents kubernetes resources mixed into the rendered manifest.
    Resources core.Resources
    // ArgoConfig represents the ArgoCD GitOps configuration for this BuildPlan.
    ArgoConfig ArgoConfig
    // CommonLabels represents common labels to manage on all rendered manifests.
    CommonLabels map[string]string
    // Namespace manages the metadata.namespace field on all resources except the
    // ArgoCD Application.
    Namespace string `json:",omitempty"`

    // KustomizeConfig represents the configuration for kustomize.
    KustomizeConfig KustomizeConfig
}
```

<a name="Fleet"></a>
## type Fleet {#Fleet}

Fleet represents a named collection of similarly configured Clusters. Useful to segregate workload clusters from their management cluster.

```go
type Fleet struct {
    Name string `json:"name"`
    // Clusters represents a mapping of Clusters by their name.
    Clusters map[string]Cluster `json:"clusters" cue:"{[Name=_]: name: Name}"`
}
```

<a name="Helm"></a>
## type Helm {#Helm}

Helm provides a [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#BuildPlan>) via the Output field which generates manifests from a helm chart with optional mix\-in resources provided directly from CUE in the Resources field.

This definition is a convenient way to produce a [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#BuildPlan>) composed of three \[Resources\] generators with one [Kustomize](<#Kustomize>) transformer.

See related:

- [ComponentConfig](<#ComponentConfig>)
- [Chart](<https://holos.run/docs/api/core/v1alpha4/#Chart>)
- [Values](<https://holos.run/docs/api/core/v1alpha4/#Values>)
- [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#BuildPlan>)

```go
type Helm struct {
    ComponentConfig `json:",inline"`

    // Chart represents a Helm chart.
    Chart core.Chart
    // Values represents data to marshal into a values.yaml for helm.
    Values core.Values
    // EnableHooks enables helm hooks when executing the `helm template` command.
    EnableHooks bool

    // BuildPlan represents the derived BuildPlan produced for the holos render
    // component command.
    BuildPlan core.BuildPlan
}
```

<a name="Hostname"></a>
## type Hostname {#Hostname}

Hostname represents the left most dns label of a domain name.

```go
type Hostname struct {
    // Name represents the subdomain to expose, e.g. "www"
    Name string
    // Namespace represents the namespace metadata.name field of backend object
    // reference.
    Namespace string
    // Service represents the Service metadata.name field of backend object
    // reference.
    Service string
    // Port represents the Service port of the backend object reference.
    Port int
}
```

<a name="Kubernetes"></a>
## type Kubernetes {#Kubernetes}

Kubernetes provides a [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#BuildPlan>) via the Output field which contains inline API Objects provided directly from CUE in the Resources field of [ComponentConfig](<#ComponentConfig>).

See related:

- [ComponentConfig](<#ComponentConfig>)
- [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#BuildPlan>)

```go
type Kubernetes struct {
    ComponentConfig `json:",inline"`

    // BuildPlan represents the derived BuildPlan produced for the holos render
    // component command.
    BuildPlan core.BuildPlan
}
```

<a name="Kustomize"></a>
## type Kustomize {#Kustomize}

Kustomize provides a [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#buildplan>) via the Output field which generates manifests from a kustomize kustomization with optional mix\-in resources provided directly from CUE in the Resources field.

See related:

- [ComponentConfig](<#ComponentConfig>)
- [BuildPlan](<https://holos.run/docs/api/core/v1alpha4/#buildplan>)

```go
type Kustomize struct {
    ComponentConfig `json:",inline"`

    // BuildPlan represents the derived BuildPlan produced for the holos render
    // component command.
    BuildPlan core.BuildPlan
}
```

<a name="KustomizeConfig"></a>
## type KustomizeConfig {#KustomizeConfig}

KustomizeConfig represents the configuration for kustomize post processing. The Files field is used to mixing in static manifest files from the component directory. The Resources field is used for mixing in manifests from network locations urls.

See related:

- [ComponentConfig](<#ComponentConfig>)
- [Kustomization](<https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/>)

```go
type KustomizeConfig struct {
    // Kustomization represents the kustomization used to transform resources.
    // Note the resources field is internally managed from the Files and Resources fields.
    Kustomization map[string]any `json:",omitempty"`
    // Files represents files to copy from the component directory for kustomization.
    Files map[string]struct{ Source string } `cue:"{[NAME=_]: Source: NAME}"`
    // Resources represents additional entries to included in the resources list.
    Resources map[string]struct{ Source string } `cue:"{[NAME=_]: Source: NAME}"`
}
```

<a name="NameLabel"></a>
## type NameLabel {#NameLabel}

NameLabel signals the common use case of converting a struct to a list where the name field of each value unifies with the field name of the outer struct.

For example:

```
S: [NameLabel=string]: name: NameLabel
S: jeff: _
S: gary: _
S: nate: _
L: [for x in S {x}]
// L is [{name: "jeff"}, {name: "gary"}, {name: "nate"}]
```

```go
type NameLabel string
```

<a name="Namespace"></a>
## type Namespace {#Namespace}

Namespace represents a Kubernetes namespace.

```go
type Namespace struct {
    Name string
}
```

<a name="Organization"></a>
## type Organization {#Organization}

Organization represents organizational metadata useful across the platform.

```go
type Organization struct {
    Name        string
    DisplayName string
    Domain      string
}
```

<a name="OrganizationStrict"></a>
## type OrganizationStrict {#OrganizationStrict}

OrganizationStrict represents organizational metadata useful across the platform. This is an example of using CUE regular expressions to constrain and validate configuration.

```go
type OrganizationStrict struct {
    Organization `json:",inline"`
    // Name represents the organization name as a resource name.  Must be 63
    // characters or less.  Must start with a letter.  May contain non-repeating
    // hyphens, letters, and numbers.  Must end with a letter or number.
    Name string `cue:"=~ \"^[a-z][0-9a-z-]{1,61}[0-9a-z]$\" & !~ \"--\""`
    // DisplayName represents the human readable organization name.
    DisplayName string `cue:"=~ \"^[0-9A-Za-z][0-9A-Za-z ]{2,61}[0-9A-Za-z]$\" & !~ \"  \""`
}
```

<a name="Owner"></a>
## type Owner {#Owner}

Owner represents the owner of a resource. For example, the name and email address of an engineering team.

```go
type Owner struct {
    Name  string
    Email string
}
```

<a name="Platform"></a>
## type Platform {#Platform}

Platform assembles a Core API [Platform](<https://holos.run/docs/api/core/v1alpha4/#Platform>) in the Resource field for the holos render platform command. Use the Components field to register components with the platform using a struct. This struct is converted into a list for final output to holos.

See related:

- [Component](<https://holos.run/docs/api/core/v1alpha4/#Component>) collection of components composing the platform.
- [Platform](<https://holos.run/docs/api/core/v1alpha4/#Platform>) resource assembled for holos to process.

```go
type Platform struct {
    Name       string
    Components map[NameLabel]core.Component
    Resource   core.Platform
}
```

<a name="Project"></a>
## type Project {#Project}

Project represents logical grouping of components owned by one or more teams. Useful for the platform team to manage resources for project teams to use.

```go
type Project struct {
    // Name represents project name.
    Name string
    // Owner represents the team who own this project.
    Owner Owner
    // Namespaces represents the namespaces assigned to this project.
    Namespaces map[NameLabel]Namespace
    // Hostnames represents the host names to expose for this project.
    Hostnames map[NameLabel]Hostname
    // CommonLabels represents common labels to manage on all rendered manifests.
    CommonLabels map[string]string
}
```

<a name="Projects"></a>
## type Projects {#Projects}

Projects represents projects managed by the platform team for use by other teams using the platform.

```go
type Projects map[NameLabel]Project
```

<a name="StandardFleets"></a>
## type StandardFleets {#StandardFleets}

StandardFleets represents the standard set of Clusters in a Platform segmented into Fleets by their purpose. The management Fleet contains a single Cluster, for example a GKE autopilot cluster with no workloads deployed for reliability and cost efficiency. The workload Fleet contains all other Clusters which contain workloads and sync Secrets from the management cluster.

```go
type StandardFleets struct {
    // Workload represents a Fleet of zero or more workload Clusters.
    Workload Fleet `json:"workload" cue:"{name: \"workload\"}"`
    // Management represents a Fleet with one Cluster named management.
    Management Fleet `json:"management" cue:"{name: \"management\"}"`
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
