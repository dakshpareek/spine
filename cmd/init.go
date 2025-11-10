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

const (
	ctxDirName         = ".ctx"
	legacyCtxDirName   = ".spine"
	configFileName     = "config.json"
	indexFileName      = "index.json"
	skeletonPromptName = skeleton.PromptFileName
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the .ctx/ workspace",
		Long: `ctx init bootstraps the ctx workspace by creating .ctx/,
writing default configuration, and preparing an index for all tracked files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	}

	return cmd
}

func runInit() error {
	wd, err := os.Getwd()
	if err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("determine working directory: %w", err)}
	}

	ctxDir := filepath.Join(wd, ctxDirName)
	legacyDir := filepath.Join(wd, legacyCtxDirName)

	// Check for existing .ctx directory
	if fs.Exists(ctxDir) {
		return &types.Error{
			Code: types.ExitCodeUserError,
			Err:  fmt.Errorf("already initialized. Use 'ctx rebuild' to reset"),
		}
	}

	// Migrate from .spine to .ctx if needed
	if fs.Exists(legacyDir) {
		fmt.Println(display.Info("Migrating from .spine/ to .ctx/..."))
		if err := os.Rename(legacyDir, ctxDir); err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("migrate .spine to .ctx: %w", err)}
		}
		fmt.Println(display.Success("Migrated to .ctx/"))
		return nil
	}

	if err := fs.EnsureDir(ctxDir); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	if err := fs.EnsureDir(filepath.Join(ctxDir, "skeletons")); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	cfg := config.GetDefaultConfig()
	cfg.RootPath = "."
	configPath := filepath.Join(ctxDir, configFileName)
	if err := fs.WriteJSON(configPath, cfg); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	promptPath := filepath.Join(ctxDir, skeletonPromptName)
	if err := fs.WriteFile(promptPath, []byte(skeleton.DefaultPrompt())); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	if err := fs.EnsureGitignoreEntry(filepath.Join(wd, ".gitignore"), ctxDirName+"/"); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	fmt.Println(display.Success("Initialized .ctx/"))
	fmt.Println(display.Success("Updated .gitignore"))
	fmt.Println(display.Info("Scanning codebase..."))

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

		fileHash, err := hash.HashFile(fullPath)
		if err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("hash %s: %w", relPath, err)}
		}

		entry := types.FileEntry{
			Path:         relPath,
			Hash:         fileHash,
			SkeletonHash: "",
			SkeletonPath: skeleton.PathForSource(relPath),
			LastModified: info.ModTime().UTC(),
			Status:       types.StatusMissing,
			Type:         scanner.DetectFileType(relPath),
			Size:         info.Size(),
		}

		if idx.Files == nil {
			idx.Files = make(map[string]types.FileEntry)
		}
		idx.Files[relPath] = entry
	}

	idx.LastSync = time.Now().UTC()
	idx.Stats = index.CalculateStats(idx)

	indexPath := filepath.Join(ctxDir, indexFileName)
	if err := index.SaveIndex(idx, indexPath); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	fmt.Printf("  Found %d files (all marked missing)\n", len(files))
	fmt.Println(display.Success("Index created"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run 'ctx generate' to create skeleton prompts")
	fmt.Println("  2. Feed prompts to your AI assistant")
	fmt.Println("  3. Run 'ctx status' to check progress")

	return nil
}
