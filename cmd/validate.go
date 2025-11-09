package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/spine/internal/display"
	"github.com/dakshpareek/spine/internal/fs"
	"github.com/dakshpareek/spine/internal/hash"
	"github.com/dakshpareek/spine/internal/index"
	"github.com/dakshpareek/spine/internal/skeleton"
	"github.com/dakshpareek/spine/internal/types"
)

type validateOptions struct {
	fix    bool
	strict bool
}

type validationIssue struct {
	message  string
	resolved bool
}

func newValidateCmd() *cobra.Command {
	opts := validateOptions{}

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Verify index, source files, and skeletons are consistent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.fix, "fix", false, "automatically update index entries for detected issues")
	cmd.Flags().BoolVar(&opts.strict, "strict", false, "exit with error if issues are found")

	return cmd
}

func runValidate(opts validateOptions) error {
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
		return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("missing index.json. Run 'ctx sync' to recreate")}
	}

	idx, err := index.LoadIndex(indexPath)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	fmt.Println("Validating code context...")
	fmt.Println()

	paths := make([]string, 0, len(idx.Files))
	for path := range idx.Files {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var (
		issues        []validationIssue
		removePaths   []string
		staleMarked   int
		missingMarked int
		currentMarked int
		removedCount  int
		modified      bool
	)

	for _, path := range paths {
		entry := idx.Files[path]
		sourcePath := filepath.Join(wd, filepath.FromSlash(path))

		info, err := os.Stat(sourcePath)
		if err != nil {
			if os.IsNotExist(err) {
				message := fmt.Sprintf("%s: source file missing", path)
				if opts.fix {
					removePaths = append(removePaths, path)
					removedCount++
					modified = true
					issues = append(issues, validationIssue{
						message:  message + " (removed from index)",
						resolved: true,
					})
				} else {
					issues = append(issues, validationIssue{
						message:  message,
						resolved: false,
					})
				}
				continue
			}
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("stat %s: %w", path, err)}
		}

		currentHash, err := hash.HashFile(sourcePath)
		if err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("hash %s: %w", path, err)}
		}

		hashChanged := currentHash != entry.Hash
		if hashChanged {
			message := fmt.Sprintf("%s: source hash mismatch", path)
			if opts.fix {
				entry.Hash = currentHash
				entry.LastModified = info.ModTime().UTC()
				entry.Size = info.Size()
				entry.Status = types.StatusStale
				staleMarked++
				modified = true
				issues = append(issues, validationIssue{
					message:  message + " (marked stale)",
					resolved: true,
				})
			} else {
				issues = append(issues, validationIssue{
					message:  message,
					resolved: false,
				})
			}
		}

		skelPath := entry.SkeletonPath
		if skelPath == "" {
			skelPath = skeleton.PathForSource(path)
		}
		entry.SkeletonPath = skelPath

		fullSkeletonPath := filepath.Join(wd, filepath.FromSlash(skelPath))

		skeletonExists := fs.Exists(fullSkeletonPath)
		skeletonHashChanged := false

		if !skeletonExists {
			message := fmt.Sprintf("%s: skeleton file missing", path)
			if opts.fix {
				entry.Status = types.StatusMissing
				entry.SkeletonHash = ""
				missingMarked++
				modified = true
				issues = append(issues, validationIssue{
					message:  message + " (marked missing)",
					resolved: true,
				})
			} else {
				issues = append(issues, validationIssue{
					message:  message,
					resolved: false,
				})
			}
		} else {
			skeletonHash, err := hash.HashFile(fullSkeletonPath)
			if err != nil {
				return &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("hash skeleton %s: %w", skelPath, err)}
			}

			if entry.SkeletonHash == "" {
				if opts.fix {
					entry.SkeletonHash = skeletonHash
					modified = true
				}
			} else if entry.SkeletonHash != skeletonHash {
				message := fmt.Sprintf("%s: skeleton hash mismatch", path)
				if opts.fix {
					entry.SkeletonHash = skeletonHash
					skeletonHashChanged = true
					modified = true
					issues = append(issues, validationIssue{
						message:  message + " (skeleton hash updated)",
						resolved: true,
					})
				} else {
					issues = append(issues, validationIssue{
						message:  message,
						resolved: false,
					})
				}
			}

			if opts.fix && !hashChanged && entry.SkeletonHash != "" &&
				entry.Status != types.StatusMissing {
				if entry.Status == types.StatusPendingGeneration || skeletonHashChanged {
					entry.Status = types.StatusCurrent
					currentMarked++
					modified = true
				}
			}
		}

		if opts.fix {
			idx.Files[path] = entry
		}
	}

	for _, path := range removePaths {
		delete(idx.Files, path)
	}

	if opts.fix && modified {
		idx.Stats = index.CalculateStats(idx)
		if err := index.SaveIndex(idx, indexPath); err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
		}
	}

	if len(issues) == 0 {
		fmt.Println(display.Success("No issues found."))
	} else {
		fmt.Println("Issues found:")
		for _, issue := range issues {
			if issue.resolved {
				fmt.Printf("  %s\n", display.Success(issue.message))
			} else {
				fmt.Printf("  %s\n", display.Warning(issue.message))
			}
		}
	}

	if opts.fix && (staleMarked > 0 || missingMarked > 0 || removedCount > 0 || currentMarked > 0) {
		fmt.Println()
		fmt.Println("Summary:")
		if staleMarked > 0 {
			fmt.Printf("  %d file(s) marked stale\n", staleMarked)
		}
		if missingMarked > 0 {
			fmt.Printf("  %d file(s) marked missing\n", missingMarked)
		}
		if currentMarked > 0 {
			fmt.Printf("  %d file(s) marked current\n", currentMarked)
		}
		if removedCount > 0 {
			fmt.Printf("  %d file(s) removed from index\n", removedCount)
		}
	}

	if len(issues) > 0 {
		fmt.Println("\nRun 'ctx generate' to resolve outstanding issues.")
	}

	if len(issues) > 0 && opts.strict {
		return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("%d validation issue(s) detected", len(issues))}
	}

	return nil
}
