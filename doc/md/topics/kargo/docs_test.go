package kargo

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/holos-run/holos/cmd"
	"github.com/rogpeppe/go-internal/testscript"

	cue "cuelang.org/go/cmd/cue/cmd"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"holos": cmd.MakeMain(),
		"cue":   cue.Main,
	}))
}

// Run these with go test -v to see the verbose names
func TestKargo(t *testing.T) {
	// The rest of the tests will run in add-on-promoter/script-setup/kargo-demo
	// which is a clone of the demo repository.
	addOnPromoterTests := []struct {
		name    string
		example string
	}{
		// Add example test cases in the same order as the document.  Each test case
		// should align to a document section.
		{"Setup", "setup"},
		{"HolosVersion", "holos-version"},
		{"GitURL", "git-url"},
		{"CertManager", "cert-manager"},
	}

	for _, tt := range addOnPromoterTests {
		t.Run("AddOnPromoter", func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				// We use an underscore so Docusaurus ignores the directory.
				testScript(t, "_add-on-promoter", tt.example)
			})
		})
	}
}

func testScript(t *testing.T, dir string, sub string) {
	workdirRoot := filepath.Join(testDir(t), dir)
	fullPath := filepath.Join(workdirRoot, sub)
	p := params(fullPath)
	p.RequireUniqueNames = false
	p.WorkdirRoot = workdirRoot
	testscript.Run(t, p)
}

// testDir returns the path of the directory containing the test cases.
func testDir(t *testing.T) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get runtime caller")
	}
	return filepath.Dir(file)
}

func params(dir string) testscript.Params {
	return testscript.Params{
		Dir:                 dir,
		RequireExplicitExec: true,
		RequireUniqueNames:  os.Getenv("HOLOS_WORKDIR_ROOT") == "",
		WorkdirRoot:         os.Getenv("HOLOS_WORKDIR_ROOT"),
		UpdateScripts:       os.Getenv("HOLOS_UPDATE_SCRIPTS") != "",
		Setup: func(env *testscript.Env) error {
			// Needed for update.sh to determine if we need to update output files.
			env.Setenv("HOLOS_UPDATE_SCRIPTS", os.Getenv("HOLOS_UPDATE_SCRIPTS"))
			// Just like cmd/cue/cmd.TestScript, set up separate cache and config dirs per test.
			env.Setenv("CUE_CACHE_DIR", filepath.Join(env.WorkDir, "tmp/cachedir"))
			configDir := filepath.Join(env.WorkDir, "tmp/configdir")
			env.Setenv("CUE_CONFIG_DIR", configDir)
			return nil
		},
	}
}
