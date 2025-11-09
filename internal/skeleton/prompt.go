package skeleton

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dakshpareek/spine/internal/types"
)

import _ "embed"

const (
	// PromptFileName is the filename used within .code-context for the prompt template.
	PromptFileName = "skeleton-prompt.txt"
)

//go:embed default_prompt.txt
var defaultPrompt string

// DefaultPrompt returns the default skeleton generation prompt template.
func DefaultPrompt() string {
	return defaultPrompt
}

// LoadPromptTemplate loads the skeleton prompt template, allowing user overrides.
// Preference order:
//  1. Config-specific override at <RootPath>/.code-context/skeleton-prompt.txt
//  2. Workspace default at .code-context/skeleton-prompt.txt
//  3. Embedded default prompt
func LoadPromptTemplate(cfg types.Config) (string, error) {
	baseDir := filepath.Dir(DirRoot) // .code-context

	candidates := []string{
		filepath.Clean(filepath.Join(baseDir, PromptFileName)),
	}

	if root := cfg.RootPath; root != "" && root != "." {
		candidates = append([]string{filepath.Clean(filepath.Join(root, baseDir, PromptFileName))}, candidates...)
	}

	for _, candidate := range uniqueStrings(candidates) {
		data, err := os.ReadFile(candidate)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("read prompt template: %w", err)
		}
		return string(data), nil
	}

	return DefaultPrompt(), nil
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	var result []string
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
