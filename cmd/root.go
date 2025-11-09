package cmd

import "github.com/spf13/cobra"

// NewRootCmd constructs the base CLI command for spine.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spine",
		Short: "Extract and maintain your codebase's architecture",
		Long: `Spine helps you maintain a lightweight, up-to-date architectural 
snapshot of your project for AI-assisted development.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Version = "1.0.0"

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_ = cmd.Help()
		return nil
	}

	cmd.AddCommand(
		newInitCmd(),
		newSyncCmd(),
		newStatusCmd(),
		newValidateCmd(),
		newCleanCmd(),
		newRebuildCmd(),
		newGenerateCmd(),
		newPipelineCmd(),
		newExportCmd(),
	)

	return cmd
}
