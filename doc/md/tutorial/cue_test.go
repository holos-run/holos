package main

import (
	"path/filepath"
	"testing"
)

// Run these with go test -v to see the verbose names
func TestCue(t *testing.T) {
	t.Run("TestCue", func(t *testing.T) {
		// Get an ordered list of test script files.
		dir := "_cue"
		for _, file := range sortedTestScripts(t, filepath.Join(dir, "examples")) {
			t.Run("examples", func(t *testing.T) {
				runOneScript(t, dir, file)
			})
		}
	})
}
