package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dakshpareek/ctx/internal/display"
	"github.com/dakshpareek/ctx/internal/fs"
	"github.com/dakshpareek/ctx/internal/index"
	"github.com/dakshpareek/ctx/internal/types"
)

type exportOptions struct {
	format string
	output string
}

func newExportCmd() *cobra.Command {
	opts := exportOptions{
		format: "markdown",
	}

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export current skeletons for AI-assisted development",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(opts)
		},
	}

	cmd.Flags().StringVar(&opts.format, "format", opts.format, "output format: markdown or json")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "write export to file instead of stdout")

	return cmd
}

func runExport(opts exportOptions) error {
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

	currentPaths := currentSkeletonPaths(idx)
	if len(currentPaths) == 0 {
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("no current skeletons to export")}
	}

	exported, err := readSkeletonContents(currentPaths, idx, wd)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	var output string

	switch strings.ToLower(opts.format) {
	case "markdown", "md":
		output = buildMarkdownExport(idx, exported, now)
	case "json":
		bytes, err := buildJSONExport(idx, exported, now)
		if err != nil {
			return err
		}
		output = string(bytes)
	default:
		return &types.Error{Code: types.ExitCodeUserError, Err: fmt.Errorf("unsupported format: %s", opts.format)}
	}

	if opts.output != "" {
		if err := fs.WriteFile(opts.output, []byte(output)); err != nil {
			return &types.Error{Code: types.ExitCodeFileSystem, Err: err}
		}
		fmt.Println(display.Success("Exported %d skeleton(s)", len(exported)))
		fmt.Println(display.Info("Export saved to %s", opts.output))
		return nil
	}

	fmt.Print(output)
	return nil
}

type exportedSkeleton struct {
	Path         string
	SkeletonPath string
	Type         string
	Status       types.Status
	Content      string
	LastModified time.Time
	Size         int64
}

func currentSkeletonPaths(idx *types.Index) []string {
	var paths []string
	for path, entry := range idx.Files {
		if entry.Status == types.StatusCurrent {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}

func readSkeletonContents(paths []string, idx *types.Index, cwd string) ([]exportedSkeleton, error) {
	result := make([]exportedSkeleton, 0, len(paths))

	for _, path := range paths {
		entry := idx.Files[path]
		if entry.SkeletonPath == "" {
			return nil, &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("missing skeleton path for %s", path)}
		}

		fullPath := filepath.Join(cwd, filepath.FromSlash(entry.SkeletonPath))
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, &types.Error{Code: types.ExitCodeFileSystem, Err: fmt.Errorf("read skeleton %s: %w", entry.SkeletonPath, err)}
		}

		result = append(result, exportedSkeleton{
			Path:         path,
			SkeletonPath: entry.SkeletonPath,
			Type:         entry.Type,
			Status:       entry.Status,
			Content:      string(data),
			LastModified: entry.LastModified,
			Size:         entry.Size,
		})
	}

	return result, nil
}

func buildMarkdownExport(idx *types.Index, skeletons []exportedSkeleton, generatedAt time.Time) string {
	var builder strings.Builder

	builder.WriteString("# Code Context Export\n\n")
	builder.WriteString(fmt.Sprintf("Generated: %s\n\n", generatedAt.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("Prompt Version: %s\n\n", idx.PromptVersion))

	builder.WriteString("## Summary\n\n")
	builder.WriteString(fmt.Sprintf("- Total files tracked: %d\n", idx.Stats.TotalFiles))
	builder.WriteString(fmt.Sprintf("- Current skeletons: %d\n", idx.Stats.Current))
	builder.WriteString(fmt.Sprintf("- Stale skeletons: %d\n", idx.Stats.Stale))
	builder.WriteString(fmt.Sprintf("- Missing skeletons: %d\n", idx.Stats.Missing))
	builder.WriteString(fmt.Sprintf("- Pending generation: %d\n\n", idx.Stats.PendingGeneration))

	builder.WriteString("## Skeletons\n\n")
	for i, skel := range skeletons {
		builder.WriteString(fmt.Sprintf("### %s\n", skel.Path))
		builder.WriteString(fmt.Sprintf("**Skeleton Path:** %s\n", skel.SkeletonPath))
		if skel.Type != "" {
			builder.WriteString(fmt.Sprintf("**Type:** %s\n", skel.Type))
		}
		builder.WriteString(fmt.Sprintf("**Last Modified:** %s\n\n", skel.LastModified.Format(time.RFC3339)))

		builder.WriteString("```")
		builder.WriteString(languageFromExtension(skel.Path))
		builder.WriteString("\n")
		builder.WriteString(skel.Content)
		if !strings.HasSuffix(skel.Content, "\n") {
			builder.WriteString("\n")
		}
		builder.WriteString("```\n\n")

		if i < len(skeletons)-1 {
			builder.WriteString("---\n\n")
		}
	}

	return builder.String()
}

func buildJSONExport(idx *types.Index, skeletons []exportedSkeleton, generatedAt time.Time) ([]byte, error) {
	payload := struct {
		GeneratedAt time.Time                   `json:"generatedAt"`
		Index       *types.Index                `json:"index"`
		Skeletons   map[string]exportedSkeleton `json:"skeletons"`
		Count       int                         `json:"count"`
	}{
		GeneratedAt: generatedAt,
		Index:       idx,
		Skeletons:   make(map[string]exportedSkeleton, len(skeletons)),
		Count:       len(skeletons),
	}

	for _, skel := range skeletons {
		payload.Skeletons[skel.Path] = skel
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, &types.Error{Code: types.ExitCodeData, Err: fmt.Errorf("encode export json: %w", err)}
	}

	data = append(data, '\n')
	return data, nil
}
