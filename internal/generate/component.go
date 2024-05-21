package generate

import (
	"bytes"
	"context"
	"embed"
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

// CueConfig represents the config values passed to cue go templates.
type CueConfig struct {
	ComponentName string
	flagSet       *flag.FlagSet
}

func (c *CueConfig) FlagSet() *flag.FlagSet {
	if c == nil {
		return nil
	}
	if c.flagSet != nil {
		return c.flagSet
	}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&c.ComponentName, "name", "example", "component name")
	c.flagSet = fs
	return fs
}

type HelmConfig struct {
	ComponentName string
	flagSet       *flag.FlagSet
}

func (c *HelmConfig) FlagSet(name string) *flag.FlagSet {
	if c == nil {
		return nil
	}
	if c.flagSet != nil {
		return c.flagSet
	}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&c.ComponentName, "name", name, "component name")
	c.flagSet = fs
	return fs
}

// CueComponents returns a slice of embedded component schematics or nil if there are none.
func CueComponents() []string {
	entries, err := fs.ReadDir(components, filepath.Join(componentsRoot, "cue"))
	if err != nil {
		return nil
	}
	dirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		dirs = append(dirs, entry.Name())
	}
	return dirs
}

// HelmComponents returns a slice of embedded component schematics or nil if there are none.
func HelmComponents() []string {
	entries, err := fs.ReadDir(components, filepath.Join(componentsRoot, "helm"))
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

// GenerateCueComponent writes the cue code for a component to the local working
// directory.
func GenerateCueComponent(ctx context.Context, name string, cfg *CueConfig) error {
	path := filepath.Join(componentsRoot, "cue", name)
	dstPath := filepath.Join(getCwd(ctx), cfg.ComponentName)
	log := logger.FromContext(ctx).With("name", cfg.ComponentName, "path", dstPath)
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

// GenerateHelmComponent writes the cue code for a component to the local working
// directory.
func GenerateHelmComponent(ctx context.Context, name string, cfg *HelmConfig) error {
	path := filepath.Join(componentsRoot, "helm", name)
	dstPath := filepath.Join(getCwd(ctx), cfg.ComponentName)
	log := logger.FromContext(ctx).With("name", cfg.ComponentName, "path", dstPath)
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
