# Command Reference

`ctx` exposes a collection of commands for managing your `.spine/` workspace. Each command accepts `--help` for detailed usage.

## `spine init`

Initializes `.spine/` in the current directory.

- Creates `config.json`, `index.json`, and `skeletons/`.
- Adds `.spine/` to `.gitignore`.
- Scans the project and marks files as `missing`.

> Run from your project root. If the workspace already exists, use `spine rebuild --confirm`.

## `spine sync`

Scans your project for changes.

Flags:

- `--full` – force a full scan (ignore Git diffs).
- `--verbose`, `-v` – list individual file changes.

Behaviour:

- Updated files → `stale`
- New files → `missing`
- Removed files → deleted from index

## `spine generate`

Builds a Markdown prompt for an AI assistant to generate skeletons.

Flags:

- `--filter` – statuses to include (`stale`, `missing`, `current`, `pending`).
- `--files` – comma-separated file paths to target explicitly.
- `--output`, `-o` – write prompt to a file (stdout by default).

Post-conditions:

- Selected entries are marked as `pendingGeneration`.
- Prompts include source code, skeleton paths, and index update instructions.

## `spine pipeline`

Runs `sync` and `generate` back-to-back.

Flags:

- `--full` – force a full scan before generating.
- `--verbose`, `-v` – show detailed sync output.
- `--filter`, `--files`, `--output` – same as `spine generate`.

Use this when you want a single step before copying the prompt to your AI assistant.

## `spine status`

Displays index statistics.

Flags:

- `--verbose`, `-v` – lists files by status.
- `--json` – emits machine-readable JSON.

## `spine validate`

Checks integrity between source files, skeletons, and index metadata.

Flags:

- `--fix` – automatically mark stale/missing entries and repair hashes.
- `--strict` – non-zero exit code if issues are detected.

## `spine clean`

Removes orphaned skeleton files within `.spine/skeletons/` and prunes empty directories.

## `spine rebuild`

Resets the `.spine/` workspace (destructive).

Flags:

- `--confirm` – required safety flag.

Steps:

1. Deletes all skeletons.
2. Resets the index.
3. Performs a full scan.

## `spine export`

Collects all `current` skeletons into a single artifact.

Flags:

- `--format` – `markdown` (default) or `json`.
- `--output`, `-o` – write export to a file (stdout by default).

Use this before long coding sessions to preload context into an AI assistant.
