package schematics

import (
	"embed"
	"io/fs"
)

// EmbeddedContent holds everything beneath the embedded directory.
//
//go:embed src/*
var EmbeddedContent embed.FS

// func GetSchematic(name string) (embed.FS, error) {
func GetSchematic(name string) (embed.FS, []fs.DirEntry, error) {
	// Check if the directory exists in the embedded content
	schematicDir, err := EmbeddedContent.ReadDir("src/" + name)
	if err != nil {
		return EmbeddedContent, schematicDir, err
	} else {
		return EmbeddedContent, schematicDir, nil
	}
}
