package render

import (
	"context"
	"sync"

	"github.com/holos-run/holos"
)

func NewArtifact() *Artifact {
	return &Artifact{
		m: make(map[holos.FilePath]holos.FileContent),
	}
}

// Artifact represents the fully rendered manifests build from the holos
// rendering pipeline.  Files are organized by keys representing paths relative
// to the current working directory.  Values represent the file string content.
type Artifact struct {
	mu sync.RWMutex
	m  map[holos.FilePath]holos.FileContent
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

type Setter interface {
	Set(path holos.FilePath, content holos.FileContent)
}

type Getter interface {
	Get(path holos.FilePath) (content holos.FileContent, ok bool)
	Keys() []holos.FilePath
}

// Builder represents the primary interface to the holos rendering pipeline.
// Pipeline stages render manifests then use [Setter] s to associate the
// rendered output with an output [holos.FilePath].
type Builder interface {
	Build(ctx context.Context, s Setter) error
}
