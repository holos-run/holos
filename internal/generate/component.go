package generate

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"

	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
)

//go:embed all:components
var components embed.FS

// componentsRoot is the root path to copy component cue code from.
const componentsRoot = "components"

func NewSchematic(root string, name string) (*Schematic, error) {
	data, err := components.ReadFile(filepath.Join(root, name, "schematic.json"))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	schematic := Schematic{Name: name}
	if err := json.Unmarshal(data, &schematic); err != nil {
		return nil, errors.Wrap(err)
	}
	return &schematic, nil
}

// Schematic represents the flags and command metadata stored in the
// schematic.yaml file along side each schematic.
type Schematic struct {
	// Name represents the name of the resource the schematic generates.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Short string `json:"short,omitempty" yaml:"short,omitempty"`
	Long  string `json:"long,omitempty" yaml:"long,omitempty"`

	Chart     *string `json:"chart,omitempty" yaml:"chart,omitempty"`
	Version   *string `json:"version,omitempty" yaml:"version,omitempty"`
	Namespace *string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	RepoName *string `json:"reponame,omitempty" yaml:"reponame,omitempty"`
	RepoURL  *string `json:"repourl,omitempty" yaml:"repourl,omitempty"`

	flagSet *flag.FlagSet
}

func (s *Schematic) FlagSet() *flag.FlagSet {
	if s == nil {
		return nil
	}
	if s.flagSet != nil {
		return s.flagSet
	}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&s.Name, "name", s.Name, "component name")
	if s.Chart != nil {
		fs.StringVar(s.Chart, "chart", *s.Chart, "chart name")
	}
	if s.Version != nil {
		fs.StringVar(s.Version, "component-version", *s.Version, "component version")
	}
	if s.Namespace != nil {
		fs.StringVar(s.Namespace, "namespace", *s.Namespace, "namespace")
	}
	if s.RepoName != nil {
		fs.StringVar(s.RepoName, "repo-name", *s.RepoName, "chart repository name")
	}
	if s.RepoURL != nil {
		fs.StringVar(s.RepoURL, "repo-url", *s.RepoURL, "chart repository url")
	}
	s.flagSet = fs
	return fs
}

// Components returns a slice of embedded component schematics or nil if there
// are none.
func Components(name string) []string {
	entries, err := fs.ReadDir(components, filepath.Join(componentsRoot, name))
	if err != nil {
		return nil
	}
	dirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		dirs = append(dirs, entry.Name())
	}
	return dirs
}

// makeRenderFunc makes a template rendering function for embedded files.
func makeRenderFunc[T any](log *slog.Logger, path string, cfg T) func([]byte) *bytes.Buffer {
	return func(content []byte) *bytes.Buffer {
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			log.Error("could not load template", "err", err)
			return bytes.NewBuffer(content)
		}

		var rendered bytes.Buffer
		if err := tmpl.Execute(&rendered, cfg); err != nil {
			log.Error("could not execute template", "err", err)
			return bytes.NewBuffer(content)
		}

		return &rendered
	}
}

// GenerateComponent writes the cue code for a component to the local working
// directory.
func GenerateComponent(ctx context.Context, kind string, name string, cfg *Schematic) error {
	// use name from args to build the source path
	path := filepath.Join(componentsRoot, kind, name)
	// write to the current directory.
	dstPath := filepath.Join(getCwd(ctx))
	log := logger.FromContext(ctx).With("name", cfg.Name, "path", dstPath)
	log.DebugContext(ctx, "mkdir")
	if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
		return errors.Wrap(err)
	}

	mapper := makeRenderFunc(log, path, cfg)
	if err := copyEmbedFS(ctx, components, path, dstPath, mapper); err != nil {
		return errors.Wrap(err)
	}

	log.InfoContext(ctx, "generated component")
	return nil
}
