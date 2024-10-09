package artifact

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/server/middleware/logger"
)

type Option func(*Artifact)

// WriteTo configures the directory to write artifacts into.
func WriteTo(path string) Option {
	return func(a *Artifact) {
		a.writeTo = path
	}
}

func New(options ...Option) *Artifact {
	a := &Artifact{
		m:       make(map[holos.FilePath][]byte),
		writeTo: "deploy",
	}
	for _, o := range options {
		o(a)
	}
	return a
}

// Artifact represents the fully rendered manifests build from the holos
// rendering pipeline.  Files are organized by keys representing paths relative
// to the current working directory.  Values represent the file string content.
type Artifact struct {
	mu      sync.RWMutex
	m       map[holos.FilePath][]byte
	writeTo string
}

// Set sets an artifact file with write locking.  Set returns an error if the
// artifact was previously set.
func (a *Artifact) Set(path holos.FilePath, data []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.m[path]; ok {
		return errors.Format("could not set artifact %s: already set", path)
	}
	a.m[path] = data
	return nil
}

// Get gets the content of an artifact with read locking.
func (a *Artifact) Get(path holos.FilePath) (data []byte, ok bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	data, ok = a.m[path]
	return
}

func (a *Artifact) Keys() []holos.FilePath {
	a.mu.RLock()
	defer a.mu.RUnlock()
	keys := make([]holos.FilePath, 0, len(a.m))
	for key := range a.m {
		keys = append(keys, key)
	}
	return keys
}

// Save saves one artifact to the current working directory, creating parent
// directories as needed.
func (a *Artifact) Save(ctx context.Context, path holos.FilePath) error {
	log := logger.FromContext(ctx)
	data, ok := a.Get(path)
	if !ok {
		return errors.Format("missing key: %s", path)
	}

	dir := filepath.Dir(string(path))
	if err := os.MkdirAll(dir, os.FileMode(0775)); err != nil {
		log.WarnContext(ctx, "could not mkdir: "+err.Error(), "path", dir, "err", err)
		return errors.Wrap(err)
	}

	if err := os.WriteFile(string(path), data, os.FileMode(0666)); err != nil {
		log.WarnContext(ctx, "could not save: "+err.Error(), "path", path, "err", err)
		return errors.Format("could not save: %w", err)
	}
	log.DebugContext(ctx, "wrote: "+string(path), "action", "write", "path", path, "status", "ok")
	return nil
}
