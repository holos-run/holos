package initialize

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/holos-run/holos/internal/cli/command"
	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/logger"
	"github.com/holos-run/holos/internal/schematics"
	"github.com/spf13/cobra"
)

// Config holds configuration parameters for initialize.
type config struct {
	schematic *string
}

// Build the shared configuration and flagset for the init subcommand.
func newConfig() (*config, *flag.FlagSet) {
	cfg := &config{}
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	cfg.schematic = flagSet.String("schematic", "bare", "The name of the schematic being used to initialize the platform.")

	return cfg, flagSet
}

// makeInitFunc returns the internal implementation of the init cli subcommand.
func makeInitializeFunc(_ *holos.Config, cfg *config) command.RunFunc {
	return func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		log := logger.FromContext(ctx)

		// Get the current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
		}

		log.Info("Starting Holos platform initialization...")
		embeddedContent, schematicFiles, err := schematics.GetSchematic(*cfg.schematic)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Printf("Writing schematic: %s\n", *cfg.schematic)
		writeFiles(currentDir, embeddedContent, schematicFiles)
		log.Info("Holos platform initialization complete.")
		return nil
	}
}

// New returns the init subcommand for the root command.
func New(hc *holos.Config) *cobra.Command {
	cmd := command.New("init")

	cfg, flagSet := newConfig()

	cmd.Short = "Initialize a new Holos platform based on a schematic."
	cmd.Flags().AddGoFlagSet(flagSet)
	cmd.RunE = makeInitializeFunc(hc, cfg)

	return cmd
}

func writeFiles(outputPath string, content embed.FS, files []fs.DirEntry) {
	for _, file := range files {
		fmt.Println("File name: ", file.Name())
		fileContent, _ := content.Open(file.Name())
		content, _ := io.ReadAll(fileContent)
		fmt.Println("File content:\n", string(content))
	}
}
