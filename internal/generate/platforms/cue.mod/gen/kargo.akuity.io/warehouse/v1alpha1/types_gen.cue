// Code generated by timoni. DO NOT EDIT.

//timoni:generate timoni vendor crd -f /Users/jeff/Holos/kargo-demo/deploy/components/kargo/kargo.gen.yaml

package v1alpha1

import "strings"

// Warehouse is a source of Freight.
#Warehouse: {
	// APIVersion defines the versioned schema of this representation
	// of an object.
	// Servers should convert recognized schemas to the latest
	// internal value, and
	// may reject unrecognized values.
	// More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	apiVersion: "kargo.akuity.io/v1alpha1"

	// Kind is a string value representing the REST resource this
	// object represents.
	// Servers may infer this from the endpoint the client submits
	// requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	kind: "Warehouse"
	metadata!: {
		name!: strings.MaxRunes(253) & strings.MinRunes(1) & {
			string
		}
		namespace!: strings.MaxRunes(63) & strings.MinRunes(1) & {
			string
		}
		labels?: {
			[string]: string
		}
		annotations?: {
			[string]: string
		}
	}

	// Spec describes sources of artifacts.
	spec!: #WarehouseSpec
}

// Spec describes sources of artifacts.
#WarehouseSpec: {
	// FreightCreationPolicy describes how Freight is created by this
	// Warehouse.
	// This field is optional. When left unspecified, the field is
	// implicitly
	// treated as if its value were "Automatic".
	// Accepted values: Automatic, Manual
	freightCreationPolicy?: "Automatic" | "Manual" | *"Automatic"

	// Interval is the reconciliation interval for this Warehouse. On
	// each
	// reconciliation, the Warehouse will discover new artifacts and
	// optionally
	// produce new Freight. This field is optional. When left
	// unspecified, the
	// field is implicitly treated as if its value were "5m0s".
	interval: =~"^([0-9]+(\\.[0-9]+)?(s|m|h))+$" | *"5m0s"

	// Shard is the name of the shard that this Warehouse belongs to.
	// This is an
	// optional field. If not specified, the Warehouse will belong to
	// the default
	// shard. A defaulting webhook will sync this field with the value
	// of the
	// kargo.akuity.io/shard label. When the shard label is not
	// present or differs
	// from the value of this field, the defaulting webhook will set
	// the label to
	// the value of this field. If the shard label is present and this
	// field is
	// empty, the defaulting webhook will set the value of this field
	// to the value
	// of the shard label.
	shard?: string

	// Subscriptions describes sources of artifacts to be included in
	// Freight
	// produced by this Warehouse.
	subscriptions: [...{
		// Chart describes a subscription to a Helm chart repository.
		chart?: {
			// DiscoveryLimit is an optional limit on the number of chart
			// versions that
			// can be discovered for this subscription. The limit is applied
			// after
			// filtering charts based on the SemverConstraint field.
			// When left unspecified, the field is implicitly treated as if
			// its value
			// were "20". The upper limit for this field is 100.
			discoveryLimit?: int & <=100 & >=1 | *20

			// Name specifies the name of a Helm chart to subscribe to within
			// a classic
			// chart repository specified by the RepoURL field. This field is
			// required
			// when the RepoURL field points to a classic chart repository and
			// MUST
			// otherwise be empty.
			name?: string

			// RepoURL specifies the URL of a Helm chart repository. It may be
			// a classic
			// chart repository (using HTTP/S) OR a repository within an OCI
			// registry.
			// Classic chart repositories can contain differently named
			// charts. When this
			// field points to such a repository, the Name field MUST also be
			// used
			// to specify the name of the desired chart within that
			// repository. In the
			// case of a repository within an OCI registry, the URL implicitly
			// points to
			// a specific chart and the Name field MUST NOT be used. The
			// RepoURL field is
			// required.
			repoURL: strings.MinRunes(1) & {
				=~"^(((https?)|(oci))://)([\\w\\d\\.\\-]+)(:[\\d]+)?(/.*)*$"
			}

			// SemverConstraint specifies constraints on what new chart
			// versions are
			// permissible. This field is optional. When left unspecified,
			// there will be
			// no constraints, which means the latest version of the chart
			// will always be
			// used. Care should be taken with leaving this field unspecified,
			// as it can
			// lead to the unanticipated rollout of breaking changes.
			// More info:
			// https://github.com/masterminds/semver#checking-version-constraints
			semverConstraint?: string
		}

		// Git describes a subscriptions to a Git repository.
		git?: {
			// AllowTags is a regular expression that can optionally be used
			// to limit the
			// tags that are considered in determining the newest commit of
			// interest. The
			// value in this field only has any effect when the
			// CommitSelectionStrategy is
			// Lexical, NewestTag, or SemVer. This field is optional.
			allowTags?: string

			// Branch references a particular branch of the repository. The
			// value in this
			// field only has any effect when the CommitSelectionStrategy is
			// NewestFromBranch or left unspecified (which is implicitly the
			// same as
			// NewestFromBranch). This field is optional. When left
			// unspecified, (and the
			// CommitSelectionStrategy is NewestFromBranch or unspecified),
			// the
			// subscription is implicitly to the repository's default branch.
			branch?: strings.MinRunes(1) & {
				=~"^\\w+([-/]\\w+)*$"
			}

			// CommitSelectionStrategy specifies the rules for how to identify
			// the newest
			// commit of interest in the repository specified by the RepoURL
			// field. This
			// field is optional. When left unspecified, the field is
			// implicitly treated
			// as if its value were "NewestFromBranch".
			// Accepted values: Lexical, NewestFromBranch, NewestTag, SemVer
			commitSelectionStrategy?: "Lexical" | "NewestFromBranch" | "NewestTag" | "SemVer" | *"NewestFromBranch"

			// DiscoveryLimit is an optional limit on the number of commits
			// that can be
			// discovered for this subscription. The limit is applied after
			// filtering
			// commits based on the AllowTags and IgnoreTags fields.
			// When left unspecified, the field is implicitly treated as if
			// its value
			// were "20". The upper limit for this field is 100.
			discoveryLimit?: int & <=100 & >=1 | *20

			// ExcludePaths is a list of selectors that designate paths in the
			// repository
			// that should NOT trigger the production of new Freight when
			// changes are
			// detected therein. When specified, changes in the identified
			// paths will not
			// trigger Freight production. When not specified, paths that
			// should trigger
			// Freight production will be defined solely by IncludePaths.
			// Selectors may be
			// defined using:
			// 1. Exact paths to files or directories (ex. "charts/foo")
			// 2. Glob patterns (prefix the pattern with "glob:"; ex.
			// "glob:*.yaml")
			// 3. Regular expressions (prefix the pattern with "regex:" or
			// "regexp:";
			// ex. "regexp:^.*\.yaml$")
			// Paths selected by IncludePaths may be unselected by
			// ExcludePaths. This
			// is a useful method for including a broad set of paths and then
			// excluding a
			// subset of them.
			excludePaths?: [...string]

			// IgnoreTags is a list of tags that must be ignored when
			// determining the
			// newest commit of interest. No regular expressions or glob
			// patterns are
			// supported yet. The value in this field only has any effect when
			// the
			// CommitSelectionStrategy is Lexical, NewestTag, or SemVer. This
			// field is
			// optional.
			ignoreTags?: [...string]

			// IncludePaths is a list of selectors that designate paths in the
			// repository
			// that should trigger the production of new Freight when changes
			// are detected
			// therein. When specified, only changes in the identified paths
			// will trigger
			// Freight production. When not specified, changes in any path
			// will trigger
			// Freight production. Selectors may be defined using:
			// 1. Exact paths to files or directories (ex. "charts/foo")
			// 2. Glob patterns (prefix the pattern with "glob:"; ex.
			// "glob:*.yaml")
			// 3. Regular expressions (prefix the pattern with "regex:" or
			// "regexp:";
			// ex. "regexp:^.*\.yaml$")
			// Paths selected by IncludePaths may be unselected by
			// ExcludePaths. This
			// is a useful method for including a broad set of paths and then
			// excluding a
			// subset of them.
			includePaths?: [...string]

			// InsecureSkipTLSVerify specifies whether certificate
			// verification errors
			// should be ignored when connecting to the repository. This
			// should be enabled
			// only with great caution.
			insecureSkipTLSVerify?: bool

			// URL is the repository's URL. This is a required field.
			repoURL: strings.MinRunes(1) & {
				=~"(?:^(https?)://(?:([\\w-]+):(.+)@)?([\\w-]+(?:\\.[\\w-]+)*)(?::(\\d{1,5}))?(/.*)$)|(?:^([\\w-]+)@([\\w+]+(?:\\.[\\w-]+)*):(/?.*))"
			}

			// SemverConstraint specifies constraints on what new tagged
			// commits are
			// considered in determining the newest commit of interest. The
			// value in this
			// field only has any effect when the CommitSelectionStrategy is
			// SemVer. This
			// field is optional. When left unspecified, there will be no
			// constraints,
			// which means the latest semantically tagged commit will always
			// be used. Care
			// should be taken with leaving this field unspecified, as it can
			// lead to the
			// unanticipated rollout of breaking changes.
			semverConstraint?: string

			// StrictSemvers specifies whether only "strict" semver tags
			// should be
			// considered. A "strict" semver tag is one containing ALL of
			// major, minor,
			// and patch version components. This is enabled by default, but
			// only has any
			// effect when the CommitSelectionStrategy is SemVer. This should
			// be disabled
			// cautiously, as it creates the potential for any tag containing
			// numeric
			// characters only to be mistaken for a semver string containing
			// the major
			// version number only.
			strictSemvers: bool | *true
		}

		// Image describes a subscription to container image repository.
		image?: {
			// AllowTags is a regular expression that can optionally be used
			// to limit the
			// image tags that are considered in determining the newest
			// version of an
			// image. This field is optional.
			allowTags?: string

			// DiscoveryLimit is an optional limit on the number of image
			// references
			// that can be discovered for this subscription. The limit is
			// applied after
			// filtering images based on the AllowTags and IgnoreTags fields.
			// When left unspecified, the field is implicitly treated as if
			// its value
			// were "20". The upper limit for this field is 100.
			discoveryLimit?: int & <=100 & >=1 | *20

			// GitRepoURL optionally specifies the URL of a Git repository
			// that contains
			// the source code for the image repository referenced by the
			// RepoURL field.
			// When this is specified, Kargo MAY be able to infer and link to
			// the exact
			// revision of that source code that was used to build the image.
			gitRepoURL?: =~"^https?://(\\w+([\\.-]\\w+)*@)?\\w+([\\.-]\\w+)*(:[\\d]+)?(/.*)?$"

			// IgnoreTags is a list of tags that must be ignored when
			// determining the
			// newest version of an image. No regular expressions or glob
			// patterns are
			// supported yet. This field is optional.
			ignoreTags?: [...string]

			// ImageSelectionStrategy specifies the rules for how to identify
			// the newest version
			// of the image specified by the RepoURL field. This field is
			// optional. When
			// left unspecified, the field is implicitly treated as if its
			// value were
			// "SemVer".
			// Accepted values: Digest, Lexical, NewestBuild, SemVer
			imageSelectionStrategy?: "Digest" | "Lexical" | "NewestBuild" | "SemVer" | *"SemVer"

			// InsecureSkipTLSVerify specifies whether certificate
			// verification errors
			// should be ignored when connecting to the repository. This
			// should be enabled
			// only with great caution.
			insecureSkipTLSVerify?: bool

			// Platform is a string of the form <os>/<arch> that limits the
			// tags that can
			// be considered when searching for new versions of an image. This
			// field is
			// optional. When left unspecified, it is implicitly equivalent to
			// the
			// OS/architecture of the Kargo controller. Care should be taken
			// to set this
			// value correctly in cases where the image referenced by this
			// ImageRepositorySubscription will run on a Kubernetes node with
			// a different
			// OS/architecture than the Kargo controller. At present this is
			// uncommon, but
			// not unheard of.
			platform?: string

			// RepoURL specifies the URL of the image repository to subscribe
			// to. The
			// value in this field MUST NOT include an image tag. This field
			// is required.
			repoURL: strings.MinRunes(1) & {
				=~"^(\\w+([\\.-]\\w+)*(:[\\d]+)?/)?(\\w+([\\.-]\\w+)*)(/\\w+([\\.-]\\w+)*)*$"
			}

			// SemverConstraint specifies constraints on what new image
			// versions are
			// permissible. The value in this field only has any effect when
			// the
			// ImageSelectionStrategy is SemVer or left unspecified (which is
			// implicitly
			// the same as SemVer). This field is also optional. When left
			// unspecified,
			// (and the ImageSelectionStrategy is SemVer or unspecified),
			// there will be no
			// constraints, which means the latest semantically tagged version
			// of an image
			// will always be used. Care should be taken with leaving this
			// field
			// unspecified, as it can lead to the unanticipated rollout of
			// breaking
			// changes. Refer to Image Updater documentation for more details.
			// More info:
			// https://github.com/masterminds/semver#checking-version-constraints
			semverConstraint?: string

			// StrictSemvers specifies whether only "strict" semver tags
			// should be
			// considered. A "strict" semver tag is one containing ALL of
			// major, minor,
			// and patch version components. This is enabled by default, but
			// only has any
			// effect when the ImageSelectionStrategy is SemVer. This should
			// be disabled
			// cautiously, as it is not uncommon to tag container images with
			// short Git
			// commit hashes, which have the potential to contain numeric
			// characters only
			// and could be mistaken for a semver string containing the major
			// version
			// number only.
			strictSemvers: bool | *true
		}
	}] & [_, ...]
}
