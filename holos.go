// Package holos defines types for the rest of the system.
package holos

// A PathCueMod is a string representing the absolute filesystem path of a cue
// module.  It is given a unique type so the API is clear.
type PathCueMod string

// A InstancePath is a string representing the absolute filesystem path of a
// holos instance.  It is given a unique type so the API is clear.
type InstancePath string

// TypeMeta represents the kind and version of a resource holos needs to
// process.  Useful to discriminate generated resources.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}
