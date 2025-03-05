package artifact

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/holos-run/holos/internal/errors"
)

// NewStore should provide a concrete Store.
var _ Store = NewStore()

// TODO(jjm): Refactor / Remove the Artifact Store.  The memory store is almost
// certainly unnecessary now that we're using a single temp directory for the
// whole build plan context as of v1alpha6.  The Artifact store was originally
// created to persist intermediate artifacts across per-task temp directories
// which we no longer have as of v1alpha6, so surely a chesterton fence we can
// tear down.

// Store sets and gets data for file artifacts.
//
// Concrete values must ensure Set is write once, returning an error if a given
// FilePath was previously Set.  Concrete values must be safe for concurrent
// reads and writes.  Use [NewStore] to create a new concrete value.
type Store interface {
	Get(path string) (data []byte, ok bool)
	Set(path string, data []byte) error
	// Save a file or directory from the store to the filesystem.
	Save(dir, path string) error
	// Load a file or directory from the filesystem into the store.
	Load(dir, path string) error
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
	slog.Debug(fmt.Sprintf("store: set path %s", path), "component", "store", "op", "set", "path", path, "bytes", len(data))
	return nil
}

// Get gets the content of an artifact with read locking.
func (a *MapStore) Get(path string) (data []byte, ok bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	data, ok = a.m[path]
	slog.Debug(fmt.Sprintf("store: get path %s ok %v", path, ok), "component", "store", "op", "get", "path", path, "bytes", len(data), "ok", ok)
	return
}

// Save writes a file or directory tree to the filesystem.
func (a *MapStore) Save(dir, path string) error {
	if strings.HasSuffix(path, "/") {
		return errors.Format("path must not end in a /")
	}

	fullPath := filepath.Join(dir, path)
	msg := fmt.Sprintf("could not save %s", fullPath)

	// Save a single file and return.
	if data, ok := a.Get(path); ok {
		if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		if err := os.WriteFile(fullPath, data, 0666); err != nil {
			return errors.Format("%s: %w", msg, err)
		}
		return nil
	}

	// Assume path is a directory, find all prefix matches.
	keys := a.Keys()
	prefix := fmt.Sprintf("%s/", path)
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			data, _ := a.Get(key)
			fullPath = filepath.Join(dir, key)
			if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
				return errors.Format("%s: %w", msg, err)
			}
			if err := os.WriteFile(fullPath, data, 0666); err != nil {
				return errors.Format("%s: %w", msg, err)
			}
		}
	}

	return nil
}

// Load saves a file or directory tree to the store.
func (a *MapStore) Load(dir, path string) error {
	fileSystem := os.DirFS(dir)
	err := fs.WalkDir(fileSystem, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err)
		}
		// Skip over directories.
		if d.IsDir() {
			return nil
		}
		// Load files into the store.
		data, err := fs.ReadFile(fileSystem, path)
		if err != nil {
			return errors.Wrap(err)
		}
		if err := a.Set(path, data); err != nil {
			return errors.Wrap(err)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err)
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
