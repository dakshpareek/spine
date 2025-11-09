package types

import "time"

// Status represents the synchronization state of a file's skeleton.
type Status string

const (
	StatusCurrent           Status = "current"
	StatusStale             Status = "stale"
	StatusMissing           Status = "missing"
	StatusPendingGeneration Status = "pendingGeneration"
)

// ExitCode represents the exit status of the CLI.
type ExitCode int

const (
	ExitCodeOK         ExitCode = 0
	ExitCodeUserError  ExitCode = 1
	ExitCodeFileSystem ExitCode = 2
	ExitCodeGit        ExitCode = 3
	ExitCodeData       ExitCode = 4
)

// Error wraps an error with an associated exit code.
type Error struct {
	Code ExitCode
	Err  error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

// Unwrap supports errors.Unwrap.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// FileEntry describes a single tracked file within the index.
type FileEntry struct {
	Path         string    `json:"path"`
	Hash         string    `json:"hash"`
	SkeletonHash string    `json:"skeletonHash"`
	SkeletonPath string    `json:"skeletonPath"`
	LastModified time.Time `json:"lastModified"`
	Status       Status    `json:"status"`
	Type         string    `json:"type"`
	Size         int64     `json:"size"`
}

// IndexStats aggregates counts of files by status.
type IndexStats struct {
	TotalFiles        int `json:"totalFiles"`
	Current           int `json:"current"`
	Stale             int `json:"stale"`
	Missing           int `json:"missing"`
	PendingGeneration int `json:"pendingGeneration"`
}

// Index is the root structure persisted as index.json.
type Index struct {
	Version       string               `json:"version"`
	PromptVersion string               `json:"promptVersion"`
	LastSync      time.Time            `json:"lastSync"`
	Config        Config               `json:"config"`
	Files         map[string]FileEntry `json:"files"`
	Stats         IndexStats           `json:"stats"`
}

// Config captures user configuration for scanning behavior.
type Config struct {
	IncludedExtensions    []string `json:"includedExtensions"`
	ExcludedPaths         []string `json:"excludedPaths"`
	SkeletonPromptVersion string   `json:"skeletonPromptVersion"`
	RootPath              string   `json:"rootPath"`
}
