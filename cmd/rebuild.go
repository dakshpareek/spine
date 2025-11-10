package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/ctx/internal/config"
	"github.com/dakshpareek/ctx/internal/display"
	"github.com/dakshpareek/ctx/internal/fs"
	"github.com/dakshpareek/ctx/internal/hash"
	"github.com/dakshpareek/ctx/internal/index"
	"github.com/dakshpareek/ctx/internal/scanner"
	"github.com/dakshpareek/ctx/internal/skeleton"
	"github.com/dakshpareek/ctx/internal/types"
)

type rebuildOptions struct {
	confirm bool
}

func newRebuildCmd() *cobra.Command {
	opts := rebuildOptions{}

	cmd := &cobra.Command{
		Use:   "rebuild",
		Short: "Reset index and skeletons (destructive)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRebuild(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.confirm, "confirm", false, "confirm rebuild and proceed without prompt")

	return cmd
}

func runRebuild(opts rebuildOptions) error {
	if !opts.confirm {
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("rebuild requires --confirm flag to proceed")}
	}

	wd, err := os.Getwd()
	if err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("determine working directory: %w", err)}
	}

	ctxDir := filepath.Join(wd, ctxDirName)
	if !fs.Exists(ctxDir) {
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("not initialized. Run 'ctx init' first")}
	}

	configPath := filepath.Join(ctxDir, configFileName)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	skeletonDir := filepath.Join(ctxDir, "skeletons")
	deletedCount, err := purgeSkeletons(skeletonDir)
	if err != nil {
		return err
	}
	if err := fs.EnsureDir(skeletonDir); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	fmt.Println(display.Warning("Rebuilding context..."))

	scanCfg := *cfg
	scanCfg.RootPath = "."

	files, err := scanner.ScanFiles(scanCfg)
	if err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("scan files: %w", err)}
	}

	idx := index.CreateEmptyIndex()
	idx.Config = *cfg

	for _, relPath := range files {
		fullPath := filepath.Join(wd, relPath)
		info, err := os.Stat(fullPath)
		if err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("stat %s: %w", relPath, err)}
		}

		hashValue, err := hash.HashFile(fullPath)
		if err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("hash %s: %w", relPath, err)}
		}

		entry := types.FileEntry{
			Path:         relPath,
			Hash:         hashValue,
			SkeletonHash: "",
			SkeletonPath: skeleton.PathForSource(relPath),
			LastModified: info.ModTime().UTC(),
			Status:       types.StatusMissing,
			Type:         scanner.DetectFileType(relPath),
			Size:         info.Size(),
		}
		idx.Files[relPath] = entry
	}

	idx.LastSync = time.Now().UTC()
	idx.Stats = index.CalculateStats(idx)

	indexPath := filepath.Join(ctxDir, indexFileName)
	if err := index.SaveIndex(idx, indexPath); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	fmt.Printf("  %s\n", display.Success(fmt.Sprintf("Deleted %d skeleton file(s)", deletedCount)))
	fmt.Printf("  %s\n", display.Success("Reset index"))
	fmt.Printf("  %s\n", display.Info("Scanning codebase..."))
	fmt.Printf("Found %d files (all marked missing)\n", len(files))
	fmt.Println()
	fmt.Println("Run 'ctx generate' to recreate skeletons.")

	return nil
}

func purgeSkeletons(root string) (int, error) {
	if !fs.Exists(root) {
		return 0, nil
	}

	count := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		count++
		return nil
	})
	if err != nil {
		return 0, &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("scan skeletons: %w", err)}
	}

	if err := os.RemoveAll(root); err != nil {
		return 0, &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("delete skeletons: %w", err)}
	}

	return count, nil
}
