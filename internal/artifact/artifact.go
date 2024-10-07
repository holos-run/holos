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
		m:       make(map[holos.FilePath]holos.FileContent),
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
	m       map[holos.FilePath]holos.FileContent
	writeTo string
}

// Set sets an artifact file with write locking.
func (a *Artifact) Set(path holos.FilePath, content holos.FileContent) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.m[path] = content
}

// Get gets the content of an artifact file path with read locking.
func (a *Artifact) Get(path holos.FilePath) (content holos.FileContent, ok bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	content, ok = a.m[path]
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

// Save writes artifacts to the output directory relative to the write to path.
func (a *Artifact) Save(ctx context.Context) error {
	log := logger.FromContext(ctx)
	keys := a.Keys()
	for _, key := range keys {
		path := filepath.Join(a.writeTo, string(key))

		content, ok := a.Get(holos.FilePath(path))
		if !ok {
			return errors.Format("missing key: %s", path)
		}

		dir := filepath.Dir(string(path))
		if err := os.MkdirAll(dir, os.FileMode(0775)); err != nil {
			log.WarnContext(ctx, "could not mkdir", "path", dir, "err", err)
			return errors.Wrap(err)
		}

		if err := os.WriteFile(path, []byte(content), os.FileMode(0666)); err != nil {
			log.WarnContext(ctx, "could not save", "path", path, "err", err)
			return errors.Format("could not save: %w", err)
		}
		log.DebugContext(ctx, "wrote "+path, "action", "write", "path", path, "status", "ok")
	}
	return nil
}
