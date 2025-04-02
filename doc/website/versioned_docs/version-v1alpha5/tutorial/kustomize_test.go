package main

import (
	"path/filepath"
	"testing"
)

// Run these with go test -v to see the verbose names
func TestKustomize(t *testing.T) {
	t.Run("TestKustomize", func(t *testing.T) {
		// Get an ordered list of test script files.
		dir := "_kustomize"
		for _, file := range sortedTestScripts(t, filepath.Join(dir, "examples")) {
			t.Run("examples", func(t *testing.T) {
				runOneScript(t, dir, file)
			})
		}
	})
}
