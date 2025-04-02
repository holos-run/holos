package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/interpreter/embed"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/util"
)

// cue context and loading is not safe for concurrent use.
var cueMutex sync.Mutex

// ExtractYAML extracts yaml encoded data from file paths.  The data is unified
// into one [cue.Value].  If a path element is a directory, all files in the
// directory are loaded non-recursively.
//
// Attribution: https://github.com/cue-lang/cue/issues/3504
// Deprecated: Use cue embed instead.
func ExtractYAML(ctxt *cue.Context, filepaths []string) (cue.Value, error) {
	value := ctxt.CompileString("")
	files := make([]string, 0, 10*len(filepaths))

	for _, path := range filepaths {
		info, err := os.Stat(path)
		if err != nil {
			return value, errors.Wrap(err)
		}

		if !info.IsDir() {
			files = append(files, path)
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return value, errors.Wrap(err)
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}

	for _, file := range files {
		f, err := yaml.Extract(file, nil)
		if err != nil {
			return value, errors.Wrap(err)
		}
		value = value.Unify(ctxt.BuildFile(f))
	}

	return value, nil
}

// BuildInstance builds the cue configuration instance at leaf relative to the
// root cue module.
func BuildInstance(root, leaf string, tags []string) (*Instance, error) {
	cueMutex.Lock()
	defer cueMutex.Unlock()
	leaf = util.DotSlash(leaf)

	cfg := &load.Config{
		Dir:        root,
		ModuleRoot: root,
		Tags:       tags,
	}
	ctxt := cuecontext.New(cuecontext.Interpreter(embed.New()))

	bis := load.Instances([]string{leaf}, cfg)
	values, err := ctxt.BuildInstances(bis)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	inst := &Instance{
		path:  leaf,
		ctx:   ctxt,
		cfg:   cfg,
		value: values[0],
	}

	return inst, nil
}

// LoadInstance loads the cue configuration instance at path.  External data
// file paths are loaded by calling [ExtractYAML] providing filepaths.  The
// extracted data values are unified with the platform configuration [cue.Value]
// in the returned [Instance].
//
// Deprecated: use BuildInstance instead.
func LoadInstance(path string, filepaths []string, tags []string) (*Instance, error) {
	cueMutex.Lock()
	defer cueMutex.Unlock()
	root, leaf, err := util.FindRootLeaf(path)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	cfg := &load.Config{
		Dir:        root,
		ModuleRoot: root,
		Tags:       tags,
	}
	ctxt := cuecontext.New(cuecontext.Interpreter(embed.New()))

	bis := load.Instances([]string{path}, cfg)
	values, err := ctxt.BuildInstances(bis)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	value, err := ExtractYAML(ctxt, filepaths)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	value = value.Unify(values[0])

	inst := &Instance{
		path:  leaf,
		ctx:   ctxt,
		cfg:   cfg,
		value: value,
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
//
// Deprecated: use loadTypeMeta() to discriminate without cue instead.
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
