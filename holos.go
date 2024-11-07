// Package holos defines types for the rest of the system.
package holos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"cuelang.org/go/cue"
	"github.com/holos-run/holos/internal/errors"
)

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
	Get(path string) (data []byte, ok bool)
	Set(path string, data []byte) error
	Save(dir, path string) error
}

// Discriminator is useful to discriminate by type meta, the kind and api
// version of something.
type Discriminator interface {
	Discriminate(ctx context.Context) (TypeMeta, error)
}

type Unifier interface {
	Unify(ctx context.Context) (BuildData, error)
}

// BuildData represents the data necessary to produce a build plan.  It is a
// convenience wrapper to store relevant fields to inform the user.
type BuildData struct {
	Value        cue.Value
	ModuleRoot   string
	InstancePath InstancePath
	Dir          string
}

func (bd *BuildData) TypeMeta() (tm TypeMeta, err error) {
	v, err := bd.value()
	if err != nil {
		return tm, errors.Wrap(err)
	}

	kind := v.LookupPath(cue.ParsePath("kind"))
	if err := kind.Err(); err != nil {
		return tm, errors.Wrap(err)
	}
	if tm.Kind, err = kind.String(); err != nil {
		return tm, errors.Wrap(err)
	}

	version := v.LookupPath(cue.ParsePath("apiVersion"))
	if err := version.Err(); err != nil {
		return tm, errors.Wrap(err)
	}
	if tm.APIVersion, err = version.String(); err != nil {
		return tm, errors.Wrap(err)
	}

	return
}

func (bd *BuildData) Decoder() (*json.Decoder, error) {
	v, err := bd.value()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	jsonBytes, err := v.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.DisallowUnknownFields()
	return decoder, nil
}

func (bd *BuildData) value() (v cue.Value, err error) {
	v = bd.Value.LookupPath(cue.ParsePath("holos"))
	if err := v.Err(); err != nil {
		if strings.HasPrefix(err.Error(), "field not found") {
			slog.Warn(fmt.Sprintf("%s: deprecated usage: nest output under holos: %s", err, bd.Dir), "err", err)
			v = bd.Value
			return v, nil
		}
		err = errors.Wrap(err)
		return v, err
	}
	return
}
