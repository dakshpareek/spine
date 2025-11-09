package hash

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestHashContent(t *testing.T) {
	hash := HashContent([]byte("hello world"))
	const expected = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	if hash != expected {
		t.Fatalf("expected hash %s, got %s", expected, hash)
	}
}

func TestHashFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := os.WriteFile(path, []byte("hello world"), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	hash, err := HashFile(path)
	if err != nil {
		t.Fatalf("HashFile error: %v", err)
	}

	if expected := HashContent([]byte("hello world")); hash != expected {
		t.Fatalf("expected hash %s, got %s", expected, hash)
	}
}

func TestHashFileLargeContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.bin")

	content := bytes.Repeat([]byte("ctx"), 512*1024) // ~1.5MB
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("failed to write large file: %v", err)
	}

	hashFromFile, err := HashFile(path)
	if err != nil {
		t.Fatalf("HashFile error: %v", err)
	}

	hashFromMemory := HashContent(content)
	if hashFromFile != hashFromMemory {
		t.Fatalf("hash mismatch for large content")
	}
}

func TestHashFileMissing(t *testing.T) {
	if _, err := HashFile("does-not-exist"); err == nil {
		t.Fatalf("expected error for missing file")
	}
}
