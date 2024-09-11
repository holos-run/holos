// Package v1alpha3 contains CUE definitions intended as convenience wrappers
// around the core data types defined in package core.  The purpose of these
// wrappers is to make life easier for platform engineers by reducing boiler
// plate code and generating component build plans in a consistent manner.
package v1alpha3

import (
	core "github.com/holos-run/holos/api/core/v1alpha3"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:generate ../../../hack/gendoc

// Helm provides a BuildPlan via the Output field which contains one HelmChart
// from package core.  Useful as a convenience wrapper to render a HelmChart
// with optional mix-in resources and Kustomization post-processing.
type Helm struct {
	// Name represents the chart name.
	Name string
	// Version represents the chart version.
	Version string
	// Namespace represents the helm namespace option when rendering the chart.
	Namespace string
	// Resources are kubernetes api objects to mix into the output.
	Resources map[string]any `cue:"{...}"`

	// Repo represents the chart repository
	Repo struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// Values represents data to marshal into a values.yaml for helm.
	Values interface{} `cue:"{...}"`

	// Chart represents the derived HelmChart for inclusion in the BuildPlan
	// Output field value.  The default HelmChart field values are derived from
	// other Helm field values and should be sufficient for most use cases.
	Chart core.HelmChart

	// EnableKustomizePostProcessor processes helm output with kustomize if true.
	EnableKustomizePostProcessor bool `cue:"true | *false"`

	// KustomizeFiles represents additional files to include in a Kustomization
	// resources list.  Useful to patch helm output.  The implementation is a
	// struct with filename keys and structs as values.  Holos encodes the struct
	// value to yaml then writes the result to the filename key.  Component
	// authors may then reference the filename in the kustomization.yaml resources
	// or patches lists.
	// Requires EnableKustomizePostProcessor: true.
	KustomizeFiles map[string]any `cue:"{[string]: {...}}"`

	// KustomizePatches represents patches to apply to the helm output.  Requires
	// EnableKustomizePostProcessor: true.
	KustomizePatches map[core.InternalLabel]any `cue:"{[string]: {...}}"`

	// KustomizeResources represents additional resources files to include in the
	// kustomize resources list.
	KustomizeResources map[string]any `cue:"{[string]: {...}}"`

	// ArgoConfig represents the ArgoCD GitOps configuration for this Component.
	ArgoConfig ArgoConfig

	// Output represents the derived BuildPlan for the Holos cli to render.
	Output core.BuildPlan
}

// Resources represents the default schema for a Kubernetes API object resource.
// For example, a Service, Namespace or Deployment.  The top level key is the
// kind of resource so default behavior and strict schema enforcement may be
// enforced for the kind.  The second level keys are an arbitrary internal
// label, which serves as the default value for the resource metadata name
// field, but may differ for situations where the same resource kind and name
// are managed in different namespaces.
//
// Refer to [definitions.cue] for the CUE schema definition as an example to
// build on when defining your own Components.
//
// [definitions.cue]: https://github.com/holos-run/holos/blob/main/internal/generate/platforms/cue.mod/pkg/github.com/holos-run/holos/api/schema/v1alpha3/definitions.cue#L9
// type Resources map[string]map[string]any

// ArgoConfig represents the ArgoCD GitOps configuration for a Component.
// Useful to define once at the root of the Platform configuration and reuse
// across all Components.
type ArgoConfig struct {
	// Enabled causes holos to render an ArgoCD Application resource for GitOps if true.
	Enabled bool `cue:"true | *false"`
	// ClusterName represents the cluster within the platform the Application
	// resource is intended for.
	ClusterName string
	// DeployRoot represents the path from the git repository root to the `deploy`
	// rendering output directory.  Used as a prefix for the
	// Application.spec.source.path field.
	DeployRoot string `cue:"string | *\".\""`
	// RepoURL represents the value passed to the Application.spec.source.repoURL
	// field.
	RepoURL string
	// TargetRevision represents the value passed to the
	// Application.spec.source.targetRevision field.  Defaults to the branch named
	// main.
	TargetRevision string `cue:"string | *\"main\""`
}

// Cluster represents a cluster managed by the Platform.
type Cluster struct {
	// Name represents the cluster name, for example "east1", "west1", or
	// "management".
	Name string `json:"name"`
	// Primary represents if the cluster is marked as the primary among a set of
	// candidate clusters.  Useful for promotion of database leaders.
	Primary bool `json:"primary" cue:"true | *false"`
}

// Fleet represents a named collection of similarly configured Clusters.  Useful
// to segregate workload clusters from their management cluster.
type Fleet struct {
	Name string `json:"name"`
	// Clusters represents a mapping of Clusters by their name.
	Clusters map[string]Cluster `json:"clusters" cue:"{[Name=_]: name: Name}"`
}

// StandardFleets represents the standard set of Clusters in a Platform
// segmented into Fleets by their purpose.  The management Fleet contains a
// single Cluster, for example a GKE autopilot cluster with no workloads
// deployed for reliability and cost efficiency.  The workload Fleet contains
// all other Clusters which contain workloads and sync Secrets from the
// management cluster.
type StandardFleets struct {
	// Workload represents a Fleet of zero or more workload Clusters.
	Workload Fleet `json:"workload" cue:"{name: \"workload\"}"`
	// Management represents a Fleet with one Cluster named management.
	Management Fleet `json:"management" cue:"{name: \"management\"}"`
}

// Platform is a convenience structure to produce a core Platform specification
// value in the Output field.  Useful to collect components at the root of the
// Platform configuration tree as a struct, which are automatically converted
// into a list for the core Platform spec output.
type Platform struct {
	// Name represents the Platform name.
	Name string `cue:"string | *\"holos\""`
	// Components is a structured map of components to manage by their name.
	Components map[string]core.PlatformSpecComponent
	// Model represents the Platform model holos gets from from the
	// PlatformService.GetPlatform rpc method and provides to CUE using a tag.
	Model structpb.Struct `cue:"{...}"`
	// Output represents the core Platform spec for the holos cli to iterate over
	// and render each listed Component, injecting the Model.
	Output core.Platform
}
