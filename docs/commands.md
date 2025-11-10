# Command Reference

`ctx` exposes a collection of commands for managing your `.ctx/` workspace. Each command accepts `--help` for detailed usage.

## `ctx init`

Initializes `.ctx/` in the current directory.

- Creates `config.json`, `index.json`, and `skeletons/`.
- Adds `.ctx/` to `.gitignore`.
- Scans the project and marks files as `missing`.

> Run from your project root. If the workspace already exists, use `ctx rebuild --confirm`.

## `ctx sync`

Scans your project for changes.

Flags:

- `--full` – force a full scan (ignore Git diffs).
- `--verbose`, `-v` – list individual file changes.

Behaviour:

- Updated files → `stale`
- New files → `missing`
- Removed files → deleted from index

## `ctx generate`

Builds a Markdown prompt for an AI assistant to generate skeletons.

Flags:

- `--filter` – statuses to include (`stale`, `missing`, `current`, `pending`).
- `--files` – comma-separated file paths to target explicitly.
- `--output`, `-o` – write prompt to a file (stdout by default).

Post-conditions:

- Selected entries are marked as `pendingGeneration`.
- Prompts include source code, skeleton paths, and index update instructions.

## `ctx pipeline`

Runs `sync` and `generate` back-to-back.

Flags:

- `--full` – force a full scan before generating.
- `--verbose`, `-v` – show detailed sync output.
- `--filter`, `--files`, `--output` – same as `ctx generate`.

Use this when you want a single step before copying the prompt to your AI assistant.

## `ctx status`

Displays index statistics.

Flags:

- `--verbose`, `-v` – lists files by status.
- `--json` – emits machine-readable JSON.

## `ctx validate`

Checks integrity between source files, skeletons, and index metadata.

Flags:

- `--fix` – automatically mark stale/missing entries and repair hashes.
- `--strict` – non-zero exit code if issues are detected.

## `ctx clean`

Removes orphaned skeleton files within `.ctx/skeletons/` and prunes empty directories.

## `ctx rebuild`

Resets the `.ctx/` workspace (destructive).

Flags:

- `--confirm` – required safety flag.

Steps:

1. Deletes all skeletons.
2. Resets the index.
3. Performs a full scan.

## `ctx export`

Collects all `current` skeletons into a single artifact.

Flags:

- `--format` – `markdown` (default) or `json`.
- `--output`, `-o` – write export to a file (stdout by default).

Use this before long coding sessions to preload context into an AI assistant.
