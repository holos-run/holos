package holos

import (
	"context"

	"cuelang.org/go/cue"
)

// Platform is the interface all Platform versions must support.
type Platform interface {
	Load(cue.Value) error
	Select(...Selector) []Component
	Export(Encoder) error
}

// Component represents the information needed to produce a BuildPlan from a
// platform component.
type Component interface {
	Describe() string
	Path() string
	Tags() []string
	WriteTo() string
	Labels() Labels
}

// BuildPlan is the interface to load, build, and export a BuildPlan from a cue
// value.
type BuildPlan interface {
	Load(cue.Value) error
	Build(context.Context) error
	Export(Encoder) error
}

// Encoder encodes to json or yaml.  Concrete values must be safe for concurrent
// execution of Encode().  Refer to [NewEncoder] to get an Encoder from a flag
// value.
type Encoder interface {
	Encode(any) error
	Close() error
}
