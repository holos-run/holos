package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
func FindCueMod(path string) (root string, err error) {
	origPath := path
	if path, err = filepath.Abs(path); err != nil {
		return "", errors.Wrap(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(path, "cue.mod")); err == nil {
			if root != "" && root != path {
				return "", fmt.Errorf("multiple modules not supported: %v is not %v", root, path)
			}
			root = path
			break
		} else if !os.IsNotExist(err) {
			return "", errors.Wrap(err)
		}
		parent := filepath.Dir(path)
		if parent == path {
			return "", fmt.Errorf("no cue.mod from root to leaf: %v", origPath)
		}
		path = parent
	}
	return root, nil
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
	leaf = DotSlash(leaf)
	return
}

// DotSlash ensures a relative path has a leading ./ needed for CUE loading.
func DotSlash(path string) string {
	clean := filepath.Clean(path)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, ".") || strings.HasPrefix(clean, string(filepath.Separator)) {
		return clean
	}
	return "." + string(filepath.Separator) + clean
}
