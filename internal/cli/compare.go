package cli

import (
	"fmt"

	"github.com/holos-run/holos/internal/compare"
	"github.com/spf13/cobra"
)

// New for the compare command.
func NewCompareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare Holos resources",
		Long:  "Compare Holos resources to verify semantic equivalence",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewCompareBuildPlansCmd())
	cmd.AddCommand(NewCompareYAMLCmd())
	return cmd
}

// New for the compare buildplans subcommand.
func NewCompareBuildPlansCmd() *cobra.Command {
	var backwardsCompatible bool

	cmd := &cobra.Command{
		Use:   "buildplans [file1] [file2]",
		Short: "Compare two BuildPlan files",
		Long:  "Compare two BuildPlan files to verify they are semantically equivalent",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := compare.New()
			return c.BuildPlans(args[0], args[1], backwardsCompatible)
		},
	}

	cmd.Flags().BoolVar(&backwardsCompatible, "backwards-compatible", false, "Enable backwards compatibility mode where file2 may have fields missing from file1")

	return cmd
}

// New for the compare yaml subcommand.
func NewCompareYAMLCmd() *cobra.Command {
	var backwardsCompatible bool

	cmd := &cobra.Command{
		Use:   "yaml [file1] [file2]",
		Short: "Compare two yaml object streams",
		Long:  "Compare two yaml object streams to verify they are structurally equivalent",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := compare.New()
			// TODO(jeff): Add a YAML() function.
			return c.BuildPlans(args[0], args[1], backwardsCompatible)
		},
	}

	cmd.Flags().BoolVar(&backwardsCompatible, "backwards-compatible", false, "Enable backwards compatibility mode where file2 may have fields missing from file1")

	return cmd
}
