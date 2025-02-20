// Package doc contains helper functions to run doc examples as testscript
// scripts to keep the docs up to date with the code at the head of the
// repository.
package doc

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/holos-run/holos/cmd"
	"github.com/rogpeppe/go-internal/testscript"

	cue "cuelang.org/go/cmd/cue/cmd"
)

func TestMain(m *testing.M) {
	holosMain := cmd.MakeMain()
	testscript.Main(m, map[string]func(){
		"holos": func() { os.Exit(holosMain()) },
		"cue":   func() { os.Exit(cue.Main()) },
	})
}

func RunOneScript(t *testing.T, dir string, file string) {
	params := testscript.Params{
		Dir:                 "",
		Files:               []string{file},
		RequireExplicitExec: true,
		RequireUniqueNames:  false,
		WorkdirRoot:         filepath.Join(testDir(t), dir),
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

	testscript.Run(t, params)
}

// testDir returns the path of the directory containing the go source file of
// the caller.
func testDir(t *testing.T) string {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		t.Fatal("could not get runtime caller")
	}
	return filepath.Dir(file)
}

// SortedTestScripts returns test scripts in sorted order for sequential
// execution.  Scripts are executed in parallel by default with testscript,
// which is a problem for the docs that build on previous examples.
func SortedTestScripts(t *testing.T, dir string) (files []string) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		// Continue to helpful error on len(files) == 0 below.
	} else if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".txtar") || strings.HasSuffix(name, ".txt") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	if len(files) == 0 {
		t.Fatalf("no txtar nor txt scripts found in dir %s", dir)
	}
	slices.Sort(files)
	return files
}
