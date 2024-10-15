// # Author API
//
// Package v1alpha4 contains ergonomic CUE definitions for Holos component
// authors.  These definitions serve as adapters to produce [Core API] resources
// for the holos command line tool.
//
// [Core API]: https://holos.run/docs/api/core/v1alpha4/
package v1alpha4

import core "github.com/holos-run/holos/api/core/v1alpha4"

//go:generate ../../../hack/gendoc

// Platform assembles a Core API [Platform] in the Resource field for the holos
// render platform command.  Use the Components field to register components
// with the platform using a struct.  This struct is converted into a list for
// final output to holos.
//
// [Platform]: https://holos.run/docs/api/core/v1alpha4/#Platform
type Platform struct {
	Name       string
	Components map[core.NameLabel]core.Component
	Resource   core.Platform
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

// ArgoConfig represents the ArgoCD GitOps configuration associated with a
// [BuildPlan].  Useful to define once at the root of the Platform configuration
// and reuse across all components.
//
// [BuildPlan]: https://holos.run/docs/api/core/v1alpha4/#buildplan
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

// Organization represents organizational metadata useful across the platform.
type Organization struct {
	Name        string
	DisplayName string
	Domain      string
}

// OrganizationStrict represents organizational metadata useful across the
// platform.  This is an example of using CUE regular expressions to constrain
// and validate configuration.
type OrganizationStrict struct {
	Organization `json:",inline"`
	// Name represents the organization name as a resource name.  Must be 63
	// characters or less.  Must start with a letter.  May contain non-repeating
	// hyphens, letters, and numbers.  Must end with a letter or number.
	Name string `cue:"=~ \"^[a-z][0-9a-z-]{1,61}[0-9a-z]$\" & !~ \"--\""`
	// DisplayName represents the human readable organization name.
	DisplayName string `cue:"=~ \"^[0-9A-Za-z][0-9A-Za-z ]{2,61}[0-9A-Za-z]$\" & !~ \"  \""`
}

// Kubernetes provides a BuildPlan via the Output field which contains inline
// API Objects provided directly from CUE in the Resources field.
//
// This definition is a convenient way to produce a [BuildPlan] composed of a
// two [Resources] generators with one [Kustomize] transformer.  The two
// generators produce Kubernetes API resources managed by an ArgoCD Application,
// transformed with Kustomize to consistently manage the metadata.namespace
// field and common labels.
//
// [BuildPlan]: https://holos.run/docs/api/core/v1alpha4/#BuildPlan
// [Resources]: https://holos.run/docs/api/core/v1alpha4/#Resources
// [Kustomize]: https://holos.run/docs/api/core/v1alpha4/#Kustomize
type Kubernetes struct {
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

	// BuildPlan represents the derived BuildPlan produced for the holos render
	// component command.
	BuildPlan core.BuildPlan
}
