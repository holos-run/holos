package artifact

import (
	"fmt"
	"os"
	"path/filepath"
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

// Save writes a file to the filesystem.
func (a *Artifact) Save(dir, path holos.FilePath) error {
	fullPath := filepath.Join(string(dir), string(path))
	msg := fmt.Sprintf("could not save %s", fullPath)
	data, ok := a.Get(path)
	if !ok {
		return errors.Format("%s: could not get %s: not set", msg, path)
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	if err := os.WriteFile(fullPath, data, 0666); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	return nil
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
