package cmd

import (
	"strings"
	"testing"
)

func TestRunGenerateErrors(t *testing.T) {
	dir := t.TempDir()
	cleanup := changeDir(t, dir)
	err := runGenerate(generateOptions{})
	if err == nil {
		t.Fatalf("expected error when not initialized")
	}
	cleanup()

	_, _ = executeCommand(t, dir, "init")

	cleanup = changeDir(t, dir)
	defer cleanup()

	if err := runGenerate(generateOptions{files: "unknown.go"}); err == nil {
		t.Fatalf("expected error for untracked file")
	}
}

func TestRunGenerateStdout(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "main.go", "package main\n")
	_, _ = executeCommand(t, dir, "init")

	cleanup := changeDir(t, dir)
	defer cleanup()

	output := captureOutput(t, func() {
		if err := runGenerate(generateOptions{}); err != nil {
			t.Fatalf("runGenerate: %v", err)
		}
	})
	if !strings.Contains(output, "Code Context Skeleton Generation") {
		t.Fatalf("expected prompt output")
	}
}
