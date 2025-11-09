package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	rootcmd "github.com/dakshpareek/spine/cmd"
	"github.com/dakshpareek/spine/internal/index"
)

func TestInitSyncGenerateFlow(t *testing.T) {
	tempDir := t.TempDir()

	srcDir := filepath.Join("testdata", "sample-project")
	if err := copyTree(srcDir, tempDir); err != nil {
		t.Fatalf("copyTree: %v", err)
	}

	withWorkingDir(t, tempDir)

	runCommand(t, "init")
	runCommand(t, "sync")
	runCommand(t, "generate", "--output", "prompt.md")

	if _, err := os.Stat(".spine/index.json"); err != nil {
		t.Fatalf("expected index.json to exist: %v", err)
	}

	if _, err := os.Stat("prompt.md"); err != nil {
		t.Fatalf("expected prompt.md to be generated: %v", err)
	}

	loaded, err := index.LoadIndex(filepath.Join(".spine", "index.json"))
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}

	var pending int
	for _, entry := range loaded.Files {
		if entry.Status == "pendingGeneration" {
			pending++
		}
	}

	if pending == 0 {
		t.Fatalf("expected at least one file marked pendingGeneration after generate")
	}

	verifyPromptContains(t, "prompt.md", "## Files to Process")
}

func runCommand(t *testing.T, args ...string) {
	t.Helper()

	cmd := rootcmd.NewRootCmd("test")
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command %v failed: %v\nstdout:\n%s\nstderr:\n%s", args, err, stdout.String(), stderr.String())
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(original)
	})
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		return nil
	})
}

func verifyPromptContains(t *testing.T, path, substring string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read prompt: %v", err)
	}
	if !bytes.Contains(data, []byte(substring)) {
		t.Fatalf("expected prompt to contain %q", substring)
	}
}
