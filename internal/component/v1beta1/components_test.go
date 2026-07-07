package v1beta1_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const apiVersion string = "v1beta1"

func TestComponents(t *testing.T) {
	tempDir := testutil.SetupPlatform(t, apiVersion)
	h := testutil.NewComponentHarness(t, tempDir, apiVersion)
	assert.NotEmpty(t, h.Root())

	t.Run("WithNoTasks", func(t *testing.T) {
		leaf := "components/minimal"
		c := h.Component(leaf)
		msg := fmt.Sprintf("Expected %s with no tasks to work, but do nothing", leaf)

		tm, err := c.TypeMeta()
		require.NoError(t, err, msg)

		t.Run("TypeMeta", func(t *testing.T) {
			assert.Equal(t, apiVersion, tm.APIVersion, msg)
			assert.Equal(t, "TaskSet", tm.Kind, msg)
		})

		t.Run("Build", func(t *testing.T) {
			bp, err := c.BuildPlan(tm, holos.NewBuildOpts(h.Root(), leaf, "deploy", t.TempDir()), holos.TagMap{})
			require.NoError(t, err, msg)
			err = bp.Build(h.Ctx())
			require.NoError(t, err, msg)
		})
	})

	t.Run("TaskSet", func(t *testing.T) {
		t.Run("Task", func(t *testing.T) {
			for _, tc := range []string{"command", "helm", "kustomize", "join", "dag"} {
				testComponent(t, h, "task", tc)
			}
		})

		t.Run("Validator", func(t *testing.T) {
			t.Run("Command", func(t *testing.T) {
				t.Run("SecretForbidden", func(t *testing.T) {
					kind := "validator"
					name := "secret"
					path := filepath.Join("components", kind, name)
					leaf := filepath.Join(h.Base(), path)
					c := h.Component(path)
					msg := fmt.Sprintf("Expected %s with %s to fail validation", path, kind)
					tm, err := c.TypeMeta()
					require.NoError(t, err, msg)
					assert.Equal(t, tm.APIVersion, apiVersion)

					t.Run("Build", func(t *testing.T) {
						bp, err := c.BuildPlan(tm, holos.NewBuildOpts(h.Root(), leaf, "deploy", t.TempDir()), holos.TagMap{})
						require.NoError(t, err, msg)
						err = bp.Build(h.Ctx())
						assert.ErrorContains(t, err, "could not run command", msg)
					})
				})
			})
		})
	})
}

// testComponent runs tests against a component directory fixture by building
// the component and asserting the rendered config manifests against a
// want_name.gen.yaml file in the component leaf directory.
func testComponent(t *testing.T, h *testutil.ComponentHarness, kind, name string) {
	t.Run(testutil.Capitalize(name), func(t *testing.T) {
		path := filepath.Join("components", kind, name)
		leaf := filepath.Join(h.Base(), path)
		c := h.Component(path)
		msg := fmt.Sprintf("Expected %s with %s to render config manifests", path, kind)
		tm, err := c.TypeMeta()
		require.NoError(t, err, msg)
		assert.Equal(t, tm.APIVersion, apiVersion)

		t.Run("Build", func(t *testing.T) {
			bp, err := c.BuildPlan(tm, holos.NewBuildOpts(h.Root(), leaf, "deploy", t.TempDir()), holos.TagMap{})
			require.NoError(t, err, msg)
			err = bp.Build(h.Ctx())
			require.NoError(t, err, msg)

			// Validate the rendered manifest
			have, err := h.Load(filepath.Join("deploy", "components", kind, name, fmt.Sprintf("%s.gen.yaml", name)))
			require.NoError(t, err, msg)
			want, err := h.Load(filepath.Join(leaf, fmt.Sprintf("want_%s.gen.yaml", name)))
			require.NoError(t, err, msg)

			// Validate in both directions
			assert.Equal(t, want, have, msg)
			assert.Equal(t, have, want, msg)
		})
	})
}
