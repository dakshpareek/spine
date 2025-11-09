# Tasks: Code Context CLI - Phase 1

## Relevant Files

- `main.go` - Entry point, CLI command routing using cobra/cli library
- `cmd/init.go` - Implementation of `ctx init` command
- `cmd/sync.go` - Implementation of `spine sync` command
- `cmd/generate.go` - Implementation of `ctx generate` command
- `cmd/status.go` - Implementation of `ctx status` command
- `cmd/validate.go` - Implementation of `ctx validate` command
- `cmd/clean.go` - Implementation of `ctx clean` command
- `cmd/rebuild.go` - Implementation of `ctx rebuild` command
- `cmd/export.go` - Implementation of `ctx export` command
- `internal/index/index.go` - Index data structures and operations (load, save, update)
- `internal/index/index_test.go` - Unit tests for index operations
- `internal/scanner/scanner.go` - File scanning with exclusions/inclusions
- `internal/scanner/scanner_test.go` - Unit tests for scanner
- `internal/hash/hash.go` - File hashing utilities (SHA-256)
- `internal/hash/hash_test.go` - Unit tests for hashing
- `internal/git/git.go` - Git integration (diff, untracked files detection)
- `internal/git/git_test.go` - Unit tests for git operations
- `internal/config/config.go` - Configuration loading and defaults
- `internal/config/config_test.go` - Unit tests for config
- `internal/skeleton/prompt.go` - Skeleton prompt template and generation
- `internal/types/types.go` - Shared type definitions (FileEntry, Index, Status enum)
- `internal/fs/fs.go` - File system utilities (create directories, write files safely)
- `internal/display/display.go` - Output formatting and colored terminal display
- `.spine/skeleton-prompt.txt` - Embedded default skeleton prompt
- `go.mod` - Go module dependencies
- `go.sum` - Dependency checksums
- `README.md` - User documentation
- `docs/commands.md` - Detailed command reference
- `docs/getting-started.md` - Step-by-step tutorial
- `.goreleaser.yml` - Multi-platform binary build configuration
- `.github/workflows/release.yml` - GitHub Actions for automated releases
- `.github/workflows/test.yml` - GitHub Actions for CI testing

### Notes

- Go testing convention: place `*_test.go` files alongside the code they test
- Run tests: `go test ./...` (all packages) or `go test ./internal/index` (specific package)
- Use `//go:embed` directive to embed skeleton-prompt.txt into binary
- Follow standard Go project layout: `cmd/` for commands, `internal/` for private packages

## Tasks

- [x] 1.0 Project Setup & Core Infrastructure
  - [x] 1.1 Initialize Go module with `go mod init github.com/yourusername/code-context`
  - [x] 1.2 Install CLI framework dependency (`github.com/spf13/cobra` for command routing)
  - [x] 1.3 Install additional dependencies: `github.com/bmatcuk/doublestar/v4` (glob matching), `github.com/fatih/color` (terminal colors)
  - [x] 1.4 Create project directory structure: `cmd/`, `internal/`, `docs/`, `tasks/`
  - [x] 1.5 Set up `.gitignore` with Go-specific entries (vendor/, *.exe, *.test, etc.)
  - [x] 1.6 Create `main.go` entry point that initializes cobra root command
  - [x] 1.7 Set up basic error handling pattern and exit codes (1-4 as per PRD)
  - [x] 1.8 Create `internal/types/types.go` with core data structures: `Index`, `FileEntry`, `Config`, `Status` enum

- [x] 2.0 Index Management System
  - [x] 2.1 Implement `internal/index/index.go` with `Index` struct matching PRD schema
  - [x] 2.2 Create `LoadIndex(path string) (*Index, error)` function to read and parse JSON
  - [x] 2.3 Create `SaveIndex(index *Index, path string) error` function to write JSON with proper formatting
  - [x] 2.4 Implement `CreateEmptyIndex() *Index` to generate initial index structure with defaults
  - [x] 2.5 Add `UpdateFileEntry(index *Index, path string, entry FileEntry)` method
  - [x] 2.6 Add `RemoveFileEntry(index *Index, path string)` method for deleted files
  - [x] 2.7 Implement `CalculateStats(index *Index)` to compute current/stale/missing counts
  - [x] 2.8 Write unit tests in `internal/index/index_test.go` covering load/save/update operations
  - [x] 2.9 Test JSON serialization edge cases (empty index, corrupted JSON, missing fields)

- [x] 3.0 File Scanning & Change Detection
  - [x] 3.1 Create `internal/config/config.go` with `Config` struct matching PRD schema
  - [x] 3.2 Implement `LoadConfig(path string) (*Config, error)` with default fallback
  - [x] 3.3 Implement `GetDefaultConfig() *Config` returning PRD-specified defaults
  - [x] 3.4 Write config tests in `internal/config/config_test.go`
  - [x] 3.5 Create `internal/scanner/scanner.go` with `ScanFiles(config Config) ([]string, error)`
  - [x] 3.6 Implement filepath walking with exclusion pattern matching using doublestar
  - [x] 3.7 Add extension filtering logic for `includedExtensions`
  - [x] 3.8 Implement file type detection based on path patterns (service, controller, repository, etc.)
  - [x] 3.9 Write scanner tests covering exclusions, inclusions, and edge cases
  - [x] 3.10 Create `internal/hash/hash.go` with `HashFile(path string) (string, error)` using SHA-256
  - [x] 3.11 Add `HashContent(content []byte) string` helper for in-memory hashing
  - [x] 3.12 Write hash tests verifying correctness and handling of large files
  - [x] 3.13 Create `internal/git/git.go` with `IsGitRepo() bool` check
  - [x] 3.14 Implement `GetModifiedFiles() ([]string, error)` using `git diff --name-only HEAD`
  - [x] 3.15 Implement `GetUntrackedFiles() ([]string, error)` using `git ls-files --others --exclude-standard`
  - [x] 3.16 Add fallback logic for non-git repos (use file mtime comparison)
  - [x] 3.17 Write git tests with mocked exec commands

- [x] 4.0 Command Implementations (init, sync, status, validate, clean, rebuild)
  - [x] 4.1 Create `cmd/init.go` implementing `ctx init` command
  - [x] 4.2 In init: Check if `.spine/` exists, error if already initialized
  - [x] 4.3 In init: Create `.spine/`, `skeletons/`, write default `config.json` and `skeleton-prompt.txt`
  - [x] 4.4 In init: Add `.spine/` to `.gitignore` (create if missing)
  - [x] 4.5 In init: Run initial scan and create `index.json` with all files marked as "missing"
  - [x] 4.6 In init: Display success message with next steps as per PRD
  - [x] 4.7 Create `cmd/sync.go` implementing `spine sync` command
  - [x] 4.8 In sync: Load existing index, detect changes via git or mtime
  - [x] 4.9 In sync: For each file, calculate hash and update status (current/stale/missing)
  - [x] 4.10 In sync: Remove deleted files from index
  - [x] 4.11 In sync: Save updated index and display change summary
  - [x] 4.12 In sync: Add `--full` flag to force full scan (ignore git diff)
  - [x] 4.13 In sync: Add `--verbose` flag for detailed file-by-file output
  - [x] 4.14 Create `cmd/status.go` implementing `ctx status` command
  - [x] 4.15 In status: Load index, calculate stats, display formatted summary
  - [x] 4.16 In status: Add `--verbose` flag to list stale/missing files
  - [x] 4.17 In status: Add `--json` flag for machine-readable output
  - [x] 4.18 Create `internal/display/display.go` with colored output helpers (✓, ⚠, ⏳ symbols)
  - [x] 4.19 Create `cmd/validate.go` implementing `ctx validate` command
  - [x] 4.20 In validate: Check source file existence, skeleton file existence, hash mismatches
  - [x] 4.21 In validate: Add `--fix` flag to auto-mark mismatched files as stale
  - [x] 4.22 In validate: Add `--strict` flag to exit with error code if issues found
  - [x] 4.23 Create `cmd/clean.go` implementing `ctx clean` command
  - [x] 4.24 In clean: Scan `.spine/skeletons/` and remove orphaned files not in index
  - [x] 4.25 In clean: Remove empty directories after cleanup
  - [x] 4.26 Create `cmd/rebuild.go` implementing `ctx rebuild` command
  - [x] 4.27 In rebuild: Require `--confirm` flag to proceed (prevent accidental data loss)
  - [x] 4.28 In rebuild: Delete all skeleton files, reset index (preserve config), run full sync

- [x] 5.0 Skeleton Generation & Export
  - [x] 5.1 Create `internal/skeleton/prompt.go` with embedded skeleton prompt template
  - [x] 5.2 Use `//go:embed skeleton-prompt.txt` to embed default prompt in binary
  - [x] 5.3 Implement `LoadPromptTemplate(config Config) (string, error)` with config override support
  - [x] 5.4 Create `cmd/generate.go` implementing `ctx generate` command
  - [x] 5.5 In generate: Load index and filter files by status (stale/missing)
  - [x] 5.6 In generate: Add `--filter` flag to specify statuses (stale,missing)
  - [x] 5.7 In generate: Add `--files` flag to generate for specific files only
  - [x] 5.8 In generate: Add `--output` flag to save prompt to file instead of stdout
  - [x] 5.9 In generate: For each file, read source content and build complete prompt
  - [x] 5.10 In generate: Format output as markdown with file sections, skeleton paths, and index update instructions
  - [x] 5.11 In generate: Mark processed files as `pendingGeneration` in index
  - [x] 5.12 In generate: Display instructions for feeding prompt to AI and next steps
  - [x] 5.13 Create `cmd/export.go` implementing `ctx export` command
  - [x] 5.14 In export: Load index and all skeleton files with status "current"
  - [x] 5.15 In export: Add `--format` flag supporting "markdown" and "json"
  - [x] 5.16 In export: Add `--output` flag to save to file instead of stdout
  - [x] 5.17 In export: Generate consolidated markdown with index summary and all skeletons
  - [x] 5.18 In export: For JSON format, output structured data with index + skeletons object

- [ ] 6.0 Testing, Documentation & Distribution
  - [ ] 6.1 Write integration tests in `main_test.go` covering full init→sync→generate flow
  - [ ] 6.2 Create test fixtures with sample files for integration testing
- [ ] 6.3 Ensure unit test coverage reaches 80% minimum (run `go test -cover ./...`) *(in progress – see README for coverage command and target)*
  - [x] 6.4 Write `README.md` with quick start (3-command example), installation, basic usage
  - [x] 6.5 Create `docs/getting-started.md` with step-by-step tutorial for new users
  - [x] 6.6 Create `docs/commands.md` with complete reference for all 8 commands
  - [x] 6.7 Create `docs/workflows.md` with common usage patterns (daily sync, initial setup, maintenance)
  - [x] 6.8 Create `docs/troubleshooting.md` with FAQ and error solutions
  - [x] 6.9 Set up `.goreleaser.yml` for multi-platform builds (Linux, macOS, Windows; amd64, arm64)
  - [ ] 6.10 Create `.github/workflows/test.yml` for automated testing on push/PR *(skipped – GitHub-specific)*
  - [ ] 6.11 Create `.github/workflows/release.yml` for automated binary releases on git tags *(skipped – GitHub-specific)*
  - [ ] 6.12 Test binary builds locally using `goreleaser build --snapshot` *(pending local tooling)*
  - [ ] 6.13 Create GitHub repository with appropriate license (MIT or Apache 2.0)
  - [x] 6.14 Write CONTRIBUTING.md with guidelines for external contributors
  - [x] 6.15 Add badges to README (build status, Go version, license) *(placeholders noted until CI is configured)*
  - [ ] 6.16 Tag v1.0.0 release and verify automated builds publish correctly
  - [ ] 6.17 Create example repository showing ctx usage with before/after context
  - [ ] 6.18 Post announcement on relevant communities (r/golang, Twitter/X, Dev.to)

---

## Execution Notes

**Recommended order:**
1. Complete tasks 1.0 → 2.0 → 3.0 sequentially (foundation)
2. Implement commands in 4.0 incrementally (start with init, sync, status)
3. Add generate/export in 5.0 once core is stable
4. Polish with 6.0 after core functionality works end-to-end

**Testing strategy:**
- Write unit tests alongside each package implementation
- Run `go test ./...` frequently to catch regressions
- Add integration tests after core commands work
- Aim for 80%+ code coverage before Phase 1 completion

**Milestone checkpoints:**
- After 1.0-3.0: Basic infrastructure complete, can scan and index files
- After 4.0: All commands functional, can manage skeletons manually
- After 5.0: Full workflow operational (init → sync → generate → export)
- After 6.0: Production-ready with docs, tests, and distribution

**Common pitfalls to avoid:**
- Don't skip unit tests - they catch edge cases early
- Test with large codebases (500+ files) to validate performance
- Verify git integration works in repos with many uncommitted changes
- Ensure proper error messages guide users to solutions
