package holos

import (
	"context"

	"cuelang.org/go/cue"
)

// Exporter is the interface used to export Platform, Component, and BuildPlan
// resources.  Useful for the show subcommand.
type Exporter interface {
	Export(Encoder) error
}

// Platform is the interface all Platform versions must support.
//
// Concrete values must preserve the relative ordering of components when
// filtering with Select.
type Platform interface {
	Load(cue.Value) error
	Select(...Selector) []Component
	Export(Encoder) error
}

// Component represents the information needed to produce a BuildPlan from a
// platform component.
type Component interface {
	// Describe represents a description of the component for user feedback.
	Describe() string
	// Path represents the component path relative to the platform root.
	Path() string
	// Tags represents the Component name, labels, annotations injected into cue
	// for BuildPlan export.
	Tags() ([]string, error)
}

// BuildPlan is the interface to load, build, and export a BuildPlan from a cue
// value.
type BuildPlan interface {
	Load(cue.Value) error
	Build(context.Context) error
	Export(idx int, encoder OrderedEncoder) error
}

// Encoder encodes to json or yaml.  Concrete values must be safe for concurrent
// execution.  Use [NewEncoder] to obtain a json or yaml encoder.
type Encoder interface {
	Encode(any) error
	Close() error
}

// OrderedEncoder encodes in sequential order from idx 0.  Concrete values must
// be safe for concurrent execution.  Use [NewSequentialEncoder] to obtain a
// json or yaml encoder.
type OrderedEncoder interface {
	Encode(idx int, v any) error
	Close() error
}
