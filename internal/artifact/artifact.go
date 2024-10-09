package artifact

import (
	"sync"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/errors"
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
