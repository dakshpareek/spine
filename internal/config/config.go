package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dakshpareek/spine/internal/types"
)

const (
	// DefaultRootPath is the default relative root used during scanning.
	DefaultRootPath = "."
	// DefaultSkeletonPromptVersion matches the bundled prompt template version.
	DefaultSkeletonPromptVersion = "2.1"
)

var (
	defaultIncludedExtensions = []string{
		".ts",
		".tsx",
		".js",
		".jsx",
		".go",
		".py",
	}

	defaultExcludedPaths = []string{
		"node_modules",
		"vendor",
		"dist",
		"build",
		".next",
		"coverage",
		".spine",
		".git",
		"*.test.*",
		"*.spec.*",
		"__tests__",
		"test",
	}
)

// GetDefaultConfig returns a new Config populated with PRD defaults.
func GetDefaultConfig() *types.Config {
	return &types.Config{
		IncludedExtensions:    cloneSlice(defaultIncludedExtensions),
		ExcludedPaths:         cloneSlice(defaultExcludedPaths),
		SkeletonPromptVersion: DefaultSkeletonPromptVersion,
		RootPath:              DefaultRootPath,
	}
}

// LoadConfig reads a configuration file if present, otherwise returning defaults.
func LoadConfig(path string) (*types.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return GetDefaultConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg types.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	applyDefaults(&cfg)
	return &cfg, nil
}

func applyDefaults(cfg *types.Config) {
	if cfg == nil {
		return
	}

	if cfg.IncludedExtensions == nil {
		cfg.IncludedExtensions = cloneSlice(defaultIncludedExtensions)
	}

	if cfg.ExcludedPaths == nil {
		cfg.ExcludedPaths = cloneSlice(defaultExcludedPaths)
	}

	if cfg.SkeletonPromptVersion == "" {
		cfg.SkeletonPromptVersion = DefaultSkeletonPromptVersion
	}

	if cfg.RootPath == "" {
		cfg.RootPath = DefaultRootPath
	}
}

func cloneSlice(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	copied := make([]string, len(values))
	copy(copied, values)
	return copied
}
