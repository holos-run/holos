package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/holos"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	ExitCode      int    `json:"exitCode"`
	File1         string `json:"file1"`
	File2         string `json:"file2"`
	ExpectedError string `json:"expectedError,omitempty"`
}

func TestCompareBuildPlans(t *testing.T) {
	fixturesDir := filepath.Join("tests", "fixtures", "compare")
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("could not read fixtures directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testName := entry.Name()
		t.Run(testName, func(t *testing.T) {
			testDir := filepath.Join(fixturesDir, testName)

			// Read the want.json file
			wantData, err := os.ReadFile(filepath.Join(testDir, "want.json"))
			if err != nil {
				t.Fatalf("could not read want.json: %v", err)
			}

			var tc testCase
			if err := json.Unmarshal(wantData, &tc); err != nil {
				t.Fatalf("could not parse want.json: %v", err)
			}

			// Prepare the command
			var stdout, stderr bytes.Buffer
			cfg := holos.New(holos.Stdout(&stdout), holos.Stderr(&stderr))
			rootCmd := New(cfg)

			// Build the full file paths
			file1Path := filepath.Join(testDir, tc.File1)
			file2Path := filepath.Join(testDir, tc.File2)

			// Set up the command arguments
			rootCmd.SetArgs([]string{"compare", "buildplans", file1Path, file2Path})

			// Run the command
			err = rootCmd.Execute()

			// Check the exit code
			if tc.ExitCode == 0 {
				assert.NoError(t, err, "command should succeed")
			} else {
				assert.Error(t, err, "command should fail")
				if tc.ExpectedError != "" {
					assert.ErrorContains(t, err, tc.ExpectedError)
				}
			}
		})
	}
}