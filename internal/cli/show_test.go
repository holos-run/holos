package cli_test

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"path/filepath"
	"testing"

	v1alpha6 "github.com/holos-run/holos/api/core/v1alpha6"
	"github.com/holos-run/holos/internal/cli"
	"github.com/holos-run/holos/internal/generate"
	"github.com/holos-run/holos/internal/platform"
	"github.com/holos-run/holos/internal/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

//go:embed all:fixtures
var fsys embed.FS

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

func TestShowBuildPlansAlpha6(t *testing.T) {
	tempDir := t.TempDir()

	// test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Initialize the platform
	if err := generate.GeneratePlatform(ctx, tempDir, "v1alpha6"); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}

	if err := fs.WalkDir(fsys, "fixtures/v1alpha6", util.MakeCopyFunc(ctx, fsys, tempDir)); err != nil {
		t.Fatalf("could not copy fixtures: %v", err)
	}

	t.Run("BlankPlatformWithNoComponents", func(t *testing.T) {
		platformDir := filepath.Join(tempDir, "platform")
		h := newHarness()
		err := h.Run(ctx, "buildplans", platformDir)
		require.NoError(t, err)
	})

	t.Run("platform1_SliceCommand", func(t *testing.T) {
		platformDir := filepath.Join(tempDir, "fixtures", "v1alpha6", "platform1")
		// Unmarshal the v1alpha6.BuildPlan we want.
		wantBytes, err := fsys.ReadFile("fixtures/v1alpha6/platform1/expected_show_buildplans.yaml")
		require.NoError(t, err)
		var want v1alpha6.BuildPlan
		err = yaml.Unmarshal(wantBytes, &want)
		require.NoError(t, err)

		t.Run("YAML", func(t *testing.T) {
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

		t.Run("JSON", func(t *testing.T) {
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
}

func TestShowPlatformAlpha5(t *testing.T) {
	tempDir := t.TempDir()

	// test cancellation
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	// Initialize the platform
	if err := generate.GeneratePlatform(ctx, tempDir, "v1alpha5"); err != nil {
		t.Fatalf("could not generate platform: %v", err)
	}
}
