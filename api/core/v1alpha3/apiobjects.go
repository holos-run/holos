package v1alpha3

import "google.golang.org/protobuf/types/known/structpb"

// InternalLabel is an arbitrary unique identifier internal to holos itself.
// The holos cli is expected to never write a InternalLabel value to rendered
// output files, therefore use a [InternalLabel] when the identifier must be
// unique and internal.  Defined as a type for clarity and type checking.
//
// A InternalLabel is useful to convert a CUE struct to a list, for example
// producing a list of [APIObject] resources from an [APIObjectMap].  A CUE
// struct using InternalLabel keys is guaranteed to not lose data when rendering
// output because a InternalLabel is expected to never be written to the final
// output.
type InternalLabel string

// Kind is a kubernetes api object kind. Defined as a type for clarity and type
// checking.
type Kind string

// APIObject represents the most basic generic form of a single kubernetes api
// object.  Represented as a JSON object internally for compatibility between
// tools, for example loading from CUE.
type APIObject structpb.Struct

// APIObjectMap represents the marshalled yaml representation of kubernetes api
// objects.  Do not produce an APIObjectMap directly, instead use [APIObjects]
// to produce the marshalled yaml representation from CUE data, then provide the
// result to [Component].
type APIObjectMap map[Kind]map[InternalLabel]string

// APIObjects represents Kubernetes API objects defined directly from CUE code.
// Useful to mix in resources to any kind of [Component], for example
// adding an ExternalSecret resource to a [HelmChart].
//
// [Kind] must be the resource kind, e.g. Deployment or Service.
//
// [InternalLabel] is an arbitrary internal identifier to uniquely identify the resource
// within the context of a `holos` command.  Holos will never write the
// intermediate label to rendered output.
//
// Refer to [Component] which accepts an [APIObjectMap] field provided by
// [APIObjects].
type APIObjects struct {
	APIObjects   map[Kind]map[InternalLabel]APIObject `json:"apiObjects"`
	APIObjectMap APIObjectMap                         `json:"apiObjectMap"`
}
