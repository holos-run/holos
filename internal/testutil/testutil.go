package testutil

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/holos-run/holos/internal/component"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/generate"
	"gopkg.in/yaml.v3"
)

//go:embed all:fixtures
var Fixtures embed.FS

// fixtures represents the name of the embedded fixtures directory
const fixtures string = "fixtures"

func SetupPlatform(t testing.TB, apiVersion string) (tempDir string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	tempDir = t.TempDir()

	// Initialize the platform
	if err := generate.GeneratePlatform(ctx, tempDir, apiVersion); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}

	// Copy the components for the test cases
	if err := fs.WalkDir(Fixtures, filepath.Join(fixtures, apiVersion), MakeCopyFunc(ctx, Fixtures, tempDir)); err != nil {
		t.Fatalf("could not prepare test directory: %v", err)
	}

	return tempDir
}

type ComponentHarness struct {
	// root represents the temp dir containing the cue.mod
	root string
	// base represents the path relative to root containing fixtures (platform/, components/, etc...)
	base string
	ctx  context.Context
}

func NewComponentHarness(t testing.TB, tempDir, apiVersion string) *ComponentHarness {
	t.Helper()

	// test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	return &ComponentHarness{
		root: tempDir,
		base: filepath.Join(fixtures, apiVersion),
		ctx:  ctx,
	}
}

func (h *ComponentHarness) Root() string {
	return filepath.Clean(h.root)
}

func (h *ComponentHarness) Ctx() context.Context {
	return h.ctx
}

// Base returns the fixture base path relative to the root, for example "fixtures/v1alpha6"
func (h *ComponentHarness) Base() string {
	return filepath.Clean(h.base)
}

// AbsBase returns the absolute path to the api version base.
func (h *ComponentHarness) AbsBase() string {
	return filepath.Join(h.root, h.base)
}

// Component returns a new component with fully qualified path
// filepath.Join(root, base, path).
func (h *ComponentHarness) Component(path string) *component.Component {
	return component.New(
		h.Root(),
		filepath.Join(h.Base(), path),
	)
}

// Load unmarshals the first document in yaml or json path relative to the
// harness root.
func (h *ComponentHarness) Load(path string) (any, error) {
	data, err := os.ReadFile(filepath.Join(h.root, path))
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var result any
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, errors.Wrap(err)
	}
	return result, nil
}

func MakeCopyFunc(ctx context.Context, fsys fs.FS, tempDir string) fs.WalkDirFunc {
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
			fullPath := filepath.Join(tempDir, path)

			switch {
			case d.IsDir():
				if err := os.MkdirAll(fullPath, 0o777); err != nil {
					return err
				}
			default:
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o777); err != nil {
					return err
				}
				srcFile, err := fsys.Open(path)
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

// Capitalize capitalizes the first letter of a UTF-8 string.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	upperRune := unicode.ToUpper(r)
	if r == upperRune {
		return s
	}
	return string(upperRune) + s[size:]
}
