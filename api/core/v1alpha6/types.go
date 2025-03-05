// Package core contains schemas for a [Platform] and [BuildPlan].  Holos takes
// a [Platform] as input, then iterates over each [Component] to produce a
// [BuildPlan].  Holos processes the [BuildPlan] to produce fully rendered
// manifests, each an [Artifact].
package core

//go:generate ../../../hack/gendoc

// BuildPlan represents an implementation of the [rendered manifest pattern].
// Holos processes a BuildPlan to produce one or more [Artifact] output files.
// BuildPlan artifact files usually contain Kubernetes manifests, but they may
// have any content.
//
// A BuildPlan usually produces two artifacts.  One artifact contains a manifest
// of resources.  A second artifact contains a GitOps resource to manage the
// first, usually an ArgoCD Application resource.
//
// Holos uses CUE to construct a BuildPlan.  A future enhancement will support
// user defined executables providing a BuildPlan to Holos in the style of an
// [external credential provider].
//
// [rendered manifest pattern]: https://akuity.io/blog/the-rendered-manifests-pattern
// [external credential provider]: https://github.com/kubernetes/enhancements/blob/313ad8b59c80819659e1fbf0f165230f633f2b22/keps/sig-auth/541-external-credential-providers/README.md
type BuildPlan struct {
	// Kind represents the type of the resource.
	Kind string `json:"kind" yaml:"kind" cue:"\"BuildPlan\""`
	// APIVersion represents the versioned schema of the resource.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1alpha6\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	// Spec specifies the desired state of the resource.
	Spec BuildPlanSpec `json:"spec" yaml:"spec"`
	// Context represents build context values owned by the holos render component
	// command.  End users should not manage context field values.  End users may
	// reference context fields from within CUE to refer to late binding concrete
	// values defined just before holos executes the build plan.
	Context BuildContext `json:"context" yaml:"context"`
}

// BuildContext represents build context values provided by the holos render
// component command.  These values are expected to be randomly generated and
// late binding, meaning they cannot be known ahead of time in a static
// configuration.  As such, CUE configuration may refer to the values here which
// will be populated by holos when the final build plan is exported from CUE.
type BuildContext struct {
	// TempDir represents the temporary directory managed and owned by the holos
	// render component command for the execution of one BuildPlan.  Multiple
	// tasks in the build plan share this temporary directory and therefore should
	// avoid reading and writing into the same sub-directories as one another.
	TempDir string `json:"tempDir" yaml:"tempDir" cue:"string | *\"${TEMP_DIR_PLACEHOLDER}\""`
}

// BuildPlanSpec represents the specification of the [BuildPlan].
type BuildPlanSpec struct {
	// Artifacts represents the artifacts for holos to build.
	Artifacts []Artifact `json:"artifacts" yaml:"artifacts"`
	// Disabled causes the holos cli to disregard the build plan.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// Artifact represents one fully rendered manifest produced by a [Transformer]
// sequence, which transforms a [Generator] collection.  A [BuildPlan] produces
// an [Artifact] collection.
//
// Each Artifact produces one manifest file or directory artifact.  Generator
// Output values are used as Transformer Inputs.  The Output field of the final
// [Transformer] should have the same value as the Artifact field.
//
// When there is more than one [Generator] there should be at least one
// [Transformer] to combine outputs into one Artifact file, or the final
// artifact should be a directory containing the outputs of the generators.  If
// there is a single Generator, it may directly produce the Artifact output.
//
// An Artifact is processed concurrently with other artifacts in the same
// [BuildPlan].  One Artifact must not use an output of another Artifact as an
// input.  Each [Generator] within an artifact also runs concurrently with
// generators of the same artifact.  Each [Transformer] is executed sequentially
// starting after all generators have completed.
//
// Output fields are write-once.  It is an error for multiple Generators or
// Transformers to produce the same Output value within the context of a
// [BuildPlan].
//
// When directories are used as inputs or outputs, they behave similar to how
// `git` works with directories.  When the output field references a directory,
// all files within the directory are recursively stored using their relative
// path as a key.  Similar to git add .  When the input field references an
// absent file, a / is appended and the resulting value is used as a prefix
// match against all previous task outputs.
type Artifact struct {
	Artifact     FileOrDirectoryPath `json:"artifact,omitempty" yaml:"artifact,omitempty"`
	Generators   []Generator         `json:"generators,omitempty" yaml:"generators,omitempty"`
	Transformers []Transformer       `json:"transformers,omitempty" yaml:"transformers,omitempty"`
	Validators   []Validator         `json:"validators,omitempty" yaml:"validators,omitempty"`
	Skip         bool                `json:"skip,omitempty" yaml:"skip,omitempty"`
}

// Generator generates Kubernetes resources.  [Helm] and [Resources] are the
// most commonly used, often paired together to mix-in resources to an
// unmodified Helm chart.  A simple [File] generator is also available for use
// with the [Kustomize] transformer.
//
// Each Generator in an [Artifact] must have a distinct Output value for a
// [Transformer] to reference.
//
//  1. [Resources] - Generates resources from CUE code.
//  2. [Helm] - Generates rendered yaml from a [Chart].
//  3. [File] - Generates data by reading a file from the component directory.
//  4. [Command] - Generates data by executing an user defined command.
type Generator struct {
	// Kind represents the kind of generator.  Must be Resources, Helm, or File.
	Kind string `json:"kind" yaml:"kind" cue:"\"Resources\" | \"Helm\" | \"File\" | \"Command\""`
	// Output represents a file for a Transformer or Artifact to consume.
	Output FileOrDirectoryPath `json:"output" yaml:"output"`
	// Resources generator. Ignored unless kind is Resources.  Resources are
	// stored as a two level struct.  The top level key is the Kind of resource,
	// e.g. Namespace or Deployment.  The second level key is an arbitrary
	// InternalLabel.  The third level is a map[string]any representing the
	// Resource.
	Resources Resources `json:"resources,omitempty" yaml:"resources,omitempty"`
	// Helm generator. Ignored unless kind is Helm.
	Helm Helm `json:"helm,omitempty" yaml:"helm,omitempty"`
	// File generator. Ignored unless kind is File.
	File File `json:"file,omitempty" yaml:"file,omitempty"`
	// Command generator. Ignored unless kind is Command.
	Command Command `json:"command,omitempty" yaml:"command,omitempty"`
}

// Resource represents one kubernetes api object.
type Resource map[string]any

// Resources represents Kubernetes resources.  Most commonly used to mix
// resources into the [BuildPlan] generated from CUE, but may be generated from
// elsewhere.
type Resources map[Kind]map[InternalLabel]Resource

// File represents a simple single file copy [Generator].  Useful with a
// [Kustomize] [Transformer] to process plain manifest files stored in the
// component directory.  Multiple File generators may be used to transform
// multiple resources.
type File struct {
	// Source represents a file sub-path relative to the component path.
	Source FilePath `json:"source" yaml:"source"`
}

// Helm represents a [Chart] manifest [Generator].
type Helm struct {
	// Chart represents a helm chart to manage.
	Chart Chart `json:"chart" yaml:"chart"`
	// Values represents values for holos to marshal into values.yaml when
	// rendering the chart.  Values follow ValueFiles when both are provided.
	Values Values `json:"values" yaml:"values"`
	// ValueFiles represents hierarchial value files passed in order to the helm
	// template -f flag.  Useful for migration from an ApplicationSet.  Use Values
	// instead.  ValueFiles precede Values when both are provided.
	ValueFiles []ValueFile `json:"valueFiles,omitempty" yaml:"valueFiles,omitempty"`
	// EnableHooks enables helm hooks when executing the `helm template` command.
	EnableHooks bool `json:"enableHooks,omitempty" yaml:"enableHooks,omitempty"`
	// Namespace represents the helm namespace flag
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// APIVersions represents the helm template --api-versions flag
	APIVersions []string `json:"apiVersions,omitempty" yaml:"apiVersions,omitempty"`
	// KubeVersion represents the helm template --kube-version flag
	KubeVersion string `json:"kubeVersion,omitempty" yaml:"kubeVersion,omitempty"`
}

// ValueFile represents one Helm value file produced from CUE.
type ValueFile struct {
	// Name represents the file name, e.g. "region-values.yaml"
	Name string `json:"name" yaml:"name"`
	// Kind is a discriminator.
	Kind string `json:"kind" yaml:"kind" cue:"\"Values\""`
	// Values represents values for holos to marshal into the file name specified
	// by Name when rendering the chart.
	Values Values `json:"values,omitempty" yaml:"values,omitempty"`
}

// Values represents [Helm] Chart values generated from CUE.
type Values map[string]any

// Chart represents a [Helm] Chart.
type Chart struct {
	// Name represents the chart name.
	Name string `json:"name" yaml:"name"`
	// Version represents the chart version.
	Version string `json:"version" yaml:"version"`
	// Release represents the chart release when executing helm template.
	Release string `json:"release" yaml:"release"`
	// Repository represents the repository to fetch the chart from.
	Repository Repository `json:"repository,omitempty" yaml:"repository,omitempty"`
}

// Repository represents a [Helm] [Chart] repository.
//
// The Auth field is useful to configure http basic authentication to the Helm
// repository.  Holos gets the username and password from the environment
// variables represented by the Auth field.
type Repository struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
	Auth Auth   `json:"auth,omitempty" yaml:"auth,omitempty"`
}

// Auth represents environment variable names containing auth credentials.
type Auth struct {
	Username AuthSource `json:"username" yaml:"username"`
	Password AuthSource `json:"password" yaml:"password"`
}

// AuthSource represents a source for the value of an [Auth] field.
type AuthSource struct {
	Value   string `json:"value,omitempty" yaml:"value,omitempty"`
	FromEnv string `json:"fromEnv,omitempty" yaml:"fromEnv,omitempty"`
}

// Transformer combines multiple inputs from prior [Generator] or [Transformer]
// outputs into one output.  [Kustomize] is the most commonly used transformer.
// A simple [Join] is also supported for use with plain manifest files.
//
//  1. [Kustomize] - Patch and transform the output from prior generators or
//     transformers.  See [Introduction to Kustomize].
//  2. [Join] - Concatenate multiple prior outputs into one output.
//  3. [Command] - Transforms data by executing an user defined command.
//
// [Introduction to Kustomize]: https://kubectl.docs.kubernetes.io/guides/config_management/introduction/
type Transformer struct {
	// Kind represents the kind of transformer. Must be Kustomize, or Join.
	Kind string `json:"kind" yaml:"kind" cue:"\"Kustomize\" | \"Join\" | \"Command\""`
	// Inputs represents the files to transform. The Output of prior Generators
	// and Transformers.
	Inputs []FileOrDirectoryPath `json:"inputs" yaml:"inputs"`
	// Output represents a file or directory for a subsequent Transformer or
	// Artifact to consume.
	Output FileOrDirectoryPath `json:"output" yaml:"output"`
	// Kustomize transformer. Ignored unless kind is Kustomize.
	Kustomize Kustomize `json:"kustomize,omitempty" yaml:"kustomize,omitempty"`
	// Join transformer. Ignored unless kind is Join.
	Join Join `json:"join,omitempty" yaml:"join,omitempty"`
	// Command transformer. Ignored unless kind is Command.
	Command Command `json:"command,omitempty" yaml:"command,omitempty"`
}

// Join represents a [Transformer] using [bytes.Join] to concatenate multiple
// inputs into one output with a separator.  Useful for combining output from
// [Helm] and [Resources] together into one [Artifact] when [Kustomize] is
// otherwise unnecessary.
//
// [bytes.Join]: https://pkg.go.dev/bytes#Join
type Join struct {
	Separator string `json:"separator,omitempty" yaml:"separator,omitempty"`
}

// Kustomize represents a kustomization [Transformer].
type Kustomize struct {
	// Kustomization represents the decoded kustomization.yaml file
	Kustomization Kustomization `json:"kustomization" yaml:"kustomization"`
	// Files holds file contents for kustomize, e.g. patch files.
	Files FileContentMap `json:"files,omitempty" yaml:"files,omitempty"`
}

// Kustomization represents a kustomization.yaml file for use with the
// [Kustomize] [Transformer].  Untyped to avoid tightly coupling holos to
// kubectl versions which was a problem for the Flux maintainers.  Type checking
// is expected to happen in CUE against the kubectl version the user prefers.
type Kustomization map[string]any

// FileContentMap represents a mapping of file paths to file contents.
type FileContentMap map[FilePath]FileContent

// FilePath represents a file path.
type FilePath string

// FileOrDirectoryPath represents a file or a directory path.
type FileOrDirectoryPath string

// FileContent represents file contents.
type FileContent string

// Validator validates files.  Useful to validate an [Artifact] prior to writing
// it out to the final destination.  Holos may execute validators concurrently.
// See the [validators] tutorial for an end to end example.
//
// [validators]: https://holos.run/docs/v1alpha6/tutorial/validators/
type Validator struct {
	// Kind represents the kind of transformer. Must be Kustomize, or Join.
	Kind string `json:"kind" yaml:"kind" cue:"\"Command\""`
	// Inputs represents the files to validate.  Usually the final Artifact.
	Inputs []FileOrDirectoryPath `json:"inputs" yaml:"inputs"`
	// Command represents a validation command.  Ignored unless kind is Command.
	Command Command `json:"command,omitempty" yaml:"command,omitempty"`
}

// Command represents a task implemented as a generic system command.  A task is
// defined as a [Generator], [Transformer], or [Validator].
type Command struct {
	// DisplayName of the command.  The basename of args[0] is used if empty.
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	// Args represents the argument vector passed to the system to execute the
	// command.
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
	// Env represents environment variables to set in the command context.
	Env []EnvVar `json:"env,omitempty" yaml:"env,omitempty"`
	// Stdout captures the command standard output for use as the task output.
	// Set to false for commands that write output to files.
	Stdout bool `json:"stdout,omitempty" yaml:"stdout,omitempty"`
}

// EnvVar represents the configuration of an environment variable within the
// context of a BuildPlan task.
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name" yaml:"name"`
	// Kind represents a discriminator.
	Kind string `json:"kind" yaml:"kind" cue:"\"Value\" | \"ValueFrom\""`
	// Value represents the concrete value of the named environment variable.
	// Ignored unless kind is Value.
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
	// ValueFrom represents the source for the named environment variable's value.
	// Ignored unless kind is ValueFrom.
	ValueFrom EnvVarSource `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

// EnvVarSource represents a source for the value of an EnvVar.
type EnvVarSource struct {
	// Kind represents a discriminator.
	Kind string `json:"kind" yaml:"kind" cue:"\"EnvRef\""`
	// EnvRef represents a reference to an environment variable.  Ignored unless
	// kind is EnvRef.
	EnvRef EnvRef `json:"envRef,omitempty" yaml:"envRef,omitempty"`
}

// EnvRef represents a reference to a value located in the environment.
type EnvRef struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name" yaml:"name"`
}

// InternalLabel is an arbitrary unique identifier internal to holos itself.
// The holos cli is expected to never write a InternalLabel value to rendered
// output files, therefore use a InternalLabel when the identifier must be
// unique and internal.  Defined as a type for clarity and type checking.
type InternalLabel string

// Kind is a discriminator. Defined as a type for clarity and type checking.
type Kind string

// Metadata represents data about the resource such as the Name.
type Metadata struct {
	// Name represents the resource name.
	Name string `json:"name" yaml:"name"`
	// Labels represents a resource selector.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// Annotations represents arbitrary non-identifying metadata.  For example
	// holos uses the `app.holos.run/description` annotation to log resources in a
	// user customized way.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Platform represents a platform to manage.  A Platform specifies a [Component]
// collection and integrates the components together into a holistic platform.
// Holos iterates over the [Component] collection producing a [BuildPlan] for
// each, which holos then executes to render manifests.
//
// Inspect a Platform resource holos would process by executing:
//
//	cue export --out yaml ./platform
type Platform struct {
	// Kind is a string value representing the resource.
	Kind string `json:"kind" yaml:"kind" cue:"\"Platform\""`
	// APIVersion represents the versioned schema of this resource.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1alpha6\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata" yaml:"metadata"`

	// Spec represents the platform specification.
	Spec PlatformSpec `json:"spec" yaml:"spec"`
}

// PlatformSpec represents the platform specification.
type PlatformSpec struct {
	// Components represents a collection of holos components to manage.
	Components []Component `json:"components" yaml:"components"`
}

// Component represents the complete context necessary to produce a [BuildPlan]
// from a path containing parameterized CUE configuration.
type Component struct {
	// Name represents the name of the component. Injected as the tag variable
	// "holos_component_name".
	Name string `json:"name" yaml:"name"`
	// Path represents the path of the component relative to the platform root.
	// Injected as the tag variable "holos_component_path".
	Path string `json:"path" yaml:"path"`
	// Instances represents additional cue instance paths to unify with Path.
	// Useful to unify data files into a component BuildPlan.  Added in holos
	// 0.101.7.
	Instances []Instance `json:"instances,omitempty" yaml:"instances,omitempty"`
	// WriteTo represents the holos render component --write-to flag.  If empty,
	// the default value for the --write-to flag is used.
	WriteTo string `json:"writeTo,omitempty" yaml:"writeTo,omitempty"`
	// Parameters represent user defined input variables to produce various
	// [BuildPlan] resources from one component path.  Injected as CUE @tag
	// variables.  Parameters with a "holos_" prefix are reserved for use by the
	// Holos Authors.  Multiple environments are a prime example of an input
	// parameter that should always be user defined, never defined by Holos.
	Parameters map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// Labels represent selector labels for the component.  Copied to the
	// resulting BuildPlan.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// Annotations represents arbitrary non-identifying metadata.  Use the
	// `app.holos.run/description` to customize the log message of each BuildPlan.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Instance represents a data instance to unify with the configuration.
//
// Useful to unify json and yaml files with cue configuration files for
// integration with other tools.  For example, executing holos render platform
// from a pull request workflow after [Kargo] executes the [yaml update] and
// [git wait for pr] promotion steps.
//
// [Kargo]: https://docs.kargo.io/
// [yaml update]: https://docs.kargo.io/references/promotion-steps#yaml-update
// [git wait for pr]: https://docs.kargo.io/references/promotion-steps#git-wait-for-pr
type Instance struct {
	// Kind is a discriminator.
	Kind string `json:"kind" yaml:"kind" cue:"\"ExtractYAML\""`
	// Ignored unless kind is ExtractYAML.
	ExtractYAML ExtractYAML `json:"extractYAML,omitempty" yaml:"extractYAML,omitempty"`
}

// ExtractYAML represents a cue data instance encoded as yaml or json. If Path
// refers to a directory all files in the directory are extracted
// non-recursively.  Otherwise, path must refer to a file.
type ExtractYAML struct {
	Path string `json:"path" yaml:"path"`
}
