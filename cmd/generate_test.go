package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"os"

	"github.com/dakshpareek/spine/internal/types"
)

func TestParseStatusFilter(t *testing.T) {
	defaults, err := parseStatusFilter("")
	if err != nil {
		t.Fatalf("parseStatusFilter error: %v", err)
	}
	if !defaults[types.StatusStale] || !defaults[types.StatusMissing] {
		t.Fatalf("expected default filter to include stale and missing")
	}

	custom, err := parseStatusFilter("current,pending")
	if err != nil {
		t.Fatalf("parseStatusFilter custom error: %v", err)
	}
	if !custom[types.StatusCurrent] || !custom[types.StatusPendingGeneration] {
		t.Fatalf("expected custom filter to include current and pendingGeneration")
	}

	if _, err := parseStatusFilter("unknown"); err == nil {
		t.Fatalf("expected error for unknown status")
	}
}

func TestParseFilesFilter(t *testing.T) {
	set, err := parseFilesFilter("src/app.go, pkg/util.go")
	if err != nil {
		t.Fatalf("parseFilesFilter error: %v", err)
	}
	if len(set) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(set))
	}

	empty, err := parseFilesFilter("")
	if err != nil {
		t.Fatalf("parseFilesFilter empty error: %v", err)
	}
	if empty != nil {
		t.Fatalf("expected nil map for empty input")
	}
}

func TestLanguageFromExtension(t *testing.T) {
	tests := map[string]string{
		"src/app.ts":      "typescript",
		"src/app.tsx":     "typescript",
		"src/app.js":      "javascript",
		"src/app.go":      "go",
		"src/app.py":      "python",
		"src/app.rs":      "rust",
		"src/app.java":    "java",
		"src/app.rb":      "ruby",
		"src/app.cs":      "csharp",
		"src/app.php":     "php",
		"src/app.swift":   "swift",
		"src/app.kt":      "kotlin",
		"src/app.scala":   "scala",
		"src/app.sh":      "bash",
		"src/app.yaml":    "yaml",
		"src/app.yml":     "yaml",
		"src/app.json":    "json",
		"src/app.md":      "markdown",
		"src/app.service": "",
	}

	for path, expected := range tests {
		if got := languageFromExtension(path); got != expected {
			t.Fatalf("languageFromExtension(%q) = %q, expected %q", path, got, expected)
		}
	}
}

func TestBuildPromptOutput(t *testing.T) {
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "src", "example.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(sourcePath, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	idx := &types.Index{
		Files: map[string]types.FileEntry{
			"src/example.go": {
				Path:         "src/example.go",
				SkeletonPath: ".code-context/skeletons/src/example.skeleton.go",
				Status:       types.StatusStale,
				Type:         "service",
			},
		},
	}

	output, err := buildPromptOutput([]string{"src/example.go"}, idx, "template", tempDir)
	if err != nil {
		t.Fatalf("buildPromptOutput error: %v", err)
	}

	if !strings.Contains(output, "### File 1: src/example.go") {
		t.Fatalf("expected file heading in output:\n%s", output)
	}
	if !strings.Contains(output, "```go") {
		t.Fatalf("expected go code fence in output:\n%s", output)
	}
	if !strings.Contains(output, "## Index Updates Required") {
		t.Fatalf("expected instructions section")
	}
}

func TestBuildJSONExport(t *testing.T) {
	idx := &types.Index{
		PromptVersion: "2.1",
		Stats: types.IndexStats{
			TotalFiles: 1,
			Current:    1,
		},
	}

	skeletons := []exportedSkeleton{
		{
			Path:         "src/example.go",
			SkeletonPath: ".code-context/skeletons/src/example.skeleton.go",
			Type:         "service",
			Status:       types.StatusCurrent,
			Content:      "content",
			LastModified: time.Now().UTC(),
		},
	}

	data, err := buildJSONExport(idx, skeletons, time.Now().UTC())
	if err != nil {
		t.Fatalf("buildJSONExport error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal export: %v", err)
	}

	if parsed["count"].(float64) != 1 {
		t.Fatalf("expected count 1, got %v", parsed["count"])
	}
}
