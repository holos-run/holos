package v1alpha1

// Platform represents a platform to manage.  A Platform resource tells holos
// which components to build.  The primary use case is to specify the cluster
// names, cluster types, and holos components to build.
type Platform struct {
	TypeMeta `json:",inline" yaml:",inline"`
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}
