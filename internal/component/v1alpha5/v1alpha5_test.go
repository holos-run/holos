package v1alpha5_test

// "github.com/stretchr/testify/require"
import (
	"context"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/component"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/generate"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

//go:embed all:platform
var f embed.FS

// must align with embed all:platform directory
const platform string = "platform"

func TestComponents(t *testing.T) {
	h := newHarness(t)

	t.Run("TypeMeta", func(t *testing.T) {
		msg := "Expected a minimal component with only typemeta.yaml to work, but do nothing"
		err := h.component("components/typemeta").Render(h.ctx)
		assert.NoError(t, err, msg)
	})

	t.Run("BasicDeployment", func(t *testing.T) {
		msg := "Expected a basic cue resources generator to render a Deployment manifest"
		c := h.component("components/basic")
		err := c.Render(h.ctx)
		assert.NoError(t, err, msg)

		// Verify the rendered artifacts.
		actual, err := h.load("deploy/components/basic/resources.gen.yaml")
		assert.NoError(t, err, msg)
		expected, err := h.load("components/basic/expected.yaml")
		assert.NoError(t, err, msg)

		// Compare in both directions, necessary to catch missing fields.
		assert.EqualValues(t, expected, actual, msg)
		assert.EqualValues(t, actual, expected, msg)
	})
}

type harness struct {
	root string
	ctx  context.Context
}

func (h *harness) component(path string) *component.Component {
	return component.New(h.root, path, component.NewConfig())
}

func (h *harness) load(path string) (any, error) {
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

func newHarness(t testing.TB) *harness {
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
	if err := fs.WalkDir(f, platform, makeCopyFunc(ctx, tempDir)); err != nil {
		t.Fatalf("could not prepare test directory: %v", err)
	}

	return &harness{
		root: root,
		ctx:  ctx,
	}
}

func makeCopyFunc(ctx context.Context, tempDir string) fs.WalkDirFunc {
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
