package scanner

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dakshpareek/spine/internal/types"
)

func TestScanFilesRespectsConfig(t *testing.T) {
	root := t.TempDir()

	makeFile(t, root, "src/service/user.service.ts")
	makeFile(t, root, "src/controller/user.controller.ts")
	makeFile(t, root, "src/repository/user.repository.go")
	makeFile(t, root, "src/dto/user.dto.ts")
	makeFile(t, root, "src/util/math.util.ts")
	makeFile(t, root, "scripts/build.go")

	// Excluded files/directories
	makeFile(t, root, "src/foo.spec.ts")
	makeFile(t, root, "node_modules/library/index.ts")
	makeFile(t, root, "test/helper.ts")

	cfg := types.Config{
		IncludedExtensions: []string{".ts", ".go"},
		ExcludedPaths:      []string{"node_modules", "*.spec.*", "*.test.*", "test"},
		RootPath:           root,
	}

	files, err := ScanFiles(cfg)
	if err != nil {
		t.Fatalf("ScanFiles error: %v", err)
	}

	expected := []string{
		"scripts/build.go",
		"src/controller/user.controller.ts",
		"src/dto/user.dto.ts",
		"src/repository/user.repository.go",
		"src/service/user.service.ts",
		"src/util/math.util.ts",
	}

	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("expected files %#v, got %#v", expected, files)
	}
}

func TestScanFilesNonexistentRoot(t *testing.T) {
	cfg := types.Config{
		IncludedExtensions: []string{".go"},
		RootPath:           filepath.Join(t.TempDir(), "missing"),
	}

	if _, err := ScanFiles(cfg); err == nil {
		t.Fatalf("expected error for missing root path")
	}
}

func TestDetectFileType(t *testing.T) {
	tests := map[string]string{
		"src/app/user.service.ts":               "service",
		"src/app/user.controller.ts":            "controller",
		"src/app/user.handler.go":               "controller",
		"src/app/user.repository.go":            "repository",
		"src/app/dto/user.dto.ts":               "dto",
		"src/app/models/user.model.ts":          "model",
		"src/app/util/math.util.ts":             "util",
		"src/app/middleware/auth.middleware.go": "middleware",
		"config/app.config.go":                  "config",
		"src/app/unknown/file.txt":              "",
	}

	for path, expected := range tests {
		if got := DetectFileType(path); got != expected {
			t.Fatalf("DetectFileType(%q) = %q, expected %q", path, got, expected)
		}
	}
}

func makeFile(t *testing.T, root, rel string) {
	t.Helper()

	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directories for %q: %v", path, err)
	}
	if err := os.WriteFile(path, []byte("sample"), 0o644); err != nil {
		t.Fatalf("failed to write file %q: %v", path, err)
	}
}
