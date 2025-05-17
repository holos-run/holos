package compare

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	ExitCode      int    `json:"exitCode"`
	File1         string `json:"file1"`
	File2         string `json:"file2"`
	ExpectedError string `json:"expectedError,omitempty"`
}

func TestBuildPlans(t *testing.T) {
	fixturesDir := "testdata"
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

			// Read the testcase.json file
			testcaseData, err := os.ReadFile(filepath.Join(testDir, "testcase.json"))
			if err != nil {
				t.Fatalf("could not read testcase.json: %v", err)
			}

			var tc testCase
			if err := json.Unmarshal(testcaseData, &tc); err != nil {
				t.Fatalf("could not parse testcase.json: %v", err)
			}

			// Build the full file paths
			file1Path := filepath.Join(testDir, tc.File1)
			file2Path := filepath.Join(testDir, tc.File2)

			// Create a new comparer and run the comparison
			c := New()
			err = c.BuildPlans(file1Path, file2Path)

			// Check the result based on expected exit code
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