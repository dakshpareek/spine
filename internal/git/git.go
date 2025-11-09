package git

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	// ErrNotGit indicates commands that require git were executed outside a repository.
	ErrNotGit = errors.New("not a git repository")
)

type commandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

type execCommandRunner struct{}

func (execCommandRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

var runner commandRunner = execCommandRunner{}

// IsGitRepo reports whether the current directory resides inside a git repository.
func IsGitRepo() bool {
	output, err := runGitCommand("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// GetModifiedFiles returns tracked files modified since HEAD. Requires git repository.
func GetModifiedFiles() ([]string, error) {
	if !IsGitRepo() {
		return nil, ErrNotGit
	}

	output, err := runGitCommand("diff", "--name-only", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("git diff --name-only: %w", err)
	}

	return parseGitList(output), nil
}

// GetUntrackedFiles returns files unknown to git (respecting exclude patterns). Requires git repository.
func GetUntrackedFiles() ([]string, error) {
	if !IsGitRepo() {
		return nil, ErrNotGit
	}

	output, err := runGitCommand("ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, fmt.Errorf("git ls-files --others: %w", err)
	}

	return parseGitList(output), nil
}

// GetModifiedFilesFallback walks the filesystem and returns files modified after the provided timestamp.
func GetModifiedFilesFallback(root string, since time.Time) ([]string, error) {
	if root == "" {
		root = "."
	}
	root = filepath.Clean(root)

	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().After(since) {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func runGitCommand(args ...string) ([]byte, error) {
	return runner.Run("git", args...)
}

func parseGitList(output []byte) []string {
	var files []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		files = append(files, line)
	}
	return files
}
