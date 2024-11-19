package holos

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/holos-run/holos/internal/artifact"
	"github.com/holos-run/holos/internal/errors"
	"gopkg.in/yaml.v3"
)

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

// TagMap represents a map of key values for CUE TagMap for flag parsing.
type TagMap map[string]string

func (t TagMap) Tags() []string {
	parts := make([]string, 0, len(t))
	for k, v := range t {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
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
	if len(parts) != 2 {
		return errors.Format("invalid format, must be tag=value")
	}
	t[parts[0]] = parts[1]
	return nil
}

func (t TagMap) Type() string {
	return "tags"
}

type feature string

const BuildFeature = feature("BUILD")
const ServerFeature = feature("SERVER")
const ClientFeature = feature("CLIENT")
const PreflightFeature = feature("PREFLIGHT")
const GenerateComponentFeature = feature("GENERATE_COMPONENT")
const SecretsFeature = feature("SECRETS")

// Flagger is the interface to check if an experimental feature is enabled.
type Flagger interface {
	Flag(name feature) bool
}

type EnvFlagger struct{}

// Flag returns true if feature name is enabled.
func (e *EnvFlagger) Flag(name feature) bool {
	return os.Getenv(fmt.Sprintf("HOLOS_FEATURE_%s", name)) != ""
}

type Labels map[string]string

type Selector struct {
	Positive map[string]string
	Negative map[string]string
}

// IsSelected returns true when the selector selects the given labels.  An empty
// selector selects.
func (s *Selector) IsSelected(labels Labels) bool {
	// Reject if any positive match is negative.
	for k, v := range s.Positive {
		if val, ok := labels[k]; ok {
			if v != val {
				// Reject - value does not match.
				return false
			} // Select - key and value match.
		} else {
			// Reject - key not present.
			return false
		}
	}

	// Reject if any negative match is positive.
	for k, v := range s.Negative {
		if val, ok := labels[k]; ok {
			if v == val {
				// Reject - value matches, negated.
				return false
			} // Select - value does not match, negated.
		}
	}

	// Select if all checks pass.  An empty selector selects.
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
	msg := "invalid value: %s, must be label=val label==val or label!=val"
	for _, str := range strings.Split(value, ",") {
		splits := map[string]map[string]string{
			"!=": s.Negative,
			"==": s.Positive, // must be before =
			"=":  s.Positive, // must be after ==
		}

		var ok bool
		for sep, m := range splits {
			elems := strings.SplitN(str, sep, 2)
			if len(elems) == 2 {
				k, v := elems[0], elems[1]
				if _, exists := m[k]; exists {
					return errors.Format("already set: %s", k)
				}
				m[k] = v
				ok = true
				break
			}
		}
		if !ok {
			return errors.Format(msg, str)
		}
	}
	return nil
}

// TypeMeta represents the kind and version of a resource holos needs to
// process.  Useful to discriminate generated resources.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

// NewEncoder returns a yaml or json encoder that writes to w.
func NewEncoder(format string, w io.Writer) (Encoder, error) {
	switch format {
	case "yaml":
		encoder := &yamlEncoder{}
		encoder.enc = yaml.NewEncoder(w)
		encoder.enc.SetIndent(2)
		return encoder, nil
	case "json":
		encoder := &jsonEncoder{}
		encoder.enc = json.NewEncoder(w)
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

// IsSelected returns true if all selectors select the given labels or no
// selectors are given.
func IsSelected(labels Labels, selectors ...Selector) bool {
	for _, selector := range selectors {
		if !selector.IsSelected(labels) {
			return false
		}
	}
	return true
}

// BuildOpts represents options common across BuildPlan api versions.  Use
// [NewBuildOpts] to create a new concrete value.
type BuildOpts struct {
	Store       artifact.Store
	Concurrency int
	Stderr      io.Writer
	WriteTo     string
	Path        string
}

func NewBuildOpts(path string) BuildOpts {
	return BuildOpts{
		Store:       artifact.NewStore(),
		Concurrency: min(runtime.NumCPU(), 8),
		Stderr:      os.Stderr,
		WriteTo:     "deploy",
		Path:        path,
	}
}
