package main

import (
	"os"
	"path/filepath"
	"testing"

	cue "cuelang.org/go/cmd/cue/cmd"
	"github.com/holos-run/holos/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	holosMain := cmd.MakeMain()
	testscript.Main(m, map[string]func(){
		"holos": func() { os.Exit(holosMain()) },
		"cue":   func() { os.Exit(cue.Main()) },
	})
}

func TestGuides_v1alpha5(t *testing.T) {
	testscript.Run(t, params(filepath.Join("v1alpha5", "guides")))
}

func TestSchemas_v1alpha5(t *testing.T) {
	testscript.Run(t, params(filepath.Join("v1alpha5", "schemas")))
}

func TestIssues_v1alpha5(t *testing.T) {
	testscript.Run(t, params(filepath.Join("v1alpha5", "issues")))
}

func TestCLI(t *testing.T) {
	testscript.Run(t, params("cli"))
}

func params(dir string) testscript.Params {
	return testscript.Params{
		Dir:                 filepath.Join("tests", dir),
		RequireExplicitExec: true,
		RequireUniqueNames:  os.Getenv("HOLOS_WORKDIR_ROOT") == "",
		WorkdirRoot:         os.Getenv("HOLOS_WORKDIR_ROOT"),
		UpdateScripts:       os.Getenv("HOLOS_UPDATE_SCRIPTS") != "",
		Setup: func(env *testscript.Env) error {
			// Just like cmd/cue/cmd.TestScript, set up separate cache and config dirs per test.
			env.Setenv("CUE_CACHE_DIR", filepath.Join(env.WorkDir, "tmp/cachedir"))
			configDir := filepath.Join(env.WorkDir, "tmp/configdir")
			env.Setenv("CUE_CONFIG_DIR", configDir)
			return nil
		},
	}
}
