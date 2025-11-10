package index

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dakshpareek/ctx/internal/types"
)

func TestCreateEmptyIndex(t *testing.T) {
	idx := CreateEmptyIndex()

	if idx == nil {
		t.Fatalf("expected non-nil index")
	}
	if idx.Version != DefaultIndexVersion {
		t.Fatalf("expected version %q, got %q", DefaultIndexVersion, idx.Version)
	}
	if idx.PromptVersion != DefaultPromptVersion {
		t.Fatalf("expected prompt version %q, got %q", DefaultPromptVersion, idx.PromptVersion)
	}
	if idx.Files == nil {
		t.Fatalf("expected files map to be initialized")
	}
	if got := idx.Config.SkeletonPromptVersion; got != DefaultPromptVersion {
		t.Fatalf("expected skeleton prompt version %q, got %q", DefaultPromptVersion, got)
	}
	if idx.Stats.TotalFiles != 0 {
		t.Fatalf("expected zero stats for new index, got %#v", idx.Stats)
	}
}

func TestUpdateAndRemoveFileEntry(t *testing.T) {
	idx := CreateEmptyIndex()

	entry := types.FileEntry{
		Hash:         "abc123",
		SkeletonHash: "def456",
		SkeletonPath: ".ctx/skeletons/src/example.skeleton.ts",
		LastModified: time.Date(2024, 6, 5, 12, 0, 0, 0, time.UTC),
		Status:       types.StatusCurrent,
		Type:         "service",
		Size:         1024,
	}

	path := "src/example.ts"
	if err := UpdateFileEntry(idx, path, entry); err != nil {
		t.Fatalf("UpdateFileEntry error: %v", err)
	}

	stored, ok := idx.Files[path]
	if !ok {
		t.Fatalf("expected file entry for %q", path)
	}
	if stored.Path != path {
		t.Fatalf("expected stored path %q, got %q", path, stored.Path)
	}

	if idx.Stats.TotalFiles != 1 || idx.Stats.Current != 1 {
		t.Fatalf("unexpected stats after update: %#v", idx.Stats)
	}

	if err := RemoveFileEntry(idx, path); err != nil {
		t.Fatalf("RemoveFileEntry error: %v", err)
	}

	if _, exists := idx.Files[path]; exists {
		t.Fatalf("expected entry %q to be removed", path)
	}

	if idx.Stats.TotalFiles != 0 {
		t.Fatalf("expected stats to reset after removal, got %#v", idx.Stats)
	}
}

func TestSaveAndLoadIndex(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "index.json")

	idx := CreateEmptyIndex()
	idx.Version = "1.0.1"
	idx.PromptVersion = "3.0"
	idx.LastSync = time.Date(2024, 5, 1, 8, 30, 0, 0, time.UTC)

	entry1 := types.FileEntry{
		Hash:         "hash1",
		SkeletonHash: "skeleton1",
		SkeletonPath: ".ctx/skeletons/src/foo.skeleton.ts",
		LastModified: time.Date(2024, 5, 1, 8, 0, 0, 0, time.UTC),
		Status:       types.StatusStale,
		Type:         "service",
		Size:         2048,
	}
	entry2 := types.FileEntry{
		Hash:         "hash2",
		SkeletonHash: "skeleton2",
		SkeletonPath: ".ctx/skeletons/src/bar.skeleton.ts",
		LastModified: time.Date(2024, 5, 1, 8, 10, 0, 0, time.UTC),
		Status:       types.StatusPendingGeneration,
		Type:         "controller",
		Size:         4096,
	}

	if err := UpdateFileEntry(idx, "src/foo.ts", entry1); err != nil {
		t.Fatalf("UpdateFileEntry error: %v", err)
	}
	if err := UpdateFileEntry(idx, "src/bar.ts", entry2); err != nil {
		t.Fatalf("UpdateFileEntry error: %v", err)
	}

	if err := SaveIndex(idx, indexPath); err != nil {
		t.Fatalf("SaveIndex error: %v", err)
	}

	loaded, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("LoadIndex error: %v", err)
	}

	if loaded.Version != idx.Version {
		t.Fatalf("expected version %q, got %q", idx.Version, loaded.Version)
	}
	if loaded.PromptVersion != idx.PromptVersion {
		t.Fatalf("expected prompt version %q, got %q", idx.PromptVersion, loaded.PromptVersion)
	}
	if !loaded.LastSync.Equal(idx.LastSync) {
		t.Fatalf("expected last sync %v, got %v", idx.LastSync, loaded.LastSync)
	}
	if loaded.Stats.TotalFiles != 2 || loaded.Stats.Stale != 1 || loaded.Stats.PendingGeneration != 1 {
		t.Fatalf("unexpected stats after load: %#v", loaded.Stats)
	}
	if loaded.Config.SkeletonPromptVersion == "" {
		t.Fatalf("expected skeleton prompt version to be set")
	}

	gotFoo, ok := loaded.Files["src/foo.ts"]
	if !ok {
		t.Fatalf("expected foo entry to exist")
	}
	if gotFoo.Status != types.StatusStale {
		t.Fatalf("expected foo status stale, got %s", gotFoo.Status)
	}
	if !gotFoo.LastModified.Equal(entry1.LastModified) {
		t.Fatalf("expected foo LastModified %v, got %v", entry1.LastModified, gotFoo.LastModified)
	}
}

func TestLoadIndexInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "index.json")

	if err := os.WriteFile(indexPath, []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	if _, err := LoadIndex(indexPath); err == nil {
		t.Fatalf("expected error when loading invalid json")
	}
}

func TestLoadIndexMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "missing.json")

	_, err := LoadIndex(indexPath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected not exist error, got %v", err)
	}
}

func TestNilIndexErrors(t *testing.T) {
	if err := SaveIndex(nil, "some/path"); !errors.Is(err, ErrNilIndex) {
		t.Fatalf("expected ErrNilIndex from SaveIndex, got %v", err)
	}

	if err := UpdateFileEntry(nil, "path", types.FileEntry{}); !errors.Is(err, ErrNilIndex) {
		t.Fatalf("expected ErrNilIndex from UpdateFileEntry, got %v", err)
	}

	if err := RemoveFileEntry(nil, "path"); !errors.Is(err, ErrNilIndex) {
		t.Fatalf("expected ErrNilIndex from RemoveFileEntry, got %v", err)
	}

	if stats := CalculateStats(nil); stats.TotalFiles != 0 {
		t.Fatalf("expected zero stats for nil index, got %#v", stats)
	}
}
