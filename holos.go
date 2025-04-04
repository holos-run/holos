// Package holos defines types for the rest of the system.
package holos

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
