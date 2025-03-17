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
	"github.com/holos-run/holos/internal/generate"
	"github.com/stretchr/testify/assert"
)

//go:embed all:platform
var f embed.FS

// must align with embed all:platform directory
const platform string = "platform"

func TestComponentAlpha6(t *testing.T) {
	h := newHarness(t, "components/v1alpha6")
	err := h.c.Render(h.ctx)
	assert.NoError(t, err)
}

type harness struct {
	c   *component.Component
	ctx context.Context
}

func newHarness(t testing.TB, leaf string) *harness {
	t.Helper()

	// temp directory for the test, contains the platform and components for the
	// test cases.
	tempDir := t.TempDir()
	root := filepath.Join(tempDir, platform)
	if err := os.MkdirAll(root, 0o777); err != nil {
		t.Fatalf("could not make platform directory: %v", err)
	}

	// test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Initialize the platform
	if err := generate.GeneratePlatform(ctx, root, "v1alpha6"); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}

	// Copy the components for the test cases
	if err := fs.WalkDir(f, platform, walkDirFunc(ctx, tempDir)); err != nil {
		t.Fatalf("could not prepare test directory: %v", err)
	}

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
				if err := os.MkdirAll(fullPath, 0o777); err != nil {
					return err
				}
			default:
				if err := os.MkdirAll(filepath.Dir(fullPath), 0o777); err != nil {
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
