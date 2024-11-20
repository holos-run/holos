package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/util"
)

func LoadInstance(path string, tags []string) (*Instance, error) {
	root, leaf, err := util.FindRootLeaf(path)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	cfg := &load.Config{
		Dir:        root,
		ModuleRoot: root,
		Tags:       tags,
	}

	ctx := cuecontext.New()

	instances := load.Instances([]string{leaf}, cfg)
	values, err := ctx.BuildInstances(instances)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	inst := &Instance{
		path:  leaf,
		ctx:   ctx,
		cfg:   cfg,
		value: values[0],
	}

	return inst, nil
}

// Instance represents a cue instance to build.  Use LoadInstance to create a
// new Instance.
type Instance struct {
	path  string
	ctx   *cue.Context
	cfg   *load.Config
	value cue.Value
}

// HolosValue returns the value of the holos field of the exported CUE instance.
func (i *Instance) HolosValue() (v cue.Value, err error) {
	v = i.value.LookupPath(cue.ParsePath("holos"))
	if err = v.Err(); err != nil {
		if strings.HasPrefix(err.Error(), "field not found") {
			slog.Warn(fmt.Sprintf("%s: deprecated usage: nest output under holos: %s", err, i.path), "err", err)
			// Return the deprecated value at the root
			return i.value, nil
		}
		err = errors.Wrap(err)
	}
	return
}

// Discriminate calls the discriminate func for side effects.  Useful to switch
// over the instance kind and apiVersion.
func (i *Instance) Discriminate(discriminate func(tm holos.TypeMeta) error) error {
	v, err := i.HolosValue()
	if err != nil {
		return errors.Wrap(err)
	}

	var tm holos.TypeMeta

	kind := v.LookupPath(cue.ParsePath("kind"))
	if err := kind.Err(); err != nil {
		return errors.Wrap(err)
	}
	if tm.Kind, err = kind.String(); err != nil {
		return errors.Wrap(err)
	}

	version := v.LookupPath(cue.ParsePath("apiVersion"))
	if err := version.Err(); err != nil {
		return errors.Wrap(err)
	}
	if tm.APIVersion, err = version.String(); err != nil {
		return errors.Wrap(err)
	}

	if err := discriminate(tm); err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (i *Instance) Decoder() (*json.Decoder, error) {
	v, err := i.HolosValue()
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

func (i *Instance) Export(enc holos.Encoder) error {
	v, err := i.HolosValue()
	if err != nil {
		return errors.Wrap(err)
	}

	var data interface{}
	if err := v.Decode(&data); err != nil {
		return errors.Wrap(err)
	}

	if err := enc.Encode(&data); err != nil {
		return errors.Wrap(err)
	}

	return nil
}