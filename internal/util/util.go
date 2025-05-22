package util

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
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

// FindRootLeaf returns the root path containing the cue.mod and the leaf path
// relative to the root for the given target path.  FindRootLeaf calls
// [filepath.Clean] on the returned paths.
func FindRootLeaf(target string) (root string, leaf string, err error) {
	if root, err = FindCueMod(target); err != nil {
		return "", "", err
	}
	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", "", err
	}
	if leaf, err = filepath.Rel(root, absPath); err != nil {
		return "", "", err
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

// MakeCopyFunc returns a [fs.WalkDirFunc] copying from embed.FS efs to dest.
// Useful for embedded test case fixtures.
func MakeCopyFunc(ctx context.Context, efs embed.FS, dest string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err != nil {
				return err
			}
			if path == "." {
				return nil
			}
			fullPath := filepath.Join(dest, path)

			switch {
			case d.IsDir():
				if err := os.MkdirAll(fullPath, 0o777); err != nil {
					return err
				}
			default:
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o777); err != nil {
					return err
				}
				srcFile, err := efs.Open(path)
				if err != nil {
					return err
				}
				defer srcFile.Close()

				dstFile, err := os.Create(fullPath)
				if err != nil {
					return err
				}
				defer dstFile.Close()

				if _, err := io.Copy(dstFile, srcFile); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// Executable returns the fully qualified path to the "holos" executable by
// searching through the PATH. If the current executable starts with "holos",
// it returns that path, otherwise it searches for "holos" in PATH.
func Executable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(filepath.Base(exe), "holos") {
		return exe, nil
	}
	return exec.LookPath("holos")
}
