package cli

import (
	"github.com/holos-run/holos/internal/compare"
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
