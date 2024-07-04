using module mode; GOMOD=/Users/jeff/Holos/holos/go.mod

...

<!-- #lowframe -->

[Go Documentation Server](/pkg/)

[GoDoc](/pkg/)

[▽](#)

Search

<!-- magnifying glass: -->

# Package v1alpha2

<!--
    	Copyright 2009 The Go Authors. All rights reserved.
    	Use of this source code is governed by a BSD-style
    	license that can be found in the LICENSE file.
    -->

<!--
    	Note: Static (i.e., not template-generated) href and id
    	attributes start with "pkg-" to make it impossible for
    	them to conflict with generated attributes (some of which
    	correspond to Go identifiers).
    -->

* `import "github.com/holos-run/holos/api/core/v1alpha2"`

- - [Overview](#pkg-overview)
  - [Index](#pkg-index)

<!-- The package's Name is printed as title by the top-level template -->

## Overview ▹

## Overview ▾

Package v1alpha2 contains the core API contract between the holos cli and CUE configuration code. Platform designers, operators, and software developers use this API to write configuration in CUE which \`holos\` loads. The overall shape of the API defines imperative actions \`holos\` should carry out to render the complete yaml that represents a Platform.

[Platform](#Platform) defines the complete configuration of a platform. With the holos reference platform this takes the shape of one management cluster and at least two workload cluster. Each cluster has multiple [HolosComponent](#HolosComponent) resources applied to it.

Each holos component path, e.g. \`components/namespaces\` produces exactly one [BuildPlan](#BuildPlan) which in turn contains a set of [HolosComponent](#HolosComponent) kinds.

The primary kinds of [HolosComponent](#HolosComponent) are:

1. [HelmChart](#HelmChart) to render config from a helm chart.
2. [KustomizeBuild](#KustomizeBuild) to render config from [Kustomize](#Kustomize)
3. [KubernetesObjects](#KubernetesObjects) to render [APIObjects](#APIObjects) defined directly in CUE configuration.

Note that Holos operates as a data pipeline, so the output of a [HelmChart](#HelmChart) may be provided to [Kustomize](#Kustomize) for post-processing.

## Index ▹

## Index ▾

<!-- Table of contents for API; must be named manual-nav to turn off auto nav. -->

* * [Constants](#pkg-constants)
  * [type APIObject](#APIObject)
  * [type APIObjectMap](#APIObjectMap)
  * [type APIObjects](#APIObjects)
  * [type BuildPlan](#BuildPlan)
  * [type BuildPlanComponents](#BuildPlanComponents)
  * [type BuildPlanSpec](#BuildPlanSpec)
  * [type Chart](#Chart)
  * [type FileContent](#FileContent)
  * [type FileContentMap](#FileContentMap)
  * [type FilePath](#FilePath)
  * [type HelmChart](#HelmChart)
  * [type HolosComponent](#HolosComponent)
  * [type Kind](#Kind)
  * [type KubernetesObjects](#KubernetesObjects)
  * [type Kustomize](#Kustomize)
  * [type KustomizeBuild](#KustomizeBuild)
  * [type Label](#Label)
  * [type Metadata](#Metadata)
  * [type Platform](#Platform)
  * [type PlatformMetadata](#PlatformMetadata)
  * [type PlatformSpec](#PlatformSpec)
  * [type PlatformSpecComponent](#PlatformSpecComponent)
  * [type Repository](#Repository)

<!-- #manual-nav -->

### Package files

[apiobjects.go](/src/github.com/holos-run/holos/api/core/v1alpha2/apiobjects.go) [buildplan.go](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go) [constants.go](/src/github.com/holos-run/holos/api/core/v1alpha2/constants.go) [core.go](/src/github.com/holos-run/holos/api/core/v1alpha2/core.go) [doc.go](/src/github.com/holos-run/holos/api/core/v1alpha2/doc.go) [helm.go](/src/github.com/holos-run/holos/api/core/v1alpha2/helm.go) [kubernetesobjects.go](/src/github.com/holos-run/holos/api/core/v1alpha2/kubernetesobjects.go) [kustomizebuild.go](/src/github.com/holos-run/holos/api/core/v1alpha2/kustomizebuild.go)

<!-- .expanded -->

<!-- #pkg-index -->

## Constants

```
const (
    APIVersion    = "v1alpha2"
    BuildPlanKind = "BuildPlan"
    HelmChartKind = "HelmChart"
    // ChartDir is the directory name created in the holos component directory to cache a chart.
    ChartDir = "vendor"
    // ResourcesFile is the file name used to store component output when post-processing with kustomize.
    ResourcesFile = "resources.yaml"
)
```

```
const KubernetesObjectsKind = "KubernetesObjects"
```

## type [APIObject](/src/github.com/holos-run/holos/api/core/v1alpha2/apiobjects.go?s=978:1008#L11) [¶](#APIObject)

APIObject represents the most basic generic form of a single kubernetes api object. Represented as a JSON object internally for compatibility between tools, for example loading from CUE.

```
type APIObject structpb.Struct
```

## type [APIObjectMap](/src/github.com/holos-run/holos/api/core/v1alpha2/apiobjects.go?s=1281:1324#L17) [¶](#APIObjectMap)

APIObjectMap represents the marshalled yaml representation of kubernetes api objects. Do not produce an APIObjectMap directly, instead use [APIObjects](#APIObjects) to produce the marshalled yaml representation from CUE data, then provide the result to [HolosComponent](#HolosComponent).

```
type APIObjectMap map[Kind]map[Label]string
```

## type [APIObjects](/src/github.com/holos-run/holos/api/core/v1alpha2/apiobjects.go?s=1901:2055#L31) [¶](#APIObjects)

APIObjects represents Kubernetes API objects defined directly from CUE code. Useful to mix in resources to any kind of [HolosComponent](#HolosComponent), for example adding an ExternalSecret resource to a [HelmChart](#HelmChart).

[Kind](#Kind) must be the resource kind, e.g. Deployment or Service.

[Label](#Label) is an arbitrary internal identifier to uniquely identify the resource within the context of a \`holos\` command. Holos will never write the intermediate label to rendered output.

Refer to [HolosComponent](#HolosComponent) which accepts an [APIObjectMap](#APIObjectMap) field provided by [APIObjects](#APIObjects).

```
type APIObjects struct {
    APIObjects   map[Kind]map[Label]APIObject `json:"apiObjects"`
    APIObjectMap APIObjectMap                 `json:"apiObjectMap"`
}
```

## type [BuildPlan](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=789:989#L11) [¶](#BuildPlan)

BuildPlan represents a build plan for the holos cli to execute. The purpose of a BuildPlan is to define one or more [HolosComponent](#HolosComponent) kinds. For example a [HelmChart](#HelmChart), [KustomizeBuild](#KustomizeBuild), or [KubernetesObjects](#KubernetesObjects).

A BuildPlan usually has an additional empty [KubernetesObjects](#KubernetesObjects) for the purpose of using the [HolosComponent](#HolosComponent) DeployFiles field to deploy an ArgoCD or Flux gitops resource for the holos component.

```
type BuildPlan struct {
    Kind       string        `json:"kind" cue:"\"BuildPlan\""`
    APIVersion string        `json:"apiVersion" cue:"string | *\"v1alpha2\""`
    Spec       BuildPlanSpec `json:"spec"`
}
```

## type [BuildPlanComponents](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=1335:1715#L25) [¶](#BuildPlanComponents)

```
type BuildPlanComponents struct {
    Resources             map[Label]KubernetesObjects `json:"resources,omitempty"`
    KubernetesObjectsList []KubernetesObjects         `json:"kubernetesObjectsList,omitempty"`
    HelmChartList         []HelmChart                 `json:"helmChartList,omitempty"`
    KustomizeBuildList    []KustomizeBuild            `json:"kustomizeBuildList,omitempty"`
}
```

## type [BuildPlanSpec](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=1056:1333#L18) [¶](#BuildPlanSpec)

BuildPlanSpec represents the specification of the build plan.

```
type BuildPlanSpec struct {
    // Disabled causes the holos cli to take no action over the [BuildPlan].
    Disabled bool `json:"disabled,omitempty"`
    // Components represents multiple [HolosComponent] kinds to manage.
    Components BuildPlanComponents `json:"components,omitempty"`
}
```

## type [Chart](/src/github.com/holos-run/holos/api/core/v1alpha2/helm.go?s=922:1304#L13) [¶](#Chart)

Chart represents a helm chart.

```
type Chart struct {
    // Name represents the chart name.
    Name string `json:"name"`
    // Version represents the chart version.
    Version string `json:"version"`
    // Release represents the chart release when executing helm template.
    Release string `json:"release"`
    // Repository represents the repository to fetch the chart from.
    Repository Repository `json:"repository,omitempty"`
}
```

## type [FileContent](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=117:140#L1) [¶](#FileContent)

FileContent represents file contents.

```
type FileContent string
```

## type [FileContentMap](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=314:358#L2) [¶](#FileContentMap)

FileContentMap represents a mapping of file paths to file contents. Paths are relative to the \`holos\` output "deploy" directory, and may contain sub-directories.

```
type FileContentMap map[FilePath]FileContent
```

## type [FilePath](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=54:74#L1) [¶](#FilePath)

FilePath represents a file path.

```
type FilePath string
```

## type [HelmChart](/src/github.com/holos-run/holos/api/core/v1alpha2/helm.go?s=415:886#L1) [¶](#HelmChart)

HelmChart represents a holos component which wraps around an upstream helm chart. Holos orchestrates helm by providing values obtained from CUE, renders the output using \`helm template\`, then post-processes the helm output yaml using the general functionality provided by [HolosComponent](#HolosComponent), for example [Kustomize](#Kustomize) post-rendering and mixing in additional kubernetes api objects.

```
type HelmChart struct {
    HolosComponent `json:",inline"`
    Kind           string `json:"kind" cue:"\"HelmChart\""`

    // Chart represents a helm chart to manage.
    Chart Chart `json:"chart"`
    // ValuesContent represents the values.yaml file holos passes to the `helm
    // template` command.
    ValuesContent string `json:"valuesContent"`
    // EnableHooks enables helm hooks when executing the `helm template` command.
    EnableHooks bool `json:"enableHooks" cue:"bool | *false"`
}
```

## type [HolosComponent](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=1851:3130#L34) [¶](#HolosComponent)

HolosComponent defines the fields common to all holos component kinds. Every holos component kind should embed HolosComponent.

```
type HolosComponent struct {
    // Kind is a string value representing the resource this object represents.
    Kind string `json:"kind"`
    // APIVersion represents the versioned schema of this representation of an object.
    APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha2\""`
    // Metadata represents data about the holos component such as the Name.
    Metadata Metadata `json:"metadata"`

    // APIObjectMap holds the marshalled representation of api objects.  Useful to
    // mix in resources to each HolosComponent type, for example adding an
    // ExternalSecret to a HelmChart HolosComponent.  Refer to [APIObjects].
    APIObjectMap APIObjectMap `json:"apiObjectMap,omitempty"`

    // DeployFiles represents file paths relative to the cluster deploy directory
    // with the value representing the file content.  Intended for defining the
    // ArgoCD Application resource or Flux Kustomization resource from within CUE,
    // but may be used to render any file related to the build plan from CUE.
    DeployFiles FileContentMap `json:"deployFiles,omitempty"`

    // Kustomize represents a kubectl kustomize build post-processing step.
    Kustomize `json:"kustomize,omitempty"`

    // Skip causes holos to take no action regarding this component.
    Skip bool `json:"skip" cue:"bool | *false"`
}
```

## type [Kind](/src/github.com/holos-run/holos/api/core/v1alpha2/apiobjects.go?s=763:779#L6) [¶](#Kind)

Kind is a kubernetes api object kind. Defined as a type for clarity and type checking.

```
type Kind string
```

## type [KubernetesObjects](/src/github.com/holos-run/holos/api/core/v1alpha2/kubernetesobjects.go?s=205:336#L1) [¶](#KubernetesObjects)

KubernetesObjects represents a [HolosComponent](#HolosComponent) composed of Kubernetes API objects provided directly from CUE using [APIObjects](#APIObjects).

```
type KubernetesObjects struct {
    HolosComponent `json:",inline"`
    Kind           string `json:"kind" cue:"\"KubernetesObjects\""`
}
```

## type [Kustomize](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=4065:4360#L81) [¶](#Kustomize)

Kustomize represents resources necessary to execute a kustomize build. Intended for at least two use cases:

1. Process a [KustomizeBuild](#KustomizeBuild) [HolosComponent](#HolosComponent) which represents raw yaml file resources in a holos component directory.
2. Post process a [HelmChart](#HelmChart) [HolosComponent](#HolosComponent) to inject istio, patch jobs, add custom labels, etc...

```
type Kustomize struct {
    // KustomizeFiles holds file contents for kustomize, e.g. patch files.
    KustomizeFiles FileContentMap `json:"kustomizeFiles,omitempty"`
    // ResourcesFile is the file name used for api objects in kustomization.yaml
    ResourcesFile string `json:"resourcesFile,omitempty"`
}
```

## type [KustomizeBuild](/src/github.com/holos-run/holos/api/core/v1alpha2/kustomizebuild.go?s=165:290#L1) [¶](#KustomizeBuild)

KustomizeBuild represents a [HolosComponent](#HolosComponent) that renders plain yaml files in the holos component directory using \`kubectl kustomize build\`.

```
type KustomizeBuild struct {
    HolosComponent `json:",inline"`
    Kind           string `json:"kind" cue:"\"KustomizeBuild\""`
}
```

## type [Label](/src/github.com/holos-run/holos/api/core/v1alpha2/apiobjects.go?s=654:671#L3) [¶](#Label)

Label is an arbitrary unique identifier internal to holos itself. The holos cli is expected to never write a Label value to rendered output files, therefore use a [Label](#Label) then the identifier must be unique and internal. Defined as a type for clarity and type checking.

A Label is useful to convert a CUE struct to a list, for example producing a list of [APIObject](#APIObject) resources from an [APIObjectMap](#APIObjectMap). A CUE struct using Label keys is guaranteed to not lose data when rendering output because a Label is expected to never be written to the final output.

```
type Label string
```

## type [Metadata](/src/github.com/holos-run/holos/api/core/v1alpha2/buildplan.go?s=3204:3702#L61) [¶](#Metadata)

Metadata represents data about the holos component such as the Name.

```
type Metadata struct {
    // Name represents the name of the holos component.
    Name string `json:"name"`
    // Namespace is the primary namespace of the holos component.  A holos
    // component may manage resources in multiple namespaces, in this case
    // consider setting the component namespace to default.
    //
    // This field is optional because not all resources require a namespace,
    // particularly CRD's and DeployFiles functionality.
    // +optional
    Namespace string `json:"namespace,omitempty"`
}
```

## type [Platform](/src/github.com/holos-run/holos/api/core/v1alpha2/core.go?s=546:1027#L5) [¶](#Platform)

Platform represents a platform to manage. A Platform resource informs holos which components to build. The platform resource also acts as a container for the platform model form values provided by the PlatformService. The primary use case is to collect the cluster names, cluster types, platform model, and holos components to build into one resource.

```
type Platform struct {
    // Kind is a string value representing the resource this object represents.
    Kind string `json:"kind" cue:"\"Platform\""`
    // APIVersion represents the versioned schema of this representation of an object.
    APIVersion string `json:"apiVersion" cue:"string | *\"v1alpha2\""`
    // Metadata represents data about the object such as the Name.
    Metadata PlatformMetadata `json:"metadata"`

    // Spec represents the specification.
    Spec PlatformSpec `json:"spec"`
}
```

## type [PlatformMetadata](/src/github.com/holos-run/holos/api/core/v1alpha2/core.go?s=76:174#L1) [¶](#PlatformMetadata)

```
type PlatformMetadata struct {
    // Name represents the Platform name.
    Name string `json:"name"`
}
```

## type [PlatformSpec](/src/github.com/holos-run/holos/api/core/v1alpha2/core.go?s=1254:1581#L20) [¶](#PlatformSpec)

PlatformSpec represents the specification of a Platform. Think of a platform specification as a list of platform components to apply to a list of kubernetes clusters combined with the user-specified Platform Model.

```
type PlatformSpec struct {
    // Model represents the platform model holos gets from from the
    // PlatformService.GetPlatform rpc method and provides to CUE using a tag.
    Model structpb.Struct `json:"model"`
    // Components represents a list of holos components to manage.
    Components []PlatformSpecComponent `json:"components"`
}
```

## type [PlatformSpecComponent](/src/github.com/holos-run/holos/api/core/v1alpha2/core.go?s=1657:1896#L29) [¶](#PlatformSpecComponent)

PlatformSpecComponent represents a holos component to build or render.

```
type PlatformSpecComponent struct {
    // Path is the path of the component relative to the platform root.
    Path string `json:"path"`
    // Cluster is the cluster name to provide when rendering the component.
    Cluster string `json:"cluster"`
}
```

## type [Repository](/src/github.com/holos-run/holos/api/core/v1alpha2/helm.go?s=1356:1435#L25) [¶](#Repository)

Repository represents a helm chart repository.

```
type Repository struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}
```

Build version go1.22.4.\
Except as [noted](https://developers.google.com/site-policies#restrictions), the content of this page is licensed under the Creative Commons Attribution 3.0 License, and code is licensed under a [BSD license](/LICENSE).\
[Terms of Service](https://golang.org/doc/tos.html) | [Privacy Policy](https://www.google.com/intl/en/policies/privacy/)

<!-- .container -->

<!-- #page -->
