package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dakshpareek/ctx/internal/types"
)

func TestValidateStatusExportCleanRebuild(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "main.go", "package main\n")

	_, _ = executeCommand(t, dir, "init")

	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n// change\n"), 0o644); err != nil {
		t.Fatalf("write change: %v", err)
	}

	_, _ = executeCommand(t, dir, "pipeline", "--output", "prompt.md")

	idx := loadIndex(t, dir)
	entry := idx.Files["main.go"]
	if entry.SkeletonPath == "" {
		t.Fatalf("expected skeleton path")
	}

	writeTempFile(t, dir, entry.SkeletonPath, "** skeleton content **\n")

	_, _ = executeCommand(t, dir, "validate", "--fix")

	idx = loadIndex(t, dir)
	if idx.Files["main.go"].Status != types.StatusCurrent {
		t.Fatalf("expected status current, got %s", idx.Files["main.go"].Status)
	}

	_, _ = executeCommand(t, dir, "status", "--json")

	_, _ = executeCommand(t, dir, "export", "--output", "context.md")
	exportData, err := os.ReadFile(filepath.Join(dir, "context.md"))
	if err != nil {
		t.Fatalf("read export: %v", err)
	}
	if !bytes.Contains(exportData, []byte("skeleton content")) {
		t.Fatalf("expected export to contain skeleton content")
	}

	orphan := filepath.Join(dir, ".ctx", "skeletons", "orphan.txt")
	if err := os.WriteFile(orphan, []byte("orphan"), 0o644); err != nil {
		t.Fatalf("write orphan: %v", err)
	}

	_, _ = executeCommand(t, dir, "clean")
	if _, err := os.Stat(orphan); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected orphan removed")
	}

	_, _ = executeCommand(t, dir, "rebuild", "--confirm")
	idx = loadIndex(t, dir)
	if idx.Files["main.go"].Status != types.StatusMissing {
		t.Fatalf("expected status missing after rebuild, got %s", idx.Files["main.go"].Status)
	}
}

func TestExportErrors(t *testing.T) {
	dir := t.TempDir()
	func() {
		cleanup := changeDir(t, dir)
		defer cleanup()
		if err := runExport(exportOptions{}); err == nil {
			t.Fatalf("expected error when not initialized")
		}
	}()

	func() {
		cleanup := changeDir(t, dir)
		defer cleanup()
		_, _ = executeCommand(t, dir, "init")
		if err := runExport(exportOptions{}); err == nil || !strings.Contains(err.Error(), "no current skeletons") {
			t.Fatalf("expected error when no current skeletons")
		}
	}()
}

func TestRunExportStdout(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "main.go", "package main\n")
	_, _ = executeCommand(t, dir, "init")
	_, _ = executeCommand(t, dir, "pipeline", "--output", "prompt.md")

	idx := loadIndex(t, dir)
	entry := idx.Files["main.go"]
	writeTempFile(t, dir, entry.SkeletonPath, "content\n")
	_, _ = executeCommand(t, dir, "validate", "--fix")

	cleanup := changeDir(t, dir)
	defer cleanup()

	output := captureOutput(t, func() {
		if err := runExport(exportOptions{format: "json"}); err != nil {
			t.Fatalf("runExport error: %v", err)
		}
	})
	if !strings.Contains(output, "\"skeletons\"") {
		t.Fatalf("expected json output")
	}
}
