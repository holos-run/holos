// Package v1alpha3 contains the core API contract between the holos cli and CUE
// configuration code.  Platform designers, operators, and software developers
// use this API to write configuration in CUE which `holos` loads.  The overall
// shape of the API defines imperative actions `holos` should carry out to
// render the complete yaml that represents a Platform.
//
// [Platform] defines the complete configuration of a platform.  With the holos
// reference platform this takes the shape of one management cluster and at
// least two workload cluster.  Each cluster has multiple [Component]
// resources applied to it.
//
// Each holos component path, e.g. `components/namespaces` produces exactly one
// [BuildPlan] which in turn contains a set of [Component] kinds.
//
// The primary kinds of [Component] are:
//
//  1. [HelmChart] to render config from a helm chart.
//  2. [KustomizeBuild] to render config from [Kustomize]
//  3. [KubernetesObjects] to render [APIObjects] defined directly in CUE
//     configuration.
//
// Note that Holos operates as a data pipeline, so the output of a [HelmChart]
// may be provided to [Kustomize] for post-processing.
package v1alpha3

//go:generate ../../../hack/gendoc
