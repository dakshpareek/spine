package cmd

import "github.com/spf13/cobra"

// NewRootCmd constructs the base CLI command for ctx.
func NewRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ctx",
		Short: "Extract and maintain your codebase's architecture",
		Long: `ctx helps you maintain a lightweight, up-to-date architectural 
snapshot of your project for AI-assisted development.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Version = version

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
