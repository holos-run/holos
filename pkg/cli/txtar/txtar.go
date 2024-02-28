package txtar

import (
	"bytes"
	"fmt"
	"github.com/holos-run/holos/pkg/cli/command"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/holos-run/holos/pkg/util"
	"github.com/holos-run/holos/pkg/wrapper"
	"github.com/spf13/cobra"
	"golang.org/x/tools/txtar"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// New returns a new txtar command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("txtar")
	cmd.Short = "trivial text-based file archives"
	cmd.Long = "writes arguments to stdout otherwise extracts"
	cmd.Args = cobra.MinimumNArgs(0)
	cmd.RunE = makeRunFunc(cfg)
	cmd.Flags().SortFlags = false
	cmd.Flags().AddGoFlagSet(cfg.TxtarFlagSet())
	return cmd
}

func makeRunFunc(cfg *holos.Config) command.RunFunc {
	return func(cmd *cobra.Command, args []string) error {
		// extract an archive
		if len(args) == 0 {
			return extract(cfg)
		}
		// create an archive
		a := &txtar.Archive{}
		for _, name := range args {
			if err := filepath.WalkDir(name, util.MakeWalkFunc(a)); err != nil {
				return wrapper.Wrap(err)
			}
		}
		if _, err := cfg.Stdout().Write(txtar.Format(a)); err != nil {
			return wrapper.Wrap(err)
		}
		return nil
	}
}

// extract files from the configured Stdin to Stdout or the filesystem.
func extract(cfg *holos.Config) error {
	input, err := io.ReadAll(cfg.Stdin())
	if err != nil {
		return wrapper.Wrap(fmt.Errorf("could not read stdin: %w", err))
	}
	archive := txtar.Parse(input)
	if idx := cfg.TxtarIndex(); idx != 0 {
		return printFile(cfg.Stdout(), idx, archive)
	}

	return writeFiles(cfg.Logger(), archive)
}

// printFile prints one file from the txtar archive by index.
func printFile(w io.Writer, idx int, a *txtar.Archive) (err error) {
	if idx == 0 {
		return wrapper.Wrap(fmt.Errorf("idx cannot be 0"))
	}
	if idx > 0 {
		_, err = w.Write(util.EnsureNewline(a.Files[idx-1].Data))
	} else {
		_, err = w.Write(util.EnsureNewline(a.Files[len(a.Files)+idx].Data))
	}
	return
}

// writeFiles writes all files in the archive.
func writeFiles(logger *slog.Logger, a *txtar.Archive) (err error) {
	var header string
	if h := bytes.Split(a.Comment, []byte{'\n'})[:1]; len(h) > 0 {
		header = string(h[0])
	}
	for _, file := range a.Files {
		log := logger.With("header", header, "path", file.Name, "bytes", len(file.Data))
		path := filepath.Join(".", file.Name)
		log.Info("writing: " + file.Name)
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not make directory: %w", err))
		}
		if err = os.WriteFile(path, file.Data, 0644); err != nil {
			return wrapper.Wrap(fmt.Errorf("could not write file: %w", err))
		}
	}
	return
}
