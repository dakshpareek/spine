package fs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureGitignoreEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	if err := os.WriteFile(path, []byte("node_modules/\n"), 0o644); err != nil {
		t.Fatalf("failed to seed gitignore: %v", err)
	}

	if err := EnsureGitignoreEntry(path, ".code-context/"); err != nil {
		t.Fatalf("EnsureGitignoreEntry error: %v", err)
	}
	if err := EnsureGitignoreEntry(path, ".code-context/"); err != nil {
		t.Fatalf("EnsureGitignoreEntry second call error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}

	content := string(data)
	if count := strings.Count(content, ".code-context/"); count != 1 {
		t.Fatalf("expected .code-context/ entry once, found %d occurrences in %q", count, content)
	}
}
