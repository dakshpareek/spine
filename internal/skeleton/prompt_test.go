package skeleton

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dakshpareek/spine/internal/types"
)

func TestLoadPromptTemplateDefault(t *testing.T) {
	tempDir := t.TempDir()
	withWorkingDir(t, tempDir)

	cfg := types.Config{}
	prompt, err := LoadPromptTemplate(cfg)
	if err != nil {
		t.Fatalf("LoadPromptTemplate error: %v", err)
	}
	if prompt != DefaultPrompt() {
		t.Fatalf("expected default prompt, got %q", prompt)
	}
}

func TestLoadPromptTemplateWorkspaceOverride(t *testing.T) {
	tempDir := t.TempDir()
	withWorkingDir(t, tempDir)

	if err := os.MkdirAll(filepath.Join(".code-context"), 0o755); err != nil {
		t.Fatalf("failed to create workspace dir: %v", err)
	}
	expected := "custom workspace prompt"
	if err := os.WriteFile(filepath.Join(".code-context", PromptFileName), []byte(expected), 0o644); err != nil {
		t.Fatalf("failed to write prompt override: %v", err)
	}

	prompt, err := LoadPromptTemplate(types.Config{})
	if err != nil {
		t.Fatalf("LoadPromptTemplate error: %v", err)
	}
	if prompt != expected {
		t.Fatalf("expected workspace prompt %q, got %q", expected, prompt)
	}
}

func TestLoadPromptTemplateRootOverride(t *testing.T) {
	tempDir := t.TempDir()
	withWorkingDir(t, tempDir)

	rootPath := filepath.Join(tempDir, "app")
	if err := os.MkdirAll(filepath.Join(rootPath, ".code-context"), 0o755); err != nil {
		t.Fatalf("failed to create root workspace: %v", err)
	}

	expected := "root override prompt"
	target := filepath.Join(rootPath, ".code-context", PromptFileName)
	if err := os.WriteFile(target, []byte(expected), 0o644); err != nil {
		t.Fatalf("failed to write root override: %v", err)
	}

	prompt, err := LoadPromptTemplate(types.Config{RootPath: rootPath})
	if err != nil {
		t.Fatalf("LoadPromptTemplate error: %v", err)
	}
	if prompt != expected {
		t.Fatalf("expected root override %q, got %q", expected, prompt)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(orig)
	})
}
