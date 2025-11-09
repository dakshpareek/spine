package fs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultDirPerm  = 0o755
	defaultFilePerm = 0o644
)

// EnsureDir creates the directory hierarchy if it does not already exist.
func EnsureDir(path string) error {
	if path == "" {
		return errors.New("directory path is empty")
	}
	if err := os.MkdirAll(path, defaultDirPerm); err != nil {
		return fmt.Errorf("create directory %q: %w", path, err)
	}
	return nil
}

// WriteFile writes the provided data to the path, creating parent directories when needed.
func WriteFile(path string, data []byte) error {
	if path == "" {
		return errors.New("file path is empty")
	}
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, defaultFilePerm); err != nil {
		return fmt.Errorf("write file %q: %w", path, err)
	}
	return nil
}

// WriteJSON writes the provided value as pretty-formatted JSON to the given path.
func WriteJSON(path string, value interface{}) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode json for %q: %w", path, err)
	}
	data = append(data, '\n')
	return WriteFile(path, data)
}

// EnsureGitignoreEntry adds the given entry to the .gitignore if it is not already present.
func EnsureGitignoreEntry(path string, entry string) error {
	if entry == "" {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, defaultFilePerm)
	if err != nil {
		return fmt.Errorf("open .gitignore: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == entry {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read .gitignore: %w", err)
	}

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("seek .gitignore: %w", err)
	}

	if _, err := file.WriteString(entry + "\n"); err != nil {
		return fmt.Errorf("append to .gitignore: %w", err)
	}

	return nil
}

// Exists reports whether the path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
