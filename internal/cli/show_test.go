package cli_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"testing"

	v1alpha5 "github.com/holos-run/holos/api/core/v1alpha5"
	v1alpha6 "github.com/holos-run/holos/api/core/v1alpha6"
	"github.com/holos-run/holos/internal/cli"
	"github.com/holos-run/holos/internal/generate"
	"github.com/holos-run/holos/internal/platform"
	"github.com/holos-run/holos/internal/testutil"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type harness struct {
	cmd    *cobra.Command
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func newHarness() *harness {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cfg := platform.NewConfig()
	cfg.Stdout = stdout
	cfg.Stderr = stderr
	cmd := cli.NewShowCmd(cfg)
	return &harness{
		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}
}

func (h *harness) Run(ctx context.Context, args ...string) error {
	h.stdout.Reset()
	h.stderr.Reset()
	h.cmd.SetArgs(args)
	return h.cmd.ExecuteContext(ctx)
}

func TestShowAlpha6(t *testing.T) {
	tempDir := t.TempDir()

	// test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Initialize the platform
	if err := generate.GeneratePlatform(ctx, tempDir, "v1alpha6"); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}

	if err := fs.WalkDir(testutil.Fixtures, "fixtures/v1alpha6", util.MakeCopyFunc(ctx, testutil.Fixtures, tempDir)); err != nil {
		t.Fatalf("could not copy fixtures: %v", err)
	}

	t.Run("BuildPlans", func(t *testing.T) {
		t.Run("EmptyPlatform", func(t *testing.T) {
			platformDir := filepath.Join(tempDir, "platform")
			h := newHarness()
			err := h.Run(ctx, "buildplans", platformDir)
			require.NoError(t, err)
		})

		t.Run("SliceComponent", func(t *testing.T) {
			platformDir := filepath.Join(tempDir, "fixtures", "v1alpha6", "platform1")
			// Unmarshal the v1alpha6.BuildPlan we want.
			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha6/platform1/expected_show_buildplans.yaml")
			require.NoError(t, err)
			var want v1alpha6.BuildPlan
			err = yaml.Unmarshal(wantBytes, &want)
			require.NoError(t, err)
			want.BuildContext.RootDir = tempDir
			want.BuildContext.HolosExecutable, err = util.Executable()
			require.NoError(t, err)
			want.BuildContext.LeafDir = "fixtures/v1alpha6/components/slice"

			t.Run("FormatYAML", func(t *testing.T) {
				h := newHarness()
				err = h.Run(ctx, "buildplans", platformDir)
				require.NoError(t, err)
				// Unmarshal what we have.
				var have v1alpha6.BuildPlan
				err = yaml.Unmarshal(h.stdout.Bytes(), &have)
				require.NoError(t, err)

				// Compare them in both directions.
				require.Equal(t, want, have)
				require.Equal(t, have, want)
			})

			t.Run("FormatJSON", func(t *testing.T) {
				h := newHarness()
				err = h.Run(ctx, "buildplans", platformDir, "--format=json")
				require.NoError(t, err)
				// Unmarshal what we have.
				var have v1alpha6.BuildPlan
				err = yaml.Unmarshal(h.stdout.Bytes(), &have)
				require.NoError(t, err)

				// Compare them in both directions.
				require.Equal(t, want, have)
				require.Equal(t, have, want)
			})
		})
	})

}

func TestShowAlpha5(t *testing.T) {
	tempDir := t.TempDir()

	// test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Initialize the platform
	if err := generate.GeneratePlatform(ctx, tempDir, "v1alpha5"); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}

	// Copy fixtures
	if err := fs.WalkDir(testutil.Fixtures, "fixtures/v1alpha5", util.MakeCopyFunc(ctx, testutil.Fixtures, tempDir)); err != nil {
		t.Fatalf("could not copy fixtures: %v", err)
	}

	// https://github.com/holos-run/holos/issues/331
	// ensure holos show components --labels selects correctly.
	t.Run("Selectors", func(t *testing.T) {
		platformDir := filepath.Join(tempDir, "fixtures", "v1alpha5", "issue331", "platform")
		t.Run("ShowPlatform", func(t *testing.T) {
			h := newHarness()
			err := h.Run(ctx, "platform", platformDir)
			require.NoError(t, err)

			var want v1alpha5.Platform
			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha5/issue331/want/platform.yaml")
			require.NoError(t, err)
			err = yaml.Unmarshal(wantBytes, &want)
			require.NoError(t, err)

			var have v1alpha5.Platform
			err = yaml.Unmarshal(h.stdout.Bytes(), &have)
			require.NoError(t, err)

			// Compare them in both directions.
			require.Equal(t, want, have)
			require.Equal(t, have, want)
		})

		t.Run("ShowBuildPlans", func(t *testing.T) {
			h := newHarness()
			err := h.Run(ctx, "buildplans", platformDir)
			require.NoError(t, err)

			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha5/issue331/want/all-buildplans.yaml")
			require.NoError(t, err)
			want := buildPlansA5(t, bytes.NewReader(wantBytes))
			have := buildPlansA5(t, bytes.NewReader(h.stdout.Bytes()))

			// Compare them in both directions.
			require.Equal(t, want, have)
			require.Equal(t, have, want)
		})

		t.Run("SelectorWithOneEqualSign", func(t *testing.T) {
			h := newHarness()
			err := h.Run(ctx, "buildplans", platformDir, "--selector", "app.holos.run/name=empty1-label")
			require.NoError(t, err)

			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha5/issue331/want/buildplans.1.yaml")
			require.NoError(t, err)
			want := buildPlansA5(t, bytes.NewReader(wantBytes))
			have := buildPlansA5(t, bytes.NewReader(h.stdout.Bytes()))

			// Compare them in both directions.
			require.Equal(t, want, have)
			require.Equal(t, have, want)
		})

		t.Run("SelectorWithTwoEqualSigns", func(t *testing.T) {
			h := newHarness()
			err := h.Run(ctx, "buildplans", platformDir, "--selector", "app.holos.run/name==empty2-label")
			require.NoError(t, err)

			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha5/issue331/want/buildplans.2.yaml")
			require.NoError(t, err)
			want := buildPlansA5(t, bytes.NewReader(wantBytes))
			have := buildPlansA5(t, bytes.NewReader(h.stdout.Bytes()))

			// Compare them in both directions.
			require.Equal(t, want, have)
			require.Equal(t, have, want)
		})

		t.Run("SelectorWithNegativeMatch_1", func(t *testing.T) {
			h := newHarness()
			err := h.Run(ctx, "buildplans", platformDir, "--selector", "app.holos.run/name!=empty3-label")
			require.NoError(t, err)

			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha5/issue331/want/buildplans.3.yaml")
			require.NoError(t, err)
			want := buildPlansA5(t, bytes.NewReader(wantBytes))
			have := buildPlansA5(t, bytes.NewReader(h.stdout.Bytes()))

			// Compare them in both directions.
			require.Equal(t, want, have)
			require.Equal(t, have, want)
		})

		t.Run("SelectorWithNegativeMatch_2", func(t *testing.T) {
			h := newHarness()
			err := h.Run(ctx, "buildplans", platformDir, "--selector", "app.holos.run/name!=something-else")
			require.NoError(t, err)

			wantBytes, err := testutil.Fixtures.ReadFile("fixtures/v1alpha5/issue331/want/buildplans.4.yaml")
			require.NoError(t, err)
			want := buildPlansA5(t, bytes.NewReader(wantBytes))
			have := buildPlansA5(t, bytes.NewReader(h.stdout.Bytes()))

			// Compare them in both directions.
			require.Equal(t, want, have)
			require.Equal(t, have, want)
		})
	})
}

// buildPlansA5 decodes v1alpha5 BuildPlans from a yaml stream.
func buildPlansA5(t testing.TB, r io.Reader) (plans []v1alpha5.BuildPlan) {
	t.Helper()
	decoder := yaml.NewDecoder(r)
	for {
		var plan v1alpha5.BuildPlan
		err := decoder.Decode(&plan)
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		plans = append(plans, plan)
	}
	return
}
