package main

import (
	"path/filepath"
	"testing"

	"github.com/holos-run/holos/internal/doc"
)

func TestMain(m *testing.M) {
	doc.TestMain(m)
}

// Run these with go test -v to see the verbose names
func TestReplaceAppSet(t *testing.T) {
	// Get an ordered list of test script files.
	dir := "_migrate_appset"
	examples := "examples"
	for _, file := range doc.SortedTestScripts(t, filepath.Join(dir, examples)) {
		t.Run(examples, func(t *testing.T) {
			doc.RunOneScript(t, dir, file)
		})
	}
}
