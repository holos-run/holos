package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// New for the compare command.
func NewCompareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare Holos resources",
		Long:  "Compare Holos resources to verify semantic equivalence",
	}

	cmd.AddCommand(NewCompareBuildPlansCmd())
	return cmd
}

// New for the compare buildplans subcommand.
func NewCompareBuildPlansCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buildplans [file1] [file2]",
		Short: "Compare two BuildPlan files",
		Long:  "Compare two BuildPlan files to verify they are semantically equivalent",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Not implemented yet
			return fmt.Errorf("compare buildplans command not implemented")
		},
	}
	return cmd
}
