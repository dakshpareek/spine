package fs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDirAndWriteFile(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, "nested")
	if err := EnsureDir(targetDir); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}
	if !Exists(targetDir) {
		t.Fatalf("expected directory to exist")
	}

	file := filepath.Join(targetDir, "file.txt")
	if err := WriteFile(file, []byte("hello")); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil || string(data) != "hello" {
		t.Fatalf("unexpected file contents")
	}
}

func TestWriteJSONAndGitignore(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "data.json")
	obj := map[string]string{"a": "b"}
	if err := WriteJSON(target, obj); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read json: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(data, &decoded); err != nil || decoded["a"] != "b" {
		t.Fatalf("unexpected json content")
	}

	gitignore := filepath.Join(dir, ".gitignore")
	if err := EnsureGitignoreEntry(gitignore, ".code-context/"); err != nil {
		t.Fatalf("EnsureGitignoreEntry: %v", err)
	}
	if err := EnsureGitignoreEntry(gitignore, ".code-context/"); err != nil {
		t.Fatalf("EnsureGitignoreEntry duplicate: %v", err)
	}
	data, err = os.ReadFile(gitignore)
	if err != nil {
		t.Fatalf("read gitignore: %v", err)
	}
	if string(data) != ".code-context/\n" {
		t.Fatalf("unexpected gitignore content: %s", string(data))
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	if Exists(filepath.Join(dir, "missing")) {
		t.Fatalf("expected missing path to return false")
	}
	if err := os.WriteFile(filepath.Join(dir, "file"), []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if !Exists(filepath.Join(dir, "file")) {
		t.Fatalf("expected file to exist")
	}
}

func TestEnsureDirErrors(t *testing.T) {
	if err := EnsureDir(""); err == nil {
		t.Fatalf("expected error for empty path")
	}
}

func TestWriteFileErrors(t *testing.T) {
	if err := WriteFile("", []byte("")); err == nil {
		t.Fatalf("expected error for empty path")
	}
}

func TestWriteJSONErrors(t *testing.T) {
	if err := WriteJSON("", map[string]string{}); err == nil {
		t.Fatalf("expected error for empty path")
	}
}
