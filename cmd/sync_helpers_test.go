package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dakshpareek/spine/internal/types"
)

func TestDetermineUpdateSetForceFull(t *testing.T) {
	files := []string{"a.go", "b.go"}
	set := determineUpdateSet(files, true, time.Now(), t.TempDir(), types.Config{})
	if len(set) != 2 {
		t.Fatalf("expected all files included")
	}
}

func TestDetermineUpdateSetFallback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.go")
	if err := os.WriteFile(path, []byte("package main"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	files := []string{"file.go"}
	set := determineUpdateSet(files, false, time.Now().Add(-time.Hour), dir, types.Config{RootPath: ""})
	if len(set) == 0 {
		t.Fatalf("expected fallback to include file")
	}
}

func TestPrintDetailedChanges(t *testing.T) {
	captured := captureOutput(t, func() {
		printDetailedChanges("Modified", []string{"a.go", "b.go"})
	})
	if len(captured) == 0 {
		t.Fatalf("expected output")
	}

	captured = captureOutput(t, func() {
		printDetailedChanges("Modified", nil)
	})
	if captured != "" {
		t.Fatalf("expected no output for empty slice")
	}
}

func TestDetermineUpdateSetGit(t *testing.T) {
	dir := t.TempDir()
	if err := runGitCommand(dir, "init"); err != nil {
		t.Skipf("git init failed: %v", err)
	}
	_ = runGitCommand(dir, "config", "user.email", "test@example.com")
	_ = runGitCommand(dir, "config", "user.name", "Test User")

	tracked := filepath.Join(dir, "tracked.go")
	if err := os.WriteFile(tracked, []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write tracked: %v", err)
	}
	if err := runGitCommand(dir, "add", "tracked.go"); err != nil {
		t.Skipf("git add failed: %v", err)
	}
	if err := runGitCommand(dir, "commit", "-m", "init"); err != nil {
		t.Skipf("git commit failed: %v", err)
	}

	if err := os.WriteFile(tracked, []byte("package main\n// modified\n"), 0o644); err != nil {
		t.Fatalf("modify tracked: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "untracked.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write untracked: %v", err)
	}

	files := []string{"tracked.go", "untracked.go"}
	set := determineUpdateSet(files, false, time.Now(), dir, types.Config{RootPath: "."})
	if len(set) != 2 {
		t.Fatalf("expected both files included, got %v", set)
	}
}

func runGitCommand(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}
