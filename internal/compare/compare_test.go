package compare

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	ExitCode       int      `json:"exitCode"`
	Name           string   `json:"name,omitempty"`
	Msg            string   `json:"msg,omitempty"`
	File1          string   `json:"file1"`
	File2          string   `json:"file2"`
	ExpectedError  string   `json:"expectedError,omitempty"` // Deprecated: use ExpectedErrors
	ExpectedErrors []string `json:"expectedErrors,omitempty"`
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

		dirName := entry.Name()
		t.Run(dirName, func(t *testing.T) {
			testDir := filepath.Join(fixturesDir, dirName)

			// Read the testcase.json file
			testcaseData, err := os.ReadFile(filepath.Join(testDir, "testcase.json"))
			if err != nil {
				t.Fatalf("could not read testcase.json: %v", err)
			}

			var tc testCase
			if err := json.Unmarshal(testcaseData, &tc); err != nil {
				t.Fatalf("could not parse testcase.json: %v", err)
			}

			// Use the test name if provided, otherwise use directory name
			testName := dirName
			if tc.Name != "" {
				testName = tc.Name
			}

			// Run the test with the appropriate name
			t.Run(testName, func(t *testing.T) {
				// Build the full file paths
				file1Path := filepath.Join(testDir, tc.File1)
				file2Path := filepath.Join(testDir, tc.File2)

				// Create a new comparer and run the comparison
				c := New()
				err := c.BuildPlans(file1Path, file2Path, false)

				// Check the result based on expected exit code
				if tc.ExitCode == 0 {
					assert.NoError(t, err, tc.Msg)
				} else {
					assert.Error(t, err, tc.Msg)
					// Support both old expectedError and new expectedErrors
					if tc.ExpectedError != "" {
						assert.ErrorContains(t, err, tc.ExpectedError, tc.Msg)
					}
					// Check each expected error substring
					for _, expectedErr := range tc.ExpectedErrors {
						assert.ErrorContains(t, err, expectedErr, tc.Msg)
					}
				}
			})
		})
	}
}
