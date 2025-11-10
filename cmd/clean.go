package cmd

import (
	"fmt"
	stdfs "io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/ctx/internal/display"
	"github.com/dakshpareek/ctx/internal/fs"
	"github.com/dakshpareek/ctx/internal/index"
	"github.com/dakshpareek/ctx/internal/types"
)

func newCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove orphaned skeleton files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClean()
		},
	}
}

func runClean() error {
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

	fmt.Println("Cleaning orphaned skeletons...")

	skeletonDir := filepath.Join(ctxDir, "skeletons")
	referenced := make(map[string]struct{}, len(idx.Files))
	for _, entry := range idx.Files {
		if entry.SkeletonPath == "" {
			continue
		}
		path := filepath.Join(wd, filepath.FromSlash(entry.SkeletonPath))
		referenced[path] = struct{}{}
	}

	var (
		orphaned int
		dirs     []string
	)

	filepath.WalkDir(skeletonDir, func(path string, d stdfs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path != skeletonDir {
				dirs = append(dirs, path)
			}
			return nil
		}

		if _, ok := referenced[path]; !ok {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
			orphaned++
		}
		return nil
	})

	// Remove directories deepest first.
	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) > len(dirs[j])
	})

	removedDirs := 0
	for _, dir := range dirs {
		if err := os.Remove(dir); err == nil {
			removedDirs++
		}
	}

	fmt.Printf("  Removed %d orphaned file(s)\n", orphaned)
	fmt.Printf("  Removed %d empty director(ies)\n", removedDirs)
	fmt.Println(display.Success("Clean complete"))

	return nil
}
