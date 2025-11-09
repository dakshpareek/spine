package scanner

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/dakshpareek/spine/internal/types"
)

var fileTypePatterns = map[string][]string{
	"service": {
		"*service.ts",
		"*service.js",
		"*service.go",
	},
	"controller": {
		"*controller.ts",
		"*controller.js",
		"*handler.go",
	},
	"repository": {
		"*repository.ts",
		"*repo.ts",
		"*repository.go",
	},
	"dto": {
		"*dto.ts",
		"*dto.go",
		"*/dto/*",
	},
	"model": {
		"*model.ts",
		"*entity.ts",
		"*model.go",
	},
	"util": {
		"*util.ts",
		"*utils.ts",
		"*helper.ts",
	},
	"middleware": {
		"*middleware.ts",
		"*middleware.go",
	},
	"config": {
		"*config.ts",
		"*config.go",
	},
}

// ScanFiles walks the configured root directory and returns matching files.
func ScanFiles(cfg types.Config) ([]string, error) {
	root := cfg.RootPath
	if root == "" {
		root = "."
	}
	root = filepath.Clean(root)

	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}
		if rel == "." {
			return nil
		}

		normalized := filepath.ToSlash(rel)

		excluded, err := isExcluded(normalized, cfg.ExcludedPaths)
		if err != nil {
			return err
		}

		if excluded {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if !shouldIncludeByExtension(normalized, cfg.IncludedExtensions) {
			return nil
		}

		files = append(files, normalized)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

// DetectFileType returns the inferred file classification based on known patterns.
func DetectFileType(p string) string {
	normalized := filepath.ToSlash(p)
	base := path.Base(normalized)
	for fileType, patterns := range fileTypePatterns {
		for _, pattern := range patterns {
			matched, err := doublestar.Match(pattern, normalized)
			if err != nil {
				continue
			}
			if matched {
				return fileType
			}

			matchedBase, err := doublestar.Match(pattern, base)
			if err != nil {
				continue
			}
			if matchedBase {
				return fileType
			}
		}
	}
	return ""
}

func isExcluded(path string, patterns []string) (bool, error) {
	if len(patterns) == 0 {
		return false, nil
	}

	segments := strings.Split(path, "/")
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}

		matched, err := doublestar.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}

		for _, segment := range segments {
			if segment == "" {
				continue
			}

			matched, err := doublestar.Match(pattern, segment)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
	}

	return false, nil
}

func shouldIncludeByExtension(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}

	for _, allowed := range extensions {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}
