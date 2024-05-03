package schematics

import (
	"embed"
	"fmt"
	"io"
)

// EmbeddedContent holds everything beneath the embedded directory.
//
//go:embed src/*
var EmbeddedContent embed.FS

// func GetSchematic(name string) (embed.FS, error) {
func GetSchematic(name string) error {
	// Check if the directory exists in the embedded content
	schematicDir, err := EmbeddedContent.ReadDir("src/" + name)
	if err != nil {
		return err
	}

	for _, file := range schematicDir {
		fmt.Println("File name: ", file.Name())
		fileContent, _ := EmbeddedContent.Open("src/" + name + "/" + file.Name())
		content, _ := io.ReadAll(fileContent)
		fmt.Println("File content:\n", string(content))
	}

	return nil
}
