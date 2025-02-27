package artifact

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/holos-run/holos/internal/errors"
)

// NewStore should provide a concrete Store.
var _ Store = NewStore()

// Store sets and gets data for file artifacts.
//
// Concrete values must ensure Set is write once, returning an error if a given
// FilePath was previously Set.  Concrete values must be safe for concurrent
// reads and writes.  Use [NewStore] to create a new concrete value.
type Store interface {
	Get(path string) (data []byte, ok bool)
	Set(path string, data []byte) error
	// Save previously set path to dir preserving directories.
	Save(dir, path string) error
}

func NewStore() *MapStore {
	return &MapStore{m: make(map[string][]byte)}
}

// MapStore represents the fully rendered manifests build from the holos
// rendering pipeline.  Files are organized by keys representing paths relative
// to the current working directory.  Values represent the file content.
type MapStore struct {
	mu sync.RWMutex
	m  map[string][]byte
}

// Set sets an artifact file with write locking.  Set returns an error if the
// artifact was previously set.
func (a *MapStore) Set(path string, data []byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.m[path]; ok {
		return errors.Format("%s already set", path)
	}
	a.m[path] = data
	return nil
}

// Get gets the content of an artifact with read locking.
func (a *MapStore) Get(path string) (data []byte, ok bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	data, ok = a.m[path]
	return
}

// Save writes a file to the filesystem.
func (a *MapStore) Save(dir, path string) error {
	fullPath := filepath.Join(dir, path)
	msg := fmt.Sprintf("could not save %s", fullPath)
	// TODO(jjm): replace slash hack with artifact directory support maybe
	if strings.HasSuffix(path, "/") {
		for _, key := range a.Keys() {
			if strings.HasPrefix(key, path) {
				if err := a.Save(dir, key); err != nil {
					return errors.Wrap(err)
				}
			}
		}
		return nil
	}
	data, ok := a.Get(path)
	if !ok {
		return errors.Format("%s: missing %s", msg, path)
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	if err := os.WriteFile(fullPath, data, 0666); err != nil {
		return errors.Format("%s: %w", msg, err)
	}
	return nil
}

func (a *MapStore) Keys() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	keys := make([]string, 0, len(a.m))
	for key := range a.m {
		keys = append(keys, key)
	}
	return keys
}
