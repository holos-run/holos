package holos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/holos-run/holos/internal/artifact"
	"github.com/holos-run/holos/internal/errors"
	"gopkg.in/yaml.v3"
)

// Interface implementation assertions.
var _ Encoder = &yamlEncoder{}
var _ Encoder = &jsonEncoder{}
var _ OrderedEncoder = &orderedEncoder{}

// StringSlice represents zero or more flag values.
type StringSlice []string

// String implements the flag.Value interface.
func (i *StringSlice) String() string {
	return fmt.Sprint(*i)
}

// Type implements the pflag.Value interface and describes the type.
func (i *StringSlice) Type() string {
	return "strings"
}

// Set implements the flag.Value interface.
func (i *StringSlice) Set(value string) error {
	for _, str := range strings.Split(value, ",") {
		*i = append(*i, str)
	}
	return nil
}

// TagMap represents a map of key values for CUE TagMap for flag parsing.  The
// values are pointers to disambiguate between the case where a tag is a boolean
// ("--inject foo") and the case where a tag has a string zero value ("--inject
// foo=").  Refer to the Tags field of [cue/load.Config]
//
// [cue/load.Config]: https://pkg.go.dev/cuelang.org/go@v0.10.1/cue/load#Config
type TagMap map[string]*string

const TagMapHelp = "set the value of a cue @tag field in the form key=value or simply key"

func (t TagMap) Tags() []string {
	parts := make([]string, 0, len(t))
	for tag, val := range t {
		if val == nil {
			parts = append(parts, tag)
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", tag, *val))
		}
	}
	return parts
}

func (t TagMap) String() string {
	return strings.Join(t.Tags(), " ")
}

// Set sets a value.  Only one value per flag is supported.  For example
// --inject=foo=bar --inject=bar=baz.  For JSON values, --inject=foo=bar,bar=baz
// is not supported.
func (t TagMap) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	switch len(parts) {
	case 1:
		t[parts[0]] = nil
	case 2:
		t[parts[0]] = &parts[1]
	default:
		return errors.Format("invalid format, must be tag=value")
	}
	return nil
}

func (t TagMap) Type() string {
	return "tags"
}

type Labels map[string]string

type Selectors []Selector

// String implements the flag.Value interface.
func (s *Selectors) String() string {
	return fmt.Sprint(*s)
}

// Type implements the pflag.Value interface and describes the type.
func (s *Selectors) Type() string {
	return "selectors"
}

// Set implements the flag.Value interface.
func (s *Selectors) Set(value string) error {
	selector := Selector{}
	if err := selector.Set(value); err != nil {
		return err
	}
	*s = append(*s, selector)
	return nil
}

type Selector struct {
	Positive map[string]string
	Negative map[string]string
}

// IsSelected returns true when the selector selects the given labels
func (s *Selector) IsSelected(labels Labels) bool {
	if s == nil {
		return true // Nil selector selects everything
	}

	// Check positive matches
	for k, v := range s.Positive {
		if val, ok := labels[k]; !ok || v != val {
			return false
		}
	}

	// Check negative matches
	for k, v := range s.Negative {
		if val, ok := labels[k]; ok && v == val {
			return false
		}
	}

	return true
}

func (s *Selector) String() string {
	elems := make([]string, 0, len(s.Positive)+len(s.Negative))
	for k, v := range s.Positive {
		elems = append(elems, fmt.Sprintf("%s==%s", k, v))
	}
	for k, v := range s.Negative {
		elems = append(elems, fmt.Sprintf("%s!=%s", k, v))
	}
	return strings.Join(elems, ",")
}

func (s *Selector) Type() string {
	return "selector"
}

func (s *Selector) Set(value string) error {
	if s.Positive == nil {
		s.Positive = map[string]string{}
	}
	if s.Negative == nil {
		s.Negative = map[string]string{}
	}

	for _, str := range strings.Split(value, ",") {
		if strings.Contains(str, "!=") {
			elems := strings.SplitN(str, "!=", 2)
			s.Negative[elems[0]] = elems[1]
			continue
		}

		// Treat both = and == as positive matches
		str = strings.ReplaceAll(str, "==", "=")
		elems := strings.SplitN(str, "=", 2)
		if len(elems) != 2 {
			return errors.Format("invalid value: %s", str)
		}
		s.Positive[elems[0]] = elems[1]
	}

	return nil
}

// TypeMeta represents the kind and version of a resource holos needs to
// process.  Useful to discriminate generated resources.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

// NewSequentialEncoder returns a yaml or json encoder that writes to w.  The
// encoding argument may be "json" or "yaml".
func NewSequentialEncoder(encoding string, w io.Writer) (OrderedEncoder, error) {
	enc, err := NewEncoder(encoding, w)
	if err != nil {
		return nil, err
	}
	seqEnc := &orderedEncoder{
		Encoder: enc,
		buf:     make(map[int]any),
	}
	return seqEnc, nil
}

// NewEncoder returns a yaml or json encoder that writes to w.  The format
// argument specifies "yaml" or "json" format output.
func NewEncoder(format string, w io.Writer) (Encoder, error) {
	switch format {
	case "yaml":
		encoder := &yamlEncoder{
			enc: yaml.NewEncoder(w),
		}
		encoder.enc.SetIndent(2)
		return encoder, nil
	case "json":
		encoder := &jsonEncoder{
			enc: json.NewEncoder(w),
		}
		encoder.enc.SetIndent("", "  ")
		return encoder, nil
	default:
		return nil, errors.Format("invalid format: %s, must be yaml or json", format)
	}
}

type jsonEncoder struct {
	mu  sync.Mutex
	enc *json.Encoder
}

func (j *jsonEncoder) Encode(v any) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return errors.Wrap(j.enc.Encode(v))
}

func (j *jsonEncoder) Close() error {
	return nil
}

type yamlEncoder struct {
	mu  sync.Mutex
	enc *yaml.Encoder
}

func (y *yamlEncoder) Encode(v any) error {
	y.mu.Lock()
	defer y.mu.Unlock()
	return errors.Wrap(y.enc.Encode(v))
}

func (y *yamlEncoder) Close() error {
	return errors.Wrap(y.enc.Close())
}

// IsSelected returns true if any one selector selects the given labels or no
// selectors are given.
func IsSelected(labels Labels, selectors ...Selector) bool {
	if len(selectors) == 0 {
		return true
	}
	for _, selector := range selectors {
		if selector.IsSelected(labels) {
			return true
		}
	}
	return false
}

type orderedEncoder struct {
	Encoder
	mu   sync.Mutex
	buf  map[int]any
	next int
}

// Encode encodes v in sequential or starting with idx 0.
//
// It is an error to provide idx values less than the next to encode.  Is is an
// error to provide the same idx value more than once.
func (s *orderedEncoder) Encode(idx int, v any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if idx < s.next {
		return fmt.Errorf("could not encode idx %d: must be greater than or equal to next idx %d", idx, s.next)
	}

	// If this is the next expected index, encode it and any buffered values
	if idx == s.next {
		if err := s.Encoder.Encode(v); err != nil {
			return errors.Wrap(err)
		}
		s.next++

		// Encode any buffered values that come next in sequence
		for {
			if v, ok := s.buf[s.next]; ok {
				if err := s.Encoder.Encode(v); err != nil {
					return errors.Wrap(err)
				}
				delete(s.buf, s.next)
				s.next++
			} else {
				break
			}
		}
		return nil
	}

	if _, ok := s.buf[idx]; ok {
		return fmt.Errorf("could not encode idx %d: already exists", idx)
	}

	// Buffer out-of-order value
	s.buf[idx] = v
	return nil
}

// TODO(jjm): consider moving the BuildOpts struct to the component package.

// BuildOpts represents options common across BuildPlan api versions.  Use
// [NewBuildOpts] to create a new concrete value.
type BuildOpts struct {
	Store       artifact.Store
	Concurrency int
	Stderr      io.Writer
	WriteTo     string
	// Path represents the component path relative to the platform module root.
	Path string
	// Tags represents user managed tags including a component name, labels, and
	// annotations.
	Tags []string

	root    string
	leaf    string
	writeTo string
	tempDir string
}

// NewBuildOpts returns a [BuildOpts] configured to build the component at leaf
// from the platform module at root writing rendered manifests into the deploy
// directory.
func NewBuildOpts(root, leaf, writeTo, tempDir string) BuildOpts {
	return BuildOpts{
		Store:       artifact.NewStore(),
		Concurrency: min(runtime.NumCPU(), 8),
		Stderr:      os.Stderr,
		Tags:        make([]string, 0, 10),

		root:    filepath.Clean(root),
		leaf:    filepath.Clean(leaf),
		writeTo: filepath.Clean(writeTo),
		tempDir: filepath.Clean(tempDir),
	}
}

// Root returns the platform root directory containing the cue module.
func (b *BuildOpts) Root() string {
	return b.root
}

// Leaf returns the cleaned component path relative to the platform root. For
// example "components/podinfo"
func (b *BuildOpts) Leaf() string {
	return b.leaf
}

// AbsLeaf returns the absolute cleaned component path.
func (b *BuildOpts) AbsLeaf() string {
	return filepath.Join(b.root, b.leaf)
}

// AbsWriteTo returns the absolute path to the write to directory, usually the
// deploy sub directory of the platform module root.
func (b *BuildOpts) AbsWriteTo() string {
	return filepath.Join(b.root, b.writeTo)
}

// TempDir returns the temporary directory managed by holos and injected into
// cue using a [BuildContext] so artifacts can refer to the same path in the
// configuration.
func (b *BuildOpts) TempDir() string {
	return b.tempDir
}

// BuildContext represents build context values provided by the holos render
// component command.  These values are expected to be randomly generated and
// late binding, meaning they cannot be known ahead of time in a static
// configuration.  As such, CUE configuration may refer to the values here which
// will be populated by holos when the final build plan is exported from CUE.
type BuildContext struct {
	// TempDir represents the temporary directory managed and owned by the holos
	// render component command for the execution of one BuildPlan.  Multiple
	// tasks in the build plan share this temporary directory and therefore should
	// avoid reading and writing into the same sub-directories as one another.
	TempDir string `json:"tempDir" yaml:"tempDir"`
}
