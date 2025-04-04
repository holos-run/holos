package v1alpha6_test

import (
	"testing"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const apiVersion string = "v1alpha6"

func TestComponents(t *testing.T) {
	tempDir := testutil.SetupPlatform(t, apiVersion)

	t.Run("Minimal", func(t *testing.T) {
		msg := "Expected a minimal component to work, but do nothing"
		h := testutil.NewComponentHarness(t, tempDir, apiVersion)
		root := h.Root()
		assert.NotEmpty(t, root)

		componentPath := "components/minimal"
		c := h.Component(componentPath)

		t.Run("TypeMeta", func(t *testing.T) {
			tm, err := c.TypeMeta()
			require.NoError(t, err, msg)
			assert.Equal(t, apiVersion, tm.APIVersion)
			assert.Equal(t, "BuildPlan", tm.Kind)
		})

		t.Run("BuildPlan", func(t *testing.T) {
			tm, err := c.TypeMeta()
			require.NoError(t, err, msg)
			bp, err := c.BuildPlan(tm, holos.NewBuildOpts(componentPath))
			require.NoError(t, err, msg)
			err = bp.Build(h.Ctx())
			require.NoError(t, err, msg)
		})
	})
}
