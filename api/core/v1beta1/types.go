// Package core contains schemas for a [Platform] and [TaskSet].  Holos takes
// a [Platform] as input, then iterates over each [Component] to produce a
// [TaskSet].  Holos merges all TaskSets into one platform-wide DAG and
// executes tasks in topological order to produce fully rendered manifests.
package core

// BuildContextTag represents the cue tag holos render component uses to inject
// the json representation of a [BuildContext] for use in a TaskSet.
const BuildContextTag string = "holos_build_context"

// ComponentNameTag represents the cue tag holos uses to inject a [Component]
// name from the holos render platform command to the holos render component
// command.
const ComponentNameTag string = "holos_component_name"

// ComponentPathTag represents the cue tag holos uses to inject a [Component]
// path relative to the cue module root from the holos render platform command
// to the holos render component command.
const ComponentPathTag string = "holos_component_path"

// ComponentLabelsTag represents the cue tag holos uses to inject the json
// representation of [Component] metadata labels from the holos render platform
// command to the holos render component command.
const ComponentLabelsTag string = "holos_component_labels"

// ComponentAnnotationsTag represents the tag holos uses to inject the json
// representation of [Component] metadata annotations from the holos render
// platform command to the holos render component command.
const ComponentAnnotationsTag = "holos_component_annotations"

//go:generate ../../../hack/gendoc

// TaskSet represents an implementation of the [rendered manifest pattern].
// A TaskSet replaces the deprecated v1alpha6 BuildPlan.  Each [Component]
// produces one TaskSet.  Holos merges all component TaskSets into one
// platform-wide DAG and executes tasks in topological order to produce fully
// rendered manifests.
//
// Holos uses CUE to construct a TaskSet.  Holos injects late binding values
// such as the build temp dir using the [BuildContext].
//
// [rendered manifest pattern]: https://akuity.io/blog/the-rendered-manifests-pattern
type TaskSet struct {
	// APIVersion represents the versioned schema of the resource.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"\"v1beta1\""`
	// Kind represents the type of the resource.
	Kind string `json:"kind" yaml:"kind" cue:"\"TaskSet\""`
	// Metadata represents data about the resource such as the Name.
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	// Spec specifies the desired state of the resource.
	Spec TaskSetSpec `json:"spec" yaml:"spec"`
	// BuildContext represents values injected by holos just before evaluating
	// a TaskSet, for example the tempDir used for the build.
	BuildContext BuildContext `json:"buildContext" yaml:"buildContext"`
}

// BuildContext represents build context values owned by the holos render
// component command.  End users should not manage context field values.  End
// users may reference fields from within CUE to refer to late binding concrete
// values defined just before holos executes a [TaskSet].
//
// Holos injects build context values by marshalling this struct to json through
// the holos_build_context cue tag.
//
// Example usage from cue to produce a [TaskSet]:
//
//	package holos
//
//	import (
//	  "encoding/json"
//	  "github.com/holos-run/holos/api/core/v1beta1:core"
//	)
//
//	_BuildContextJSON: string | *"{}" @tag(holos_build_context, type=string)
//	BuildContext: core.#BuildContext & json.Unmarshal(_BuildContextJSON)
//
//	holos: core.#TaskSet & {
//	  buildContext: BuildContext
//	  spec: tasks: {
//	    slice: {
//	      kind: "Command"
//	      inputs: ["resources.gen.yaml"]
//	      output: "components/slice"
//	      command: args: [
//	        "kubectl-slice",
//	        "-f",
//	        "\(buildContext.tempDir)/resources.gen.yaml",
//	        "-o",
//	        "\(buildContext.tempDir)/components/slice",
//	      ]
//	    }
//	  }
//	}
type BuildContext struct {
	// TempDir represents the temporary directory managed and owned by the holos
	// render component command for the execution of one TaskSet.  Multiple tasks
	// in the task set share this temporary directory and therefore should avoid
	// reading and writing into the same sub-directories as one another.
	TempDir string `json:"tempDir" yaml:"tempDir" cue:"string | *\"${TEMP_DIR_PLACEHOLDER}\""`
	// RootDir represents the fully qualified path to the platform root directory.
	// Useful to construct arguments for commands in TaskSet tasks.
	RootDir string `json:"rootDir" yaml:"rootDir" cue:"string | *\"${ROOT_DIR_PLACEHOLDER}\""`
	// LeafDir represents the cleaned path to the holos component relative to the
	// platform root.  Useful to construct arguments for commands in TaskSet
	// tasks.
	LeafDir string `json:"leafDir" yaml:"leafDir" cue:"string | *\"${LEAF_DIR_PLACEHOLDER}\""`
	// HolosExecutable represents the fully qualified path to the holos
	// executable.  Useful to execute tools embedded as subcommands such as holos
	// cue vet.
	HolosExecutable string `json:"holosExecutable" yaml:"holosExecutable" cue:"string | *\"holos\""`
}

// TaskSetSpec represents the specification of the [TaskSet].
type TaskSetSpec struct {
	// Tasks represents the tasks for holos to execute, keyed by name.  Tasks are
	// structs, not lists, so TaskSets compose by CUE unification.
	Tasks map[string]Task `json:"tasks" yaml:"tasks"`
	// Disabled causes the holos render platform command to skip the TaskSet.
	Disabled bool `json:"disabled,omitempty" yaml:"disabled,omitempty"`
}

// Task represents one unit of work in a [TaskSet].  Task unifies the v1alpha6
// Generator, Transformer, and Validator concepts.  A task declares the
// artifact-store paths it consumes and produces; the executor derives DAG
// edges from those declarations.
//
// Exactly one of the kind-specific config fields must be set, matching Kind:
//
//  1. [Resources] - Export Kubernetes resources defined in CUE.
//  2. [Helm] - Render a Helm chart.
//  3. [File] - Read a file from the component directory.
//  4. [Kustomize] - Patch and transform prior outputs.
//  5. [Join] - Concatenate prior outputs.
//  6. [Command] - Execute a user defined command.
//  7. [Artifact] - Write the final artifact (sink).
type Task struct {
	// Kind discriminates the task behavior.
	Kind string `json:"kind" yaml:"kind" cue:"\"Resources\" | \"Helm\" | \"File\" | \"Kustomize\" | \"Join\" | \"Command\" | \"Artifact\""`
	// DependsOn declares tasks that must complete before this task runs, keyed
	// by task name or canonical ID — a struct, not a list, so mixins compose
	// ordering edges by unification.  Use for ordering constraints with no data
	// flow; data-flow edges are derived from Inputs and Output.
	DependsOn map[string]Dependency `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`
	// Inputs are artifact-store paths consumed by the task.
	Inputs []FileOrDirectoryPath `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	// Output is the artifact-store path produced by the task.  Output values are
	// write-once: it is an error for two tasks to declare the same Output within
	// one TaskSet.  The platform merge namespaces store paths by component,
	// extending the rule platform-wide.
	Output FileOrDirectoryPath `json:"output,omitempty" yaml:"output,omitempty"`
	// Resources task config.  Ignored unless kind is Resources.
	Resources Resources `json:"resources,omitempty" yaml:"resources,omitempty"`
	// Helm task config.  Ignored unless kind is Helm.
	Helm Helm `json:"helm,omitempty" yaml:"helm,omitempty"`
	// File task config.  Ignored unless kind is File.
	File File `json:"file,omitempty" yaml:"file,omitempty"`
	// Kustomize task config.  Ignored unless kind is Kustomize.
	Kustomize Kustomize `json:"kustomize,omitempty" yaml:"kustomize,omitempty"`
	// Join task config.  Ignored unless kind is Join.
	Join Join `json:"join,omitempty" yaml:"join,omitempty"`
	// Command task config.  Ignored unless kind is Command.
	Command Command `json:"command,omitempty" yaml:"command,omitempty"`
	// Artifact task config.  Ignored unless kind is Artifact.
	Artifact Artifact `json:"artifact,omitempty" yaml:"artifact,omitempty"`
}

// Dependency represents one explicit ordering edge declared in
// [Task.DependsOn].  It is deliberately empty — the edge is the struct key —
// so future fields (for example an optional edge) may be added without
// breaking composition.
type Dependency struct{}

// Command represents a [Task] implemented by executing an user defined system
// command.  Command is a first-class Task kind in v1beta1.  Commands execute
// with the working directory set to the platform root.
//
// A command with an output generates or transforms; a command with only inputs
// validates, gating downstream tasks through [Task.DependsOn] edges.
type Command struct {
	// DisplayName of the command.  The basename of args[0] is used if empty.
	DisplayName string `json:"displayName,omitempty" yaml:"displayName,omitempty"`
	// Args represents the argument vector passed to the os to execute the
	// command.
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
	// Stdin names a task input wired to the command's standard input.  Must be
	// one of the task's declared Inputs.
	Stdin FileOrDirectoryPath `json:"stdin,omitempty" yaml:"stdin,omitempty"`
	// IsStdoutOutput captures the command stdout as the task Output if true.
	IsStdoutOutput bool `json:"isStdoutOutput,omitempty" yaml:"isStdoutOutput,omitempty"`
}

// Artifact represents the sink [Task] kind.  It writes its single input from
// the artifact store to the final artifact path once every task it depends on
// has completed successfully.
type Artifact struct {
	// Path represents the final artifact path relative to the write-to
	// directory (deploy by default).  Defaults to the task's single input path
	// when empty.
	Path FileOrDirectoryPath `json:"path,omitempty" yaml:"path,omitempty"`
}

// Resource represents one kubernetes api object.
type Resource map[string]any

// Resources represents Kubernetes resources.  Most commonly used to mix
// resources into the [TaskSet] generated from CUE, but may be generated from
// elsewhere.  Resources are stored as a two level struct.  The top level key
// is the Kind of resource, e.g. Namespace or Deployment.  The second level key
// is an arbitrary [InternalLabel].  The third level is a map[string]any
// representing the [Resource].
type Resources map[Kind]map[InternalLabel]Resource

// File represents a simple single file copy [Task].  Useful with a
// [Kustomize] task to process plain manifest files stored in the component
// directory.  Multiple File tasks may be used to transform multiple resources.
type File struct {
	// Source represents a file sub-path relative to the component path.
	Source FilePath `json:"source" yaml:"source"`
}

// Helm represents a [Task] that renders a helm [Chart].
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
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
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

// Join represents a [Task] using [bytes.Join] to concatenate multiple inputs
// into one output with a separator.  Useful for combining output from [Helm]
// and [Resources] tasks together into one output when [Kustomize] is otherwise
// unnecessary.
//
// [bytes.Join]: https://pkg.go.dev/bytes#Join
type Join struct {
	Separator string `json:"separator,omitempty" yaml:"separator,omitempty"`
}

// Kustomize represents a kustomization [Task] to patch and transform prior
// task outputs.
type Kustomize struct {
	// Kustomization represents the decoded kustomization.yaml file
	Kustomization Kustomization `json:"kustomization" yaml:"kustomization"`
	// Files holds file contents for kustomize, e.g. patch files.
	Files FileContentMap `json:"files,omitempty" yaml:"files,omitempty"`
}

// Kustomization represents a kustomization.yaml file for use with the
// [Kustomize] [Task].  Untyped to avoid tightly coupling holos to kubectl
// versions which was a problem for the Flux maintainers.  Type checking is
// expected to happen in CUE against the kubectl version the user prefers.
type Kustomization map[string]any

// FileContentMap represents a mapping of file paths to file contents.
type FileContentMap map[FilePath]FileContent

// FilePath represents a file path.
type FilePath string

// FileOrDirectoryPath represents a file or a directory path.
type FileOrDirectoryPath string

// FileContent represents file contents.
type FileContent string

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
// Holos iterates over the [Component] collection producing a [TaskSet] for
// each, which holos then executes to render manifests.
//
// Inspect a Platform resource holos would process by executing:
//
//	cue export --out yaml ./platform
type Platform struct {
	// APIVersion represents the versioned schema of this resource.
	APIVersion string `json:"apiVersion" yaml:"apiVersion" cue:"string | *\"v1beta1\""`
	// Kind is a string value representing the resource.
	Kind string `json:"kind" yaml:"kind" cue:"\"Platform\""`
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

// Component represents the complete context necessary to produce a [TaskSet]
// from a path containing parameterized CUE configuration.
type Component struct {
	// Name represents the name of the component. Injected as the tag variable
	// "holos_component_name".
	Name string `json:"name" yaml:"name"`
	// Path represents the path of the component relative to the platform root.
	// Injected as the tag variable "holos_component_path".
	Path string `json:"path" yaml:"path"`
	// Parameters represent user defined input variables to produce various
	// [TaskSet] resources from one component path.  Injected as CUE @tag
	// variables.  Parameters with a "holos_" prefix are reserved for use by the
	// Holos Authors.  Multiple environments are a prime example of an input
	// parameter that should always be user defined, never defined by Holos.
	Parameters map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// Labels represent selector labels for the component.  Holos copies Labels
	// from the Component to the resulting TaskSet.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// Annotations represents arbitrary non-identifying metadata.  Use the
	// `app.holos.run/description` to customize the log message of each TaskSet.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
