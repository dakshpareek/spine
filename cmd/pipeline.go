package cmd

import (
	"github.com/spf13/cobra"
)

func newPipelineCmd() *cobra.Command {
	syncOpts := syncOptions{}
	genOpts := generateOptions{
		filter: "stale,missing",
	}

	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Run sync and generate in one step",
		Long: `ctx pipeline performs a sync to detect changes and immediately builds the
generate prompt, reducing the manual steps in the daily workflow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runSync(syncOpts); err != nil {
				return err
			}
			return runGenerate(genOpts)
		},
	}

	cmd.Flags().BoolVar(&syncOpts.full, "full", false, "force full scan (ignore git diff)")
	cmd.Flags().BoolVarP(&syncOpts.verbose, "verbose", "v", false, "show detailed file changes during sync")
	cmd.Flags().StringVar(&genOpts.filter, "filter", genOpts.filter, "comma-separated statuses to include (stale,missing)")
	cmd.Flags().StringVar(&genOpts.files, "files", "", "comma-separated list of specific files to include in the prompt")
	cmd.Flags().StringVarP(&genOpts.output, "output", "o", "", "write prompt to file instead of stdout")

	return cmd
}
