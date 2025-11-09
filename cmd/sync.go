package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/spine/internal/config"
	"github.com/dakshpareek/spine/internal/display"
	"github.com/dakshpareek/spine/internal/fs"
	"github.com/dakshpareek/spine/internal/git"
	"github.com/dakshpareek/spine/internal/hash"
	"github.com/dakshpareek/spine/internal/index"
	"github.com/dakshpareek/spine/internal/scanner"
	"github.com/dakshpareek/spine/internal/skeleton"
	"github.com/dakshpareek/spine/internal/types"
)

type syncOptions struct {
	full    bool
	verbose bool
}

func newSyncCmd() *cobra.Command {
	opts := syncOptions{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Scan codebase for changes and update index",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.full, "full", false, "force full scan (ignore git diff)")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "show detailed file changes")

	return cmd
}

func runSync(opts syncOptions) error {
	wd, err := os.Getwd()
	if err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("determine working directory: %w", err)}
	}

	ctxDir := filepath.Join(wd, ctxDirName)
	if !fs.Exists(ctxDir) {
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("not initialized. Run 'spine init' first")}
	}

	configPath := filepath.Join(ctxDir, configFileName)
	if !fs.Exists(configPath) {
		return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("missing config.json. Run 'spine rebuild --confirm' to restore")}
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	indexPath := filepath.Join(ctxDir, indexFileName)
	if !fs.Exists(indexPath) {
		return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("missing index.json. Run 'spine rebuild --confirm' to restore")}
	}
	idx, err := index.LoadIndex(indexPath)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	scanCfg := *cfg
	scanCfg.RootPath = "."

	fmt.Println(display.Info("Scanning codebase..."))
	files, err := scanner.ScanFiles(scanCfg)
	if err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("scan files: %w", err)}
	}
	fmt.Println(display.Success("%d files scanned", len(files)))

	fileSet := make(map[string]struct{}, len(files))
	for _, f := range files {
		fileSet[f] = struct{}{}
	}

	rootDir := wd

	updateSet := determineUpdateSet(files, opts.full, idx.LastSync, rootDir, scanCfg)

	var (
		modified []string
		added    []string
	)

	for path := range updateSet {
		if _, ok := fileSet[path]; !ok {
			continue
		}

		fullPath := filepath.Join(rootDir, filepath.FromSlash(path))
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("stat %s: %w", path, err)}
		}

		hashValue, err := hash.HashFile(fullPath)
		if err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("hash %s: %w", path, err)}
		}

		existing, ok := idx.Files[path]
		if !ok {
			entry := types.FileEntry{
				Path:         path,
				Hash:         hashValue,
				SkeletonHash: "",
				SkeletonPath: skeleton.PathForSource(path),
				LastModified: info.ModTime().UTC(),
				Status:       types.StatusMissing,
				Type:         scanner.DetectFileType(path),
				Size:         info.Size(),
			}
			idx.Files[path] = entry
			added = append(added, path)
			continue
		}

		if existing.Hash != hashValue {
			existing.Hash = hashValue
			existing.Status = types.StatusStale
			modified = append(modified, path)
		}

		existing.LastModified = info.ModTime().UTC()
		existing.Size = info.Size()
		existing.Type = scanner.DetectFileType(path)
		existing.Path = path
		if existing.SkeletonPath == "" {
			existing.SkeletonPath = skeleton.PathForSource(path)
		}

		idx.Files[path] = existing
	}

	var deleted []string
	for path := range idx.Files {
		if _, ok := fileSet[path]; !ok {
			deleted = append(deleted, path)
		}
	}
	for _, path := range deleted {
		delete(idx.Files, path)
	}

	idx.LastSync = time.Now().UTC()
	idx.Stats = index.CalculateStats(idx)

	if err := index.SaveIndex(idx, indexPath); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	if len(modified)+len(added)+len(deleted) == 0 {
		fmt.Println(display.Success("No changes detected"))
	} else {
		fmt.Println(display.Bold("Changes detected:"))
		if len(modified) > 0 {
			fmt.Printf("  • %d modified (marked stale)\n", len(modified))
		}
		if len(added) > 0 {
			fmt.Printf("  • %d new file(s) (marked missing)\n", len(added))
		}
		if len(deleted) > 0 {
			fmt.Printf("  • %d deleted\n", len(deleted))
		}
	}

	if opts.verbose {
		printDetailedChanges("Modified", modified)
		printDetailedChanges("Added", added)
		printDetailedChanges("Deleted", deleted)
	}

	fmt.Println()
	fmt.Println(display.Bold("Status:"))
	stats := idx.Stats
	fmt.Println(display.Success("%d current", stats.Current))
	if stats.Stale > 0 {
		fmt.Println(display.Warning("%d stale", stats.Stale))
	} else {
		fmt.Println(display.Success("0 stale"))
	}
	if stats.Missing > 0 {
		fmt.Println(display.Warning("%d missing", stats.Missing))
	} else {
		fmt.Println(display.Success("0 missing"))
	}
	if stats.PendingGeneration > 0 {
		fmt.Println(display.Info("%d pending generation", stats.PendingGeneration))
	}
	fmt.Printf("  Total tracked: %d\n", stats.TotalFiles)

	return nil
}

func determineUpdateSet(files []string, forceFull bool, lastSync time.Time, root string, cfg types.Config) map[string]struct{} {
	set := make(map[string]struct{})
	if forceFull {
		for _, f := range files {
			set[f] = struct{}{}
		}
		return set
	}

	if git.IsGitRepo() && !lastSync.IsZero() {
		if modified, err := git.GetModifiedFiles(); err == nil {
			for _, f := range modified {
				normalized := filepath.ToSlash(f)
				set[normalized] = struct{}{}
			}
		} else if !errors.Is(err, git.ErrNotGit) {
			for _, f := range files {
				set[f] = struct{}{}
			}
			return set
		}

		if untracked, err := git.GetUntrackedFiles(); err == nil {
			for _, f := range untracked {
				normalized := filepath.ToSlash(f)
				set[normalized] = struct{}{}
			}
		}

		if len(set) == 0 {
			for _, f := range files {
				set[f] = struct{}{}
			}
		}
		return set
	}

	fallbackRoot := root
	if cfg.RootPath != "." {
		fallbackRoot = filepath.Join(root, cfg.RootPath)
	}

	fallbackFiles, err := git.GetModifiedFilesFallback(fallbackRoot, lastSync)
	if err != nil || len(fallbackFiles) == 0 {
		for _, f := range files {
			set[f] = struct{}{}
		}
		return set
	}

	for _, f := range fallbackFiles {
		set[filepath.ToSlash(f)] = struct{}{}
	}
	return set
}

func printDetailedChanges(label string, items []string) {
	if len(items) == 0 {
		return
	}
	sort.Strings(items)
	fmt.Printf("%s:\n", label)
	for _, item := range items {
		fmt.Printf("  - %s\n", item)
	}
}
