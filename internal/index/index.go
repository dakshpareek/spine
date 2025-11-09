package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dakshpareek/spine/internal/types"
)

const (
	// DefaultIndexVersion represents the initial version written to a new index.
	DefaultIndexVersion = "1.0.0"
	// DefaultPromptVersion represents the default prompt template version.
	DefaultPromptVersion = "2.1"
)

var (
	// ErrNilIndex is returned when operations receive a nil index reference.
	ErrNilIndex = errors.New("index is nil")
)

// CreateEmptyIndex constructs an index pre-populated with default values and empty collections.
func CreateEmptyIndex() *types.Index {
	return &types.Index{
		Version:       DefaultIndexVersion,
		PromptVersion: DefaultPromptVersion,
		LastSync:      time.Time{},
		Config: types.Config{
			IncludedExtensions:    nil,
			ExcludedPaths:         nil,
			SkeletonPromptVersion: DefaultPromptVersion,
		},
		Files: make(map[string]types.FileEntry),
		Stats: types.IndexStats{},
	}
}

// LoadIndex reads an index file from disk and unmarshals it into memory.
func LoadIndex(path string) (*types.Index, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open index: %w", err)
	}
	defer file.Close()

	var idx types.Index
	if err := json.NewDecoder(file).Decode(&idx); err != nil {
		return nil, fmt.Errorf("decode index: %w", err)
	}

	ensureIndexInitialized(&idx)
	idx.Stats = CalculateStats(&idx)

	return &idx, nil
}

// SaveIndex writes the provided index structure to disk with stable formatting.
func SaveIndex(idx *types.Index, path string) error {
	if idx == nil {
		return ErrNilIndex
	}

	ensureIndexInitialized(idx)
	idx.Stats = CalculateStats(idx)

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("encode index: %w", err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create index directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write index: %w", err)
	}

	return nil
}

// UpdateFileEntry upserts a file entry within the index.
func UpdateFileEntry(idx *types.Index, path string, entry types.FileEntry) error {
	if idx == nil {
		return ErrNilIndex
	}

	ensureIndexInitialized(idx)
	entry.Path = path
	idx.Files[path] = entry
	idx.Stats = CalculateStats(idx)

	return nil
}

// RemoveFileEntry deletes a file entry from the index if it exists.
func RemoveFileEntry(idx *types.Index, path string) error {
	if idx == nil {
		return ErrNilIndex
	}

	if idx.Files == nil {
		return nil
	}

	delete(idx.Files, path)
	idx.Stats = CalculateStats(idx)

	return nil
}

// CalculateStats recomputes index statistics based on the current file set.
func CalculateStats(idx *types.Index) types.IndexStats {
	if idx == nil {
		return types.IndexStats{}
	}

	var stats types.IndexStats
	for _, entry := range idx.Files {
		stats.TotalFiles++
		switch entry.Status {
		case types.StatusCurrent:
			stats.Current++
		case types.StatusStale:
			stats.Stale++
		case types.StatusMissing:
			stats.Missing++
		case types.StatusPendingGeneration:
			stats.PendingGeneration++
		}
	}

	idx.Stats = stats
	return stats
}

func ensureIndexInitialized(idx *types.Index) {
	if idx.Files == nil {
		idx.Files = make(map[string]types.FileEntry)
	}

	if idx.Config.SkeletonPromptVersion == "" {
		if idx.PromptVersion != "" {
			idx.Config.SkeletonPromptVersion = idx.PromptVersion
		} else {
			idx.Config.SkeletonPromptVersion = DefaultPromptVersion
		}
	}

	if idx.Version == "" {
		idx.Version = DefaultIndexVersion
	}

	if idx.PromptVersion == "" {
		idx.PromptVersion = DefaultPromptVersion
	}

	if idx.LastSync.IsZero() {
		idx.LastSync = time.Time{}
	} else {
		idx.LastSync = idx.LastSync.UTC()
	}
}
