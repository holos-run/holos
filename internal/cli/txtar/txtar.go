package txtar

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/errors"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/util"
	"github.com/rogpeppe/go-internal/txtar"
	"github.com/spf13/cobra"
)

//go:embed long.txt
var longHelp string

// New returns a new txtar command.
func New(cfg *holos.Config) *cobra.Command {
	cmd := command.New("txtar")
	cmd.Use = "txtar [flags] [FILES...]"
	cmd.Short = "create and extract trivial text-based file archives"
	cmd.Long = longHelp
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
				return errors.Wrap(err)
			}
		}
		if cfg.TxtarQuote() {
			names := make([]string, 0, len(a.Files))
			for _, file := range a.Files {
				if txtar.NeedsQuote(file.Data) {
					data, err := txtar.Quote(file.Data)
					if err != nil {
						return errors.Wrap(err)
					}
					file.Data = data
					names = append(names, file.Name)
				}
			}
			if len(names) > 0 {
				a.Comment = []byte(fmt.Sprintf("unquote %s\n", strings.Join(names, " ")))
			}
		}
		if _, err := cfg.Stdout().Write(txtar.Format(a)); err != nil {
			return errors.Wrap(err)
		}
		return nil
	}
}

// extract files from the configured Stdin to Stdout or the filesystem.
func extract(cfg *holos.Config) error {
	slog.Debug("extracting from stdin...")
	input, err := io.ReadAll(cfg.Stdin())
	if err != nil {
		return errors.Wrap(fmt.Errorf("could not read stdin: %w", err))
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
		return errors.Wrap(fmt.Errorf("idx cannot be 0"))
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
			return errors.Wrap(fmt.Errorf("could not make directory: %w", err))
		}
		if err = os.WriteFile(path, file.Data, 0644); err != nil {
			return errors.Wrap(fmt.Errorf("could not write file: %w", err))
		}
	}
	return
}
