package util

import (
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/errors"
)

// EnsureNewline adds a trailing newline if not already there.
func EnsureNewline(b []byte) []byte {
	if len(b) > 0 && b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	return b
}

// FindCueMod returns the root module location containing the cue.mod.
func FindCueMod(target string) (dir string, err error) {
	dir, err = filepath.Abs(target)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		if _, err = os.Stat(filepath.Join(dir, "cue.mod")); err == nil {
			break
		} else if !os.IsNotExist(err) {
			return "", errors.Wrap(err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.Format("no cue.mod from root to leaf: %v", target)
		}
		dir = parent
	}

	return
}

func FindRootLeaf(target string) (root string, leaf string, err error) {
	if root, err = FindCueMod(target); err != nil {
		return "", "", errors.Wrap(err)
	}
	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", "", errors.Wrap(err)
	}
	if leaf, err = filepath.Rel(root, absPath); err != nil {
		return "", "", errors.Wrap(err)
	}
	// Needed for CUE to load the path properly.
	leaf = "." + string(os.PathSeparator) + leaf
	return
}
