package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	indexpkg "github.com/dakshpareek/spine/internal/index"
	"github.com/dakshpareek/spine/internal/types"
)

func TestStatusVerboseOutput(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "main.go", "package main\n")

	_, _ = executeCommand(t, dir, "init")

	idx := loadIndex(t, dir)
	idx.Files["main.go"] = types.FileEntry{
		Path:   "main.go",
		Status: types.StatusCurrent,
	}
	idx.Files["missing.go"] = types.FileEntry{
		Path:   "missing.go",
		Status: types.StatusMissing,
	}
	idx.Files["pending.go"] = types.FileEntry{
		Path:   "pending.go",
		Status: types.StatusPendingGeneration,
	}
	idx.LastSync = time.Now().Add(-2 * time.Hour)
	if err := indexpkg.SaveIndex(idx, filepath.Join(dir, ".spine", "index.json")); err != nil {
		t.Fatalf("save index: %v", err)
	}

	output := captureOutput(t, func() {
		_, _ = executeCommand(t, dir, "status", "--verbose")
	})
	if !strings.Contains(output, "Missing:") {
		t.Fatalf("expected verbose output to include Missing list")
	}
	if !strings.Contains(output, "Pending generation:") {
		t.Fatalf("expected verbose output to include pending list")
	}
	if !strings.Contains(output, "0 stale") {
		t.Fatalf("expected overview to include zero stale message")
	}
}

func TestStatusJSONOutput(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "main.go", "package main\n")

	_, _ = executeCommand(t, dir, "init")

	output := captureOutput(t, func() {
		_, _ = executeCommand(t, dir, "status", "--json")
	})
	if !strings.Contains(output, "\"version\"") {
		t.Fatalf("expected json output")
	}
}

func TestDisplayWithWarningBranches(t *testing.T) {
	if !strings.Contains(displayWithWarning(0, "items"), "0 items") {
		t.Fatalf("expected success branch")
	}
	if !strings.Contains(displayWithWarning(2, "items"), "2 items") {
		t.Fatalf("expected warning branch")
	}
}

func TestPrintList(t *testing.T) {
	output := captureOutput(t, func() {
		printList("Section", []string{"a", "b"})
	})
	if !strings.Contains(output, "Section:") || !strings.Contains(output, "a") {
		t.Fatalf("expected list output")
	}

	output = captureOutput(t, func() {
		printList("Empty", nil)
	})
	if len(output) != 0 {
		t.Fatalf("expected no output for empty list")
	}
}

func TestRunStatusErrors(t *testing.T) {
	dir := t.TempDir()
	cleanup := changeDir(t, dir)
	defer cleanup()

	if err := runStatus(statusOptions{}); err == nil {
		t.Fatalf("expected error when not initialized")
	}

	// Create .spine but no index to trigger missing index path.
	if err := os.MkdirAll(filepath.Join(dir, ctxDirName), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := runStatus(statusOptions{}); err == nil {
		t.Fatalf("expected error when index missing")
	}
}
