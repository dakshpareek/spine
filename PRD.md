# Code Context CLI - Phase 1 PRD

## Executive Summary

**Product Name:** ctx (`ctx`)
**Version:** 1.0.0 (Phase 1 - Manual Mode)
**Language:** Go
**Target Users:** Developers using AI coding assistants (ChatGPT, Claude, etc.)
**Core Problem:** AI agents lack architectural context when writing code, leading to inconsistent patterns, duplicate implementations, and context-unaware solutions.
**Solution:** Maintain a lightweight, always-fresh mini-replica of your codebase with structural skeletons that AI agents can consume for full architectural awareness.

---

## Product Overview

Code Context CLI creates and maintains a shadow filesystem (`.ctx/`) that mirrors your codebase structure with compact "skeleton" files. These skeletons capture the essential structure, patterns, and conventions of each file without implementation details. AI agents consume these skeletons to understand your entire architecture before writing code.

### Key Principles
- **Zero coupling:** Tool operates in isolated `.ctx/` directory
- **Manual first:** Phase 1 uses AI copy-paste workflow (no API integration)
- **Git-aware:** Leverages git for efficient change detection
- **Single source of truth:** One mini-replica serves all features/tasks
- **Minimal storage:** Skeletons are ~6% of original code size

---

## Phase 1 Scope

### In Scope
✅ Initialize `.ctx/` structure
✅ Scan codebase and build file index
✅ Detect file changes via git diff and hashing
✅ Generate AI prompts for skeleton creation (manual copy-paste)
✅ Import and validate AI-generated skeletons
✅ Update index status after skeleton generation
✅ Configuration for file extensions and exclusions
✅ Status reporting and validation commands

### Out of Scope (Phase 2+)
❌ Automated AI API integration
❌ Watch mode / real-time updates
❌ Git hooks for auto-sync
❌ Dependency graph analysis
❌ Smart context selection for features
❌ Team collaboration features

---

## Technical Architecture

### Directory Structure
```
project-root/
├── .ctx/              # Hidden directory (gitignored by default)
│   ├── index.json             # Master index with file metadata
│   ├── config.json            # User configuration
│   ├── skeleton-prompt.txt    # Default skeleton generation prompt
│   └── skeletons/             # Mirror of source structure
│       └── src/
│           └── booking/
│               └── booking.service.skeleton.ts
├── .gitignore                 # Auto-updated to ignore .ctx/
└── src/                       # User's actual codebase
    └── booking/
        └── booking.service.ts
```

### Data Models

#### Index Schema (`index.json`)
```json
{
  "version": "1.0.0",
  "promptVersion": "2.1",
  "lastSync": "2025-11-04T10:30:00Z",
  "config": {
    "includedExtensions": [".ts", ".tsx", ".js", ".jsx", ".go"],
    "excludedPaths": [
      "node_modules",
      "dist",
      "build",
      ".next",
      "vendor",
      "*.test.ts",
      "*.spec.ts"
    ]
  },
  "files": {
    "src/booking/booking.service.ts": {
      "path": "src/booking/booking.service.ts",
      "hash": "a3f2e1c9...",
      "skeletonHash": "b7d4a8f2...",
      "skeletonPath": ".ctx/skeletons/src/booking/booking.service.skeleton.ts",
      "lastModified": "2025-11-04T09:15:00Z",
      "status": "current",
      "type": "service",
      "size": 15234
    }
  },
  "stats": {
    "totalFiles": 247,
    "current": 245,
    "stale": 2,
    "missing": 0,
    "pendingGeneration": 0
  }
}
```

#### File Status Types
- `current` - Skeleton is up-to-date with source file
- `stale` - Source file modified, skeleton needs regeneration
- `missing` - No skeleton exists yet
- `pendingGeneration` - Prompt generated, waiting for AI response

#### Config Schema (`config.json`)
```json
{
  "includedExtensions": [".ts", ".tsx", ".js", ".jsx", ".go"],
  "excludedPaths": [
    "node_modules",
    "dist",
    "build",
    ".next",
    "coverage",
    "vendor",
    "*.test.*",
    "*.spec.*"
  ],
  "skeletonPromptVersion": "2.1",
  "rootPath": "."
}
```

---

## Command Specifications

### 1. `ctx init`

**Purpose:** Initialize `.ctx/` structure in current directory

**Behavior:**
1. Check if `.ctx/` already exists (error if yes)
2. Create directory structure:
   ```
   .ctx/
   ├── index.json (empty initial state)
   ├── config.json (default config)
   ├── skeleton-prompt.txt (embedded default prompt)
   └── skeletons/ (empty directory)
   ```
3. Add `.ctx/` to `.gitignore` (create if missing)
4. Run initial scan (same as `ctx sync`)

**Output:**
```
✓ Initialized .ctx/
✓ Updated .gitignore
✓ Scanning codebase...
  Found 247 files
✓ Index created

Next steps:
  1. Run 'ctx generate' to create skeleton prompts
  2. Feed prompts to your AI assistant
  3. Run 'ctx status' to check progress
```

**Error Cases:**
- `.ctx/` already exists → "Already initialized. Use 'ctx rebuild' to reset."
- Not in a project root → "Run from project root directory"
- No write permissions → "Permission denied"

---

### 2. `ctx sync`

**Purpose:** Scan codebase for changes and update index

**Behavior:**
1. Load existing `index.json`
2. Get list of modified files:
   - If git repo: `git diff --name-only HEAD` + untracked files
   - If no git: Compare file mtimes with index
3. For each file in codebase matching config:
   - Calculate SHA-256 hash of content
   - Compare with index hash
   - Update status:
     - Hash unchanged + skeleton exists → `current`
     - Hash changed → `stale`
     - New file → `missing`
4. Remove deleted files from index
5. Update `index.json` with new state
6. Calculate and display stats

**Flags:**
- `--full` - Scan all files (ignore git diff)
- `--verbose` - Show detailed file-by-file changes

**Output:**
```
Scanning codebase...
✓ 247 files scanned

Changes detected:
  • 3 modified (marked stale)
  • 1 new file (marked missing)
  • 1 deleted

Status:
  ✓ 242 current
  ⚠ 3 stale
  ⚠ 2 missing

Run 'ctx generate' to create skeleton prompts.
```

---

### 3. `ctx generate`

**Purpose:** Generate AI prompts for skeleton creation (manual mode)

**Behavior:**
1. Load `index.json`
2. Filter files with status: `stale` or `missing`
3. For each file:
   - Read source file content
   - Load skeleton prompt template
   - Generate complete prompt with instructions
4. Output consolidated prompt to stdout
5. Mark files as `pendingGeneration` in index

**Flags:**
- `--filter=stale,missing` - Generate only for specific statuses
- `--files=path1,path2` - Generate for specific files only
- `--output=prompt.txt` - Save to file instead of stdout

**Output Format:**
```markdown
# Code Context Skeleton Generation

You are generating structural skeletons for a codebase. Follow these instructions carefully.

## Instructions

1. For each file below, generate a skeleton using the template provided
2. Create/update skeleton files in the specified paths
3. Update the index.json file with the provided updates

## Skeleton Generation Template

[Full skeleton prompt from skeleton-prompt.txt]

## Files to Process

### File 1: src/booking/booking.service.ts
**Status:** stale
**Skeleton Path:** .ctx/skeletons/src/booking/booking.service.skeleton.ts

**Source Code:**
```typescript
[Full file content]
```

---

### File 2: src/voucher/voucher.repository.ts
**Status:** missing
**Skeleton Path:** .ctx/skeletons/src/voucher/voucher.repository.skeleton.ts

**Source Code:**
```typescript
[Full file content]
```

---

## Index Updates Required

After generating skeletons, update `.ctx/index.json`:

For each file processed:
1. Create/update the skeleton file at the specified path
2. Calculate SHA-256 hash of the skeleton content
3. Update the file entry:
   ```json
   "src/booking/booking.service.ts": {
     "status": "current",
     "skeletonHash": "<calculated-hash>",
     "lastModified": "<current-timestamp>"
   }
   ```
4. Update stats section:
   - Decrement stale/missing counts
   - Increment current count

**Index Path:** .ctx/index.json

## Verification

After completion, the user will run `ctx status` to verify all files are current.
```

**Post-Generation:**
- Update index: mark files as `pendingGeneration`
- Show instruction to user about next steps

---

### 4. `ctx status`

**Purpose:** Display current state of code context

**Behavior:**
1. Load `index.json`
2. Calculate statistics
3. Display formatted summary
4. Optionally show file-by-file details

**Flags:**
- `--verbose` - Show list of stale/missing files
- `--json` - Output as JSON for scripting

**Output:**
```
Code Context Status

Overview:
  Total files: 247
  ✓ Current: 245 (99%)
  ⚠ Stale: 2
  ⚠ Missing: 0
  ⏳ Pending: 0

Last sync: 2 minutes ago
Skeleton prompt: v2.1

Next steps:
  Run 'ctx generate' to create skeletons for stale files.
```

**Verbose Output:**
```
Stale files:
  • src/booking/booking.service.ts (modified 5 minutes ago)
  • src/voucher/voucher.repo.ts (modified 1 hour ago)
```

---

### 5. `ctx validate`

**Purpose:** Verify integrity of skeletons and index

**Behavior:**
1. Load `index.json`
2. For each file entry:
   - Check if source file exists
   - Check if skeleton file exists
   - Verify source hash matches current file
   - Verify skeleton hash matches skeleton file
3. Report any inconsistencies
4. Optionally auto-fix issues

**Flags:**
- `--fix` - Automatically mark mismatched files as stale
- `--strict` - Exit with error code if issues found

**Output:**
```
Validating code context...

Issues found:
  ⚠ src/booking/service.ts: source hash mismatch (marked stale)
  ⚠ src/old/legacy.ts: source file deleted (removed from index)
  ✓ src/voucher/repo.ts: skeleton file missing (marked missing)

Summary:
  2 files marked stale
  1 file removed from index

Run 'ctx generate' to resolve issues.
```

---

### 6. `ctx clean`

**Purpose:** Remove orphaned skeleton files

**Behavior:**
1. Scan `.ctx/skeletons/` directory
2. Compare with files in `index.json`
3. Delete skeleton files not referenced in index
4. Remove empty directories

**Output:**
```
Cleaning orphaned skeletons...
  Removed 3 orphaned files
  Removed 1 empty directory
✓ Clean complete
```

---

### 7. `ctx rebuild`

**Purpose:** Reset and regenerate entire context

**Behavior:**
1. Confirm destructive action (require `--confirm` flag)
2. Delete all skeleton files
3. Reset index (keep config)
4. Run full sync
5. Show generate instructions

**Flags:**
- `--confirm` - Skip confirmation prompt

**Output:**
```
⚠ This will delete all existing skeletons.
Continue? (y/N): y

Rebuilding...
  ✓ Deleted 245 skeleton files
  ✓ Reset index
  ✓ Scanning codebase...

Found 247 files (all marked missing)

Run 'ctx generate' to recreate skeletons.
```

---

### 8. `ctx export`

**Purpose:** Export context for AI consumption

**Behavior:**
1. Load `index.json` and all skeleton files
2. Generate consolidated markdown document
3. Output to stdout or file

**Flags:**
- `--output=context.md` - Save to file
- `--format=json|markdown` - Output format

**Markdown Output:**
```markdown
# Codebase Context

**Generated:** 2025-11-04T10:30:00Z
**Files:** 247 (245 current, 2 stale)
**Prompt Version:** 2.1

---

## Index Summary

- Total Files: 247
- Current Skeletons: 245
- Coverage: 99%

## File Skeletons

### src/booking/booking.service.ts

[Full skeleton content]

---

### src/voucher/voucher.repository.ts

[Full skeleton content]

---

[... rest of skeletons ...]
```

---

## Default Skeleton Prompt

Embedded in tool as `skeleton-prompt.txt`:

```
You are documenting a code file for downstream AI agents. Given the file contents, produce a concise structural summary that enables AI agents to understand existing patterns, reusable components, and architectural approaches when writing new code.

Use this exact template:

**<Display Name>**
- File: `<PATH_FROM_ROOT>:1`
- Imports: <one-sentence summary of import groups and architectural purpose>
- State: <describe module/class state and lifecycle; use "None" if stateless>
- Constructor: <one sentence on dependencies injected; omit if no constructor>

**Public Methods**
- Method: `<name(params) -> returnType>` — <purpose, key behaviors, return semantics in 1-2 sentences>.

**Private Helpers**
- Helper: `<name(params?) -> returnType>` — <purpose and trigger context in one sentence>.
  (Include only non-trivial helpers; omit section if none.)

**Key Dependencies**
- <One bullet per notable collaborator describing role and integration pattern.>

**Patterns & Conventions**
- Error handling: <exception types, error mapping approach>
- Validation: <strategy and sequencing>
- Persistence: <transaction/repository patterns>
- Naming: <suffixes, conventions>
- (Omit if no clear patterns)

**Usage Signals**
- Side effects: <database writes, API calls, events>
- Resilience: <idempotency, retry, caching mechanisms>
- Observability: <logging approach>
- (Use "None" if nothing notable)

**Guidelines:**
- Maximum 2 sentences per method/helper description
- Focus on *what* and *why*, not implementation details
- Highlight reusable patterns for consistency
- Note architectural decisions (e.g., "uses idempotency keys")
- Call out implicit contracts (e.g., "expects pre-validated voucher")
- Omit empty sections entirely
```

---

## File Type Detection

Auto-detect file type from path patterns:

```go
var fileTypes = map[string][]string{
    "service":     {"*service.ts", "*service.js", "*service.go"},
    "controller":  {"*controller.ts", "*controller.js", "*handler.go"},
    "repository":  {"*repository.ts", "*repo.ts", "*repository.go"},
    "dto":         {"*dto.ts", "*dto.go", "*/dto/*"},
    "model":       {"*model.ts", "*entity.ts", "*model.go"},
    "util":        {"*util.ts", "*utils.ts", "*helper.ts"},
    "middleware":  {"*middleware.ts", "*middleware.go"},
    "config":      {"*config.ts", "*config.go"},
}
```

---

## Technical Implementation Details

### Hashing Strategy
```go
import "crypto/sha256"

func hashFile(path string) (string, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    hash := sha256.Sum256(content)
    return fmt.Sprintf("%x", hash), nil
}
```

### Git Integration
```go
func getModifiedFiles() ([]string, error) {
    // Check if git repo
    if !isGitRepo() {
        return nil, errors.New("not a git repository")
    }

    // Get modified tracked files
    cmd := exec.Command("git", "diff", "--name-only", "HEAD")
    output, _ := cmd.Output()

    // Get untracked files
    cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
    untracked, _ := cmd.Output()

    return append(
        strings.Split(string(output), "\n"),
        strings.Split(string(untracked), "\n")...,
    ), nil
}
```

### File Scanning with Exclusions
```go
import "github.com/bmatcuk/doublestar/v4"

func scanFiles(config Config) ([]string, error) {
    var files []string

    err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        // Skip excluded paths
        for _, pattern := range config.ExcludedPaths {
            matched, _ := doublestar.Match(pattern, path)
            if matched {
                if d.IsDir() {
                    return filepath.SkipDir
                }
                return nil
            }
        }

        // Check extension
        if !d.IsDir() {
            ext := filepath.Ext(path)
            if contains(config.IncludedExtensions, ext) {
                files = append(files, path)
            }
        }

        return nil
    })

    return files, err
}
```

---

## User Workflows

### Initial Setup
```bash
# Day 1: Initialize
cd my-project/
ctx init

# Output: Found 247 files, all marked missing

ctx generate > skeleton-prompt.txt
# Copy skeleton-prompt.txt to Claude/ChatGPT
# AI generates skeletons and updates index

ctx status
# Output: 247/247 current ✓
```

### Daily Usage
```bash
# Morning: Sync changes
ctx sync
# Output: 3 stale, 1 new file

ctx generate > prompt.txt
# Feed to AI, AI handles skeleton + index update

# Start coding with full context
ctx export > context.md
# Feed context.md to AI when writing new code
```

### Maintenance
```bash
# Check integrity
ctx validate

# Clean orphaned files
ctx clean

# Full rebuild (e.g., after prompt update)
ctx rebuild --confirm
```

---

## Error Handling

### Error Categories
1. **User Errors** (exit code 1)
   - Not initialized (`ctx sync` before `ctx init`)
   - Invalid config
   - Missing required flags

2. **File System Errors** (exit code 2)
   - Permission denied
   - Disk full
   - File not found

3. **Git Errors** (exit code 3)
   - Not a git repo (when using git-based sync)
   - Git command failed

4. **Data Errors** (exit code 4)
   - Corrupted index.json
   - Invalid JSON format

### Error Messages
- Always actionable: "Run 'ctx init' to initialize"
- Include context: "Failed to read src/booking/service.ts: permission denied"
- Suggest fixes: "Index is corrupted. Run 'ctx rebuild --confirm' to reset."

---

## Configuration

### Default Config
```json
{
  "includedExtensions": [".ts", ".tsx", ".js", ".jsx", ".go", ".py"],
  "excludedPaths": [
    "node_modules",
    "vendor",
    "dist",
    "build",
    ".next",
    "coverage",
    "*.test.*",
    "*.spec.*",
    "__tests__",
    "test"
  ],
  "skeletonPromptVersion": "2.1",
  "rootPath": "."
}
```

### User Overrides
Users can edit `.ctx/config.json` directly. Changes take effect on next `ctx sync`.

---

## Testing Requirements

### Unit Tests
- File hashing (SHA-256 correctness)
- Path matching (exclusions/inclusions)
- Index serialization/deserialization
- Status transitions (current → stale → missing)

### Integration Tests
- Full init → sync → generate → validate flow
- Git integration (mocked git commands)
- File system operations (temp directories)
- Error handling for edge cases

### Test Coverage Target
- Minimum 80% code coverage
- 100% coverage for critical paths (hashing, index updates)

---

## Distribution

### Binary Names
- `ctx` (primary)
- `code-context` (alternative)

### Installation Methods
```bash
# Homebrew (macOS/Linux)
brew install code-context

# Go install
go install github.com/yourusername/code-context@latest

# Download binary
curl -L https://github.com/.../releases/latest/download/ctx-linux-amd64 -o ctx
chmod +x ctx
sudo mv ctx /usr/local/bin/
```

### Releases
- GitHub Releases with binaries for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)

---

## Success Metrics (Phase 1)

### Must Have
- ✅ Initialize and scan 500+ file codebase in <5 seconds
- ✅ Detect changes via git diff in <1 second
- ✅ Generate prompts for 50 files in <2 seconds
- ✅ Binary size <10MB
- ✅ Zero external dependencies (statically linked)

### Nice to Have
- User feedback: "Saves 2+ hours per week"
- Adoption: 100+ GitHub stars in first month
- Community: 5+ external contributors

---

## Documentation Requirements

### README.md
- Quick start (3 commands to working state)
- Installation instructions
- Basic workflow example
- Link to full docs

### docs/
- `getting-started.md` - Step-by-step tutorial
- `commands.md` - Complete command reference
- `workflows.md` - Common usage patterns
- `troubleshooting.md` - FAQ and error solutions
- `architecture.md` - Technical design decisions

---

## Phase 2 Preview (Future)

After Phase 1 validates the concept:

1. **Auto Mode:** Direct AI API integration
   ```bash
   ctx generate --auto --provider=openai
   ```

2. **Watch Mode:** Real-time skeleton updates
   ```bash
   ctx watch
   ```

3. **Smart Context:** Dependency-aware exports
   ```bash
   ctx export --related-to=booking.service.ts
   ```

4. **Team Sync:** Shared context in git
   ```bash
   ctx push  # Commit .ctx/ to repo
   ```

---

## Open Questions (Resolve Before Implementation)

1. **Prompt Storage:** Embed in binary or require external file?
   - **Decision:** Embed default, allow override via config

2. **Gitignore Behavior:** Always ignore `.ctx/` or make configurable?
   - **Decision:** Default ignore, document how to commit if desired

3. **Binary Name:** `ctx` vs `code-context` vs `skeleton`?
   - **Decision:** `ctx` (short, memorable)

4. **Skeleton Extension:** `.skeleton.ts` vs `.skel.ts` vs `.ctx.ts`?
   - **Decision:** `.skeleton.ts` (clear intent)

5. **Index Format:** JSON vs YAML vs custom?
   - **Decision:** JSON (better Go support, universal)

---

## Development Milestones

### Milestone 1: Core Engine (Week 1)
- [ ] Project setup (Go modules, CI/CD)
- [ ] Index data structures and serialization
- [ ] File scanning with exclusions
- [ ] Hash calculation
- [ ] Basic `init` and `sync` commands

### Milestone 2: Generation Flow (Week 2)
- [ ] Skeleton prompt template
- [ ] `generate` command with file filtering
- [ ] Prompt formatting and output
- [ ] Status transition logic
- [ ] Git integration for change detection

### Milestone 3: Utilities (Week 3)
- [ ] `status` command with formatting
- [ ] `validate` command with auto-fix
- [ ] `clean` command
- [ ] `rebuild` command
- [ ] `export` command

### Milestone 4: Polish (Week 4)
- [ ] Error handling and messages
- [ ] Unit and integration tests
- [ ] Documentation (README, command docs)
- [ ] Binary builds for all platforms
- [ ] GitHub release automation

---

## Appendix: Example Session

```bash
# Initialize new project
$ ctx init
✓ Initialized .ctx/
✓ Scanning codebase...
  Found 247 files (all marked missing)

Run 'ctx generate' to create skeleton prompts.

# Generate prompts for AI
$ ctx generate > prompt.txt
✓ Generated prompts for 247 files
✓ Marked as pending generation

Feed prompt.txt to your AI assistant.

# [User copies prompt.txt to Claude]
# [Claude generates skeletons and updates index]

# Verify completion
$ ctx status
Code Context Status

Overview:
  Total files: 247
  ✓ Current: 247 (100%)

Last sync: just now
Skeleton prompt: v2.1

✓ All skeletons up to date!

# Daily workflow: detect changes
$ ctx sync
Scanning codebase...
✓ 247 files scanned

Changes detected:
  • 3 modified (marked stale)
  • 1 new file (marked missing)

Status:
  ✓ 243 current
  ⚠ 4 need regeneration

Run 'ctx generate' to update skeletons.

# Export context for coding session
$ ctx export --output=context.md
✓ Exported 247 skeletons to context.md

# [User feeds context.md to AI while coding]
# [AI writes code consistent with existing patterns]
```

---

## Conclusion

Phase 1 delivers a minimal, functional tool that solves the core problem: maintaining a fresh, AI-consumable view of codebase architecture. The manual workflow validates the concept without complex API integration, enabling rapid iteration based on user feedback.

**Next Steps:**
1. Review and approve PRD
2. Set up Go project structure
3. Begin Milestone 1 development
4. Create public GitHub repository
5. Build community around the tool

**Success Criteria:**
Phase 1 is complete when a developer can:
1. Initialize context in <30 seconds
2. Sync daily changes in <5 seconds
3. Generate AI prompts for any number of files
4. Export complete context for AI coding sessions
5. All with zero external dependencies or configuration
