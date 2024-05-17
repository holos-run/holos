package v1alpha1

import "google.golang.org/protobuf/types/known/structpb"

// Platform represents a platform to manage.  A Platform resource informs holos
// which components to build.  The platform resource also acts as a container
// for the platform model form values provided by the PlatformService.  The
// primary use case is to collect the cluster names, cluster types, platform
// model, and holos components to build into one resource.
type Platform struct {
	TypeMeta `json:",inline" yaml:",inline"`
	Metadata ObjectMeta   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec     PlatformSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// PlatformSpec represents the platform build plan specification.
type PlatformSpec struct {
	// Model represents the platform model holos gets from from the
	// holos.platform.v1alpha1.PlatformService.GetPlatform method and provides to
	// CUE using a tag.
	Model structpb.Struct `json:"model,omitempty" yaml:"model,omitempty"`
}
