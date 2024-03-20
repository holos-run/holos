package v1alpha1

// APIObjectMap is the shape of marshalled api objects returned from cue to the
// holos cli. A map is used to improve the clarity of error messages from cue.
type APIObjectMap map[string]map[string]string

// FileContentMap is a map of file names to file contents.
type FileContentMap map[string]string
