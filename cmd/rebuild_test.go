package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPurgeSkeletons(t *testing.T) {
	root := filepath.Join(t.TempDir(), "skeletons")

	paths := []string{
		filepath.Join(root, "src", "a.skeleton.go"),
		filepath.Join(root, "src", "nested", "b.skeleton.ts"),
		filepath.Join(root, "pkg", "c.skeleton.py"),
	}

	for _, p := range paths {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(p, []byte("data"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	count, err := purgeSkeletons(root)
	if err != nil {
		t.Fatalf("purgeSkeletons error: %v", err)
	}
	if count != len(paths) {
		t.Fatalf("expected %d deleted files, got %d", len(paths), count)
	}

	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("expected skeleton root to be removed, stat err=%v", err)
	}

	// Calling again on missing directory should be a no-op.
	count, err = purgeSkeletons(root)
	if err != nil {
		t.Fatalf("second purge error: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected zero deletions on second purge, got %d", count)
	}
}
