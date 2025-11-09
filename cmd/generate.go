package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/spine/internal/display"
	"github.com/dakshpareek/spine/internal/fs"
	"github.com/dakshpareek/spine/internal/index"
	"github.com/dakshpareek/spine/internal/scanner"
	"github.com/dakshpareek/spine/internal/skeleton"
	"github.com/dakshpareek/spine/internal/types"
)

type generateOptions struct {
	filter string
	files  string
	output string
}

func newGenerateCmd() *cobra.Command {
	opts := generateOptions{
		filter: "stale,missing",
	}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate AI prompts for creating/updating skeletons",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(opts)
		},
	}

	cmd.Flags().StringVar(&opts.filter, "filter", opts.filter, "comma-separated statuses to include (stale,missing)")
	cmd.Flags().StringVar(&opts.files, "files", "", "comma-separated list of specific files to include")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "write prompt to file instead of stdout")

	return cmd
}

func runGenerate(opts generateOptions) error {
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
		return &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("missing index.json. Run 'ctx sync' to restore")}
	}
	idx, err := index.LoadIndex(indexPath)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	statuses, err := parseStatusFilter(opts.filter)
	if err != nil {
		return &types.Error{Code: types.ExitCodeUserError, Err: err}
	}

	fileFilter, err := parseFilesFilter(opts.files)
	if err != nil {
		return &types.Error{Code: types.ExitCodeUserError, Err: err}
	}

	var selected []string
	if len(fileFilter) > 0 {
		for file := range fileFilter {
			entry, exists := idx.Files[file]
			if !exists {
				return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("file not tracked in index: %s", file)}
			}
			if !statuses[entry.Status] {
				return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("file %s does not match filter statuses", file)}
			}
			selected = append(selected, file)
		}
		sort.Strings(selected)
	} else {
		for path, entry := range idx.Files {
			if statuses[entry.Status] {
				selected = append(selected, path)
			}
		}
		sort.Strings(selected)
	}

	if len(selected) == 0 {
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("no files match the requested filters")}
	}

	promptTemplate, err := skeleton.LoadPromptTemplate(idx.Config)
	if err != nil {
		return &types.Error{Code: types.ExitCodeData, Err: err}
	}

	output, err := buildPromptOutput(selected, idx, promptTemplate, wd)
	if err != nil {
		return err
	}

	for _, path := range selected {
		entry := idx.Files[path]
		entry.Status = types.StatusPendingGeneration
		entry.LastModified = entry.LastModified.UTC()
		if entry.SkeletonPath == "" {
			entry.SkeletonPath = skeleton.PathForSource(path)
		}
		idx.Files[path] = entry
	}

	idx.LastSync = time.Now().UTC()
	idx.Stats = index.CalculateStats(idx)

	if err := index.SaveIndex(idx, indexPath); err != nil {
		return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
	}

	if opts.output != "" {
		if err := fs.WriteFile(opts.output, []byte(output)); err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
		}
		fmt.Println(display.Success("Generated prompts for %d file(s)", len(selected)))
		fmt.Println(display.Info("Prompt saved to %s", opts.output))
		fmt.Println("Next steps:")
		fmt.Println("  1. Paste the prompt into your AI assistant")
		fmt.Println("  2. Create/update skeleton files at the specified paths")
		fmt.Println("  3. Update index entries with new skeleton hashes")
		return nil
	}

	fmt.Print(output)

	fmt.Fprintln(os.Stderr, display.Success("Generated prompts for %d file(s)", len(selected)))
	fmt.Fprintln(os.Stderr, display.Info("Copy the prompt above into your AI assistant"))
	fmt.Fprintln(os.Stderr, "Next steps:")
	fmt.Fprintln(os.Stderr, "  1. Create/update skeleton files at the specified paths")
	fmt.Fprintln(os.Stderr, "  2. Update index.json with new skeleton hashes and statuses")
	fmt.Fprintln(os.Stderr, "  3. Run 'ctx status' to verify progress")

	return nil
}

func buildPromptOutput(paths []string, idx *types.Index, promptTemplate, cwd string) (string, error) {
	var builder strings.Builder

	builder.WriteString("# Code Context Skeleton Generation\n\n")
	builder.WriteString("You are generating structural skeletons for this codebase. Follow these instructions carefully.\n\n")

	builder.WriteString("## Instructions\n\n")
	builder.WriteString("1. For each file below, generate a skeleton using the provided template.\n")
	builder.WriteString("2. Create or update skeleton files at the specified paths.\n")
	builder.WriteString("3. Update `.code-context/index.json` with the new skeleton hash and status.\n\n")

	builder.WriteString("## Skeleton Generation Template\n\n")
	builder.WriteString("```text\n")
	builder.WriteString(promptTemplate)
	if !strings.HasSuffix(promptTemplate, "\n") {
		builder.WriteString("\n")
	}
	builder.WriteString("```\n\n")

	builder.WriteString("## Files to Process\n\n")

	for i, path := range paths {
		entry := idx.Files[path]
		sourcePath := filepath.Join(cwd, filepath.FromSlash(path))
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return "", &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("read %s: %w", path, err)}
		}

		entryType := entry.Type
		if entryType == "" {
			entryType = scanner.DetectFileType(path)
		}

		builder.WriteString(fmt.Sprintf("### File %d: %s\n", i+1, path))
		builder.WriteString(fmt.Sprintf("**Status:** %s\n", entry.Status))
		if entryType != "" {
			builder.WriteString(fmt.Sprintf("**Type:** %s\n", entryType))
		}
		builder.WriteString(fmt.Sprintf("**Skeleton Path:** %s\n\n", entry.SkeletonPath))

		builder.WriteString("**Source Code:**\n")
		lang := languageFromExtension(path)
		builder.WriteString("```")
		builder.WriteString(lang)
		builder.WriteString("\n")
		builder.Write(content)
		if len(content) == 0 || content[len(content)-1] != '\n' {
			builder.WriteString("\n")
		}
		builder.WriteString("```\n\n")

		if i < len(paths)-1 {
			builder.WriteString("---\n\n")
		}
	}

	builder.WriteString("## Index Updates Required\n\n")
	builder.WriteString("After generating skeletons, update `.code-context/index.json`:\n\n")
	builder.WriteString("1. Write the skeleton file to the specified path.\n")
	builder.WriteString("2. Calculate the SHA-256 hash of the skeleton content.\n")
	builder.WriteString("3. Update the file entry with `status: \"current\"`, the new `skeletonHash`, and the current timestamp for `lastModified`.\n")
	builder.WriteString("4. Recalculate the index stats to reflect the changes.\n\n")
	builder.WriteString("**Index Path:** .code-context/index.json\n\n")

	builder.WriteString("## Verification\n\n")
	builder.WriteString("After completion, run `ctx status` to verify that all files are marked current.\n")

	return builder.String(), nil
}

func parseStatusFilter(raw string) (map[types.Status]bool, error) {
	if strings.TrimSpace(raw) == "" {
		return map[types.Status]bool{
			types.StatusStale:   true,
			types.StatusMissing: true,
		}, nil
	}

	statuses := map[types.Status]bool{}
	for _, part := range strings.Split(raw, ",") {
		value := strings.TrimSpace(strings.ToLower(part))
		switch value {
		case "stale":
			statuses[types.StatusStale] = true
		case "missing":
			statuses[types.StatusMissing] = true
		case "pending", "pendinggeneration":
			statuses[types.StatusPendingGeneration] = true
		case "current":
			statuses[types.StatusCurrent] = true
		case "":
			continue
		default:
			return nil, fmt.Errorf("unknown status filter: %s", part)
		}
	}

	if len(statuses) == 0 {
		return nil, fmt.Errorf("no valid statuses provided")
	}
	return statuses, nil
}

func parseFilesFilter(raw string) (map[string]struct{}, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	files := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		value = filepath.ToSlash(value)
		files[value] = struct{}{}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no valid file paths provided")
	}

	return files, nil
}

func languageFromExtension(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".rb":
		return "ruby"
	case ".cs":
		return "csharp"
	case ".php":
		return "php"
	case ".swift":
		return "swift"
	case ".kt":
		return "kotlin"
	case ".scala":
		return "scala"
	case ".sh":
		return "bash"
	case ".yml", ".yaml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	default:
		return ""
	}
}
