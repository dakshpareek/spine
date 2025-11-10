package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	indexpkg "github.com/dakshpareek/ctx/internal/index"
	"github.com/dakshpareek/ctx/internal/types"
)

func changeDir(t *testing.T, dir string) func() {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	return func() {
		_ = os.Chdir(original)
	}
}

func executeCommand(t *testing.T, dir string, args ...string) (string, string) {
	t.Helper()
	cleanup := changeDir(t, dir)
	defer cleanup()

	cmd := NewRootCmd("test")
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("command %v failed: %v\nstdout:\n%s\nstderr:\n%s", args, err, stdout.String(), stderr.String())
	}

	return stdout.String(), stderr.String()
}

func writeTempFile(t *testing.T, dir, rel, content string) string {
	t.Helper()
	target := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return target
}

func loadIndex(t *testing.T, dir string) *types.Index {
	t.Helper()
	idx, err := indexpkg.LoadIndex(filepath.Join(dir, ".ctx", "index.json"))
	if err != nil {
		t.Fatalf("load index: %v", err)
	}
	return idx
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = oldStdout
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(data)
}
