package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPipelineGeneratesPrompt(t *testing.T) {
	tempDir := t.TempDir()
	writeTempFile(t, tempDir, "main.go", "package main\n")

	_, _ = executeCommand(t, tempDir, "init")

	if err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte("package main\n// change\n"), 0o644); err != nil {
		t.Fatalf("write change: %v", err)
	}

	_, _ = executeCommand(t, tempDir, "pipeline", "--output", "pipeline.md")

	data, err := os.ReadFile(filepath.Join(tempDir, "pipeline.md"))
	if err != nil {
		t.Fatalf("read pipeline prompt: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected pipeline prompt to have content")
	}
}
