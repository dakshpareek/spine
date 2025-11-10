package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dakshpareek/ctx/internal/types"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	expect := &types.Config{
		IncludedExtensions:    defaultIncludedExtensions,
		ExcludedPaths:         defaultExcludedPaths,
		SkeletonPromptVersion: DefaultSkeletonPromptVersion,
		RootPath:              DefaultRootPath,
	}

	if !reflect.DeepEqual(cfg.IncludedExtensions, expect.IncludedExtensions) {
		t.Fatalf("included extensions mismatch:\nexpected: %#v\nactual:   %#v", expect.IncludedExtensions, cfg.IncludedExtensions)
	}

	if !reflect.DeepEqual(cfg.ExcludedPaths, expect.ExcludedPaths) {
		t.Fatalf("excluded paths mismatch:\nexpected: %#v\nactual:   %#v", expect.ExcludedPaths, cfg.ExcludedPaths)
	}

	if cfg.SkeletonPromptVersion != expect.SkeletonPromptVersion {
		t.Fatalf("expected skeleton prompt version %q, got %q", expect.SkeletonPromptVersion, cfg.SkeletonPromptVersion)
	}

	if cfg.RootPath != expect.RootPath {
		t.Fatalf("expected root path %q, got %q", expect.RootPath, cfg.RootPath)
	}

	// Ensure slices are cloned (mutating result should not alter defaults).
	cfg.IncludedExtensions[0] = ".changed"
	another := GetDefaultConfig()
	if another.IncludedExtensions[0] == ".changed" {
		t.Fatalf("default included extensions mutated after modification")
	}
}

func TestLoadConfigMissingFileReturnsDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	expect := GetDefaultConfig()
	if !reflect.DeepEqual(cfg, expect) {
		t.Fatalf("expected default config %#v, got %#v", expect, cfg)
	}
}

func TestLoadConfigPartialOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{
  "includedExtensions": [".rb"],
  "skeletonPromptVersion": "3.0"
}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if !reflect.DeepEqual(cfg.IncludedExtensions, []string{".rb"}) {
		t.Fatalf("expected included extensions to persist override, got %#v", cfg.IncludedExtensions)
	}

	if !reflect.DeepEqual(cfg.ExcludedPaths, defaultExcludedPaths) {
		t.Fatalf("expected excluded paths defaults, got %#v", cfg.ExcludedPaths)
	}

	if cfg.SkeletonPromptVersion != "3.0" {
		t.Fatalf("expected skeleton prompt version 3.0, got %q", cfg.SkeletonPromptVersion)
	}

	if cfg.RootPath != DefaultRootPath {
		t.Fatalf("expected default root path %q, got %q", DefaultRootPath, cfg.RootPath)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := os.WriteFile(path, []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	if _, err := LoadConfig(path); err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}
