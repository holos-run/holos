package v1alpha6_test

import (
	"fmt"
	"path/filepath"
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
	assert.NotEmpty(t, h.Root())

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
			bp, err := c.BuildPlan(tm, holos.NewBuildOpts(h.Root(), leaf, "deploy", t.TempDir()))
			require.NoError(t, err, msg)
			err = bp.Build(h.Ctx())
			require.NoError(t, err, msg)
		})
	})

	t.Run("BuildPlan", func(t *testing.T) {
		t.Run("Generator", func(t *testing.T) {
			for _, tc := range []string{"simple", "directory"} {
				testComponent(t, h, "generator", tc)
			}
		})

		t.Run("Transformer", func(t *testing.T) {
			for _, tc := range []string{"simple", "kustomize", "join"} {
				testComponent(t, h, "transformer", tc)
			}
		})
		t.Run("Validator", func(t *testing.T) {
			t.Run("Command", func(t *testing.T) {

			})
		})
	})
}

// testComponent runs tests against a component directory fixtures by building
// the component and asserting the rendered config manifests against a
// want_name.gen.yaml file in the component leaf directory.
func testComponent(t *testing.T, h *testutil.ComponentHarness, kind, name string) {
	t.Run(testutil.Capitalize(name), func(t *testing.T) {
		leaf := filepath.Join("components", kind, name)
		c := h.Component(leaf)
		msg := fmt.Sprintf("Expected %s with %s to render config manifests", leaf, kind)
		tm, err := c.TypeMeta()
		require.NoError(t, err, msg)
		assert.Equal(t, tm.APIVersion, apiVersion)

		t.Run("Build", func(t *testing.T) {
			bp, err := c.BuildPlan(tm, holos.NewBuildOpts(h.Root(), leaf, "deploy", t.TempDir()))
			require.NoError(t, err, msg)
			err = bp.Build(h.Ctx())
			require.NoError(t, err, msg)

			// Validate the rendered manifest
			have, err := h.Load(filepath.Join("deploy", "components", name, fmt.Sprintf("%s.gen.yaml", name)))
			require.NoError(t, err, msg)
			want, err := h.Load(filepath.Join(h.Base(), leaf, fmt.Sprintf("want_%s.gen.yaml", name)))
			require.NoError(t, err, msg)

			// Validate in both directions
			assert.Equal(t, want, have, msg)
			assert.Equal(t, have, want, msg)
		})
	})
}
