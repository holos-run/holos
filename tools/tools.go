//go:build tools

package tools

// Refer to "How can I track tool dependencies for a module?"
// https://go.dev/wiki/Modules

import (
	_ "cuelang.org/go/cmd/cue"
	_ "github.com/princjef/gomarkdoc/cmd/gomarkdoc"
	_ "github.com/rogpeppe/go-internal/cmd/testscript"
	_ "sigs.k8s.io/kustomize/kustomize/v5"
)
