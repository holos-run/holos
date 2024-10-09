// Package holos defines types for the rest of the system.
package holos

import "context"

// A PathCueMod is a string representing the absolute filesystem path of a cue
// module.  It is given a unique type so the API is clear.
type PathCueMod string

// A InstancePath is a string representing the absolute filesystem path of a
// holos instance.  It is given a unique type so the API is clear.
type InstancePath string

// FilePath represents the path of a file relative to the current working
// directory of holos at runtime.
type FilePath string

// FileContent represents the contents of a file as a string.
type FileContent string

// TypeMeta represents the kind and version of a resource holos needs to
// process.  Useful to discriminate generated resources.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

// Builder builds file artifacts.
type Builder interface {
	Build(context.Context, ArtifactMap) error
}

// ArtifactMap sets and gets data for file artifacts.
//
// Concrete values must ensure Set is write once, returning an error if a given
// FilePath was previously Set.  Concrete values must be safe for concurrent
// reads and writes.
type ArtifactMap interface {
	Get(FilePath) ([]byte, bool)
	Set(FilePath, []byte) error
}
