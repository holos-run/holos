package v1alpha2

import "google.golang.org/protobuf/types/known/structpb"

type PlatformMetadata struct {
	// Name represents the Platform name.
	Name string `json:"name"`
}

// Platform represents a platform to manage.  A Platform resource informs holos
// which components to build.  The platform resource also acts as a container
// for the platform model form values provided by the PlatformService.  The
// primary use case is to collect the cluster names, cluster types, platform
// model, and holos components to build into one resource.
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

// PlatformSpec represents the specification of a Platform.  Think of a platform
// specification as a list of platform components to apply to a list of
// kubernetes clusters combined with the user-specified Platform Model.
type PlatformSpec struct {
	// Model represents the platform model holos gets from from the
	// PlatformService.GetPlatform rpc method and provides to CUE using a tag.
	Model structpb.Struct `json:"model"`
	// Components represents a list of holos components to manage.
	Components []PlatformSpecComponent `json:"components"`
}

// PlatformSpecComponent represents a holos component to build or render.
type PlatformSpecComponent struct {
	// Path is the path of the component relative to the platform root.
	Path string `json:"path"`
	// Cluster is the cluster name to provide when rendering the component.
	Cluster string `json:"cluster"`
}
