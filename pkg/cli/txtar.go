package cli

import (
	"bytes"
	"fmt"
	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/config"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
	"golang.org/x/tools/txtar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func newTxtarCmd(cfg *config.Config) *cobra.Command {
	cmd := command.New("txtar")
	cmd.Short = "trivial text-based file archives"
	cmd.Long = "writes arguments to stdout otherwise extracts"
	cmd.Args = cobra.MinimumNArgs(0)
	cmd.RunE = makeTxtarRun(cfg)
	cmd.Flags().SortFlags = false
	cmd.Flags().AddGoFlagSet(cfg.TxtarFlagSet())
	return cmd
}

func makeTxtarRun(cfg *config.Config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return txExtract(cfg)
		}
		a := &txtar.Archive{}
		for _, name := range args {
			if err := filepath.WalkDir(name, makeWalkFunc(a)); err != nil {
				return wrapper.Wrap(err)
			}
		}
		cfg.Write(txtar.Format(a))
		return nil
	}
}

func makeWalkFunc(a *txtar.Archive) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return wrapper.Wrap(err)
		}

		if !d.IsDir() {
			if file, err := txFile(path); err != nil {
				return wrapper.Wrap(err)
			} else {
				a.Files = append(a.Files, file)
			}
		}

		return nil
	}
}

func txFile(path string) (file txtar.File, err error) {
	file.Name = path
	file.Data, err = os.ReadFile(path)
	return
}

func txExtract(cfg *config.Config) error {
	input, err := io.ReadAll(cfg.Stdin())
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not read stdin: %w", err))
	}
	archive := txtar.Parse(input)
	header := bytes.Split(archive.Comment, []byte{'\n'})[:1]
	if len(header) == 0 {
		header = append(header, []byte{})
	}

	// Print one file to stdout
	idx := cfg.TxtarIndex()
	if idx > 0 {
		cfg.Write(command.EnsureNewline(archive.Files[idx-1].Data))
		return nil
	}
	if idx < 0 {
		tail := len(archive.Files)
		cfg.Write(command.EnsureNewline(archive.Files[tail+idx].Data))
		return nil
	}

	// Write all files
	for _, file := range archive.Files {
		log := cfg.Logger().With("header", string(header[0]), "path", file.Name, "bytes", len(file.Data))
		path := filepath.Join(".", file.Name)
		log.Info("writing: " + file.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not make directory: %w", err))
		}
		if err := os.WriteFile(path, file.Data, 0644); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not write file: %w", err))
		}
	}
	return nil
}
