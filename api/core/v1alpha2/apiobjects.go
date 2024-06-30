package v1alpha2

import "google.golang.org/protobuf/types/known/structpb"

// Label is an arbitrary unique identifier.  Defined as a type for clarity and type checking.
type Label string

// Kind is a kubernetes api object kind. Defined as a type for clarity and type checking.
type Kind string

// APIObjectMap represents the marshalled yaml representation of kubernetes api
// objects.  Do not produce an APIObjectMap directly, instead use [APIObjects]
// to produce the marshalled yaml representation from CUE data.
//
// Example:
//
//	# CUE
//	apiObjectMap: (#APIObjects & {apiObjects: Resources}).apiObjectMap
type APIObjectMap map[Kind]map[Label]string

// APIObjects represents kubernetes api objects to apply to the api server.
// Useful to mix in resources to each HolosComponent type, for example adding an
// ExternalSecret to a HelmChart HolosComponent.
//
// Kind must be the resource kind, e.g. Deployment or Service.
//
// Label is an arbitrary internal identifier to uniquely identify the resource
// within the context of a `holos` command.  Holos will never write the
// intermediate label to rendered output.
//
// Refer to [HolosComponent] which accepts an [APIObjectMap] field provided by
// [APIObjects].
type APIObjects struct {
	// APIObjects represents Kubernetes API objects defined directly from CUE
	// code.  Useful to mix in resources, for example adding an ExternalSecret
	// resource to a HelmChart HolosComponent.
	APIObjects map[Kind]map[Label]structpb.Struct `json:"apiObjects"`
	// APIObjectMap represents the marshalled yaml representation of APIObjects,
	// useful to inspect the rendered representation of the resource which will be
	// sent to the kubernetes API server.
	APIObjectMap APIObjectMap `json:"apiObjectMap"`
}
