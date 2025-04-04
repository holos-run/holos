package v1alpha6_test

import (
	"fmt"
	"testing"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const apiVersion string = "v1alpha6"

func TestComponents(t *testing.T) {
	tempDir := testutil.SetupPlatform(t, apiVersion)
	h := testutil.NewComponentHarness(t, tempDir, apiVersion)
	root := h.Root()
	assert.NotEmpty(t, root)

	t.Run("WithNoArtifacts", func(t *testing.T) {
		leaf := "components/minimal"
		c := h.Component(leaf)
		msg := fmt.Sprintf("Expected %s with no artifacts to work, but do nothing", leaf)

		tm, err := c.TypeMeta()
		require.NoError(t, err, msg)

		t.Run("TypeMeta", func(t *testing.T) {
			assert.Equal(t, apiVersion, tm.APIVersion, msg)
			assert.Equal(t, "BuildPlan", tm.Kind, msg)
		})

		t.Run("BuildPlan", func(t *testing.T) {
			bp, err := c.BuildPlan(tm, holos.NewBuildOpts(root, leaf, "deploy", t.TempDir()))
			require.NoError(t, err, msg)
			err = bp.Build(h.Ctx())
			require.NoError(t, err, msg)
		})
	})

	t.Run("BuildPlan", func(t *testing.T) {
		t.Run("Command", func(t *testing.T) {
			t.Run("Generator", func(t *testing.T) {
				leaf := "components/commands/generator/simple"
				c := h.Component(leaf)
				msg := fmt.Sprintf("Expected %s with command generator to render config manifests", leaf)
				tm, err := c.TypeMeta()
				require.NoError(t, err, msg)
				assert.Equal(t, tm.APIVersion, apiVersion)

				t.Run("Build", func(t *testing.T) {
					bp, err := c.BuildPlan(tm, holos.NewBuildOpts(root, leaf, "deploy", t.TempDir()))
					require.NoError(t, err, msg)
					err = bp.Build(h.Ctx())
					require.NoError(t, err, msg)
					// TODO: Check the rendered manifests.
				})
			})
			t.Run("Transformer", func(t *testing.T) {})
			t.Run("Validator", func(t *testing.T) {})
		})
	})
}
