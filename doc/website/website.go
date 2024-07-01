// Package website embeds the docs website for the server subcommand.  Docs are
// served at /docs similar to how the ui is served at /ui.
package website

import (
	"embed"
	"io/fs"
)

// Output must be the relative path to where the build tool places the static
// site index.html file.
const OutputPath = "build"

//go:embed all:build
var Dist embed.FS

// Root returns the static site root directory.
func Root() fs.FS {
	sub, err := fs.Sub(Dist, OutputPath)
	if err != nil {
		panic(err)
	}
	return sub
}
