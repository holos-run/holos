package v1alpha1

// Label is an arbitrary unique identifier.  Defined as a type for clarity and type checking.
type Label string

// Kind is a kubernetes api object kind. Defined as a type for clarity and type checking.
type Kind string

// APIObjectMap is the shape of marshalled api objects returned from cue to the
// holos cli. A map is used to improve the clarity of error messages from cue.
type APIObjectMap map[Kind]map[Label]string

// FileContentMap is a map of file names to file contents.
type FileContentMap map[string]string
