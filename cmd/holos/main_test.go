package main

import (
	"os"
	"path/filepath"
	"testing"

	cue "cuelang.org/go/cmd/cue/cmd"
	"github.com/holos-run/holos/internal/cli"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"holos": cli.MakeMain(),
		"cue":   cue.Main,
	}))
}

func TestGuides(t *testing.T) {
	testscript.Run(t, params(filepath.Join("v1alpha4", "guides")))
}

func TestCLI(t *testing.T) {
	testscript.Run(t, params("cli"))
}

func params(dir string) testscript.Params {
	return testscript.Params{
		Dir:                 filepath.Join("tests", dir),
		RequireExplicitExec: true,
		RequireUniqueNames:  true,
		Setup: func(env *testscript.Env) error {
			// Just like cmd/cue/cmd.TestScript, set up separate cache and config dirs per test.
			env.Setenv("CUE_CACHE_DIR", filepath.Join(env.WorkDir, "tmp/cachedir"))
			configDir := filepath.Join(env.WorkDir, "tmp/configdir")
			env.Setenv("CUE_CONFIG_DIR", configDir)
			return nil
		},
	}
}
