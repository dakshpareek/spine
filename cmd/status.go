package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/spine/internal/display"
	"github.com/dakshpareek/spine/internal/fs"
	"github.com/dakshpareek/spine/internal/index"
	"github.com/dakshpareek/spine/internal/types"
)

type statusOptions struct {
	verbose bool
	asJSON  bool
}

func newStatusCmd() *cobra.Command {
	opts := statusOptions{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display index summary and file states",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "list stale and missing files")
	cmd.Flags().BoolVar(&opts.asJSON, "json", false, "output status as JSON")

	return cmd
}

func runStatus(opts statusOptions) error {
	wd, err := os.Getwd()
	if err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("determine working directory: %w", err)}
	}

	ctxDir := filepath.Join(wd, ctxDirName)
	if !fs.Exists(ctxDir) {
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("not initialized. Run 'ctx init' first")}
	}

	indexPath := filepath.Join(ctxDir, indexFileName)
	if !fs.Exists(indexPath) {
		return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("missing index.json. Run 'ctx sync' to rebuild")}
	}

	idx, err := index.LoadIndex(indexPath)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	idx.Stats = index.CalculateStats(idx)

	if opts.asJSON {
		data, err := json.MarshalIndent(idx, "", "  ")
		if err != nil {
			return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("encode status json: %w", err)}
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println(display.Bold("Code Context Status"))
	fmt.Println()

	stats := idx.Stats

	fmt.Println("Overview:")
	fmt.Printf("  Total files: %d\n", stats.TotalFiles)
	fmt.Printf("  %s\n", display.Success(fmt.Sprintf("%d current", stats.Current)))
	fmt.Printf("  %s\n", displayWithWarning(stats.Stale, "stale"))
	fmt.Printf("  %s\n", displayWithWarning(stats.Missing, "missing"))
	if stats.PendingGeneration > 0 {
		fmt.Printf("  %s\n", display.Info("%d pending generation", stats.PendingGeneration))
	}

	if !idx.LastSync.IsZero() {
		fmt.Printf("\nLast sync: %s\n", humanizeDuration(time.Since(idx.LastSync)))
	}

	fmt.Printf("Prompt version: %s\n", idx.PromptVersion)

	if opts.verbose {
		printStatusLists(idx)
	}

	return nil
}

func displayWithWarning(count int, label string) string {
	if count > 0 {
		return display.Warning("%d %s", count, label)
	}
	return display.Success("0 %s", label)
}

func humanizeDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	return fmt.Sprintf("%d days ago", days)
}

func printStatusLists(idx *types.Index) {
	var (
		stale   []string
		missing []string
		pending []string
	)

	for path, entry := range idx.Files {
		switch entry.Status {
		case types.StatusStale:
			stale = append(stale, path)
		case types.StatusMissing:
			missing = append(missing, path)
		case types.StatusPendingGeneration:
			pending = append(pending, path)
		}
	}

	sort.Strings(stale)
	sort.Strings(missing)
	sort.Strings(pending)

	printList("Stale", stale)
	printList("Missing", missing)
	printList("Pending generation", pending)
}

func printList(label string, items []string) {
	if len(items) == 0 {
		return
	}
	fmt.Printf("\n%s:\n", label)
	for _, item := range items {
		fmt.Printf("  - %s\n", item)
	}
}
