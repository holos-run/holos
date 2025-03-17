package component_test

// "github.com/stretchr/testify/require"
import (
	"context"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/cli/render/component"
	"github.com/stretchr/testify/assert"
)

//go:embed all:platform
var f embed.FS

func TestComponentAlpha6(t *testing.T) {
	h := newHarness(t, "components/v1alpha6")
	assert.NoError(t, h.c.Render(h.ctx))
}

type harness struct {
	c   *component.Component
	ctx context.Context
}

func newHarness(t testing.TB, leaf string) *harness {
	t.Helper()
	tempDir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	if err := fs.WalkDir(f, ".", walkDirFunc(ctx, tempDir)); err != nil {
		t.Fatalf("could not prepare test directory: %v", err)
	}
	root := filepath.Join(tempDir, "platform")
	return &harness{
		c:   component.New(root, leaf, component.NewConfig()),
		ctx: ctx,
	}
}

func walkDirFunc(ctx context.Context, tempDir string) fs.WalkDirFunc {
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
				if err := os.MkdirAll(fullPath, 0777); err != nil {
					return err
				}
			default:
				if err := os.MkdirAll(filepath.Dir(fullPath), 0777); err != nil {
					return err
				}
				srcFile, err := f.Open(path)
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
