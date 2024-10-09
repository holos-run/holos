package artifact

import (
	"sync"

	"github.com/holos-run/holos"
	"github.com/holos-run/holos/internal/errors"
)

func New() *Artifact {
	return &Artifact{m: make(map[holos.FilePath][]byte)}
}

// Artifact represents the fully rendered manifests build from the holos
// rendering pipeline.  Files are organized by keys representing paths relative
// to the current working directory.  Values represent the file content.
type Artifact struct {
	mu sync.RWMutex
	m  map[holos.FilePath][]byte
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
