package main

import (
	"github.com/holos-run/holos/pkg/cli"
	"github.com/rogpeppe/go-internal/testscript"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"holos": cli.MakeMain(),
	}))
}

func TestGetSecrets(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
	})
}
