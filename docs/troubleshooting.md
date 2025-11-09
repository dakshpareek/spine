# Troubleshooting

Common problems and how to fix them.

## `spine init` Fails with “Already initialized”

`spine init` aborts if `.spine/` exists.

- Run `spine rebuild --confirm` to reset the workspace.
- Or delete `.spine/` manually if you are sure it’s stale.

## `spine sync` Doesn’t Detect Changes

- Ensure you are running from the project root.
- For git repositories, confirm the project is a valid Git worktree (`git status`).
- Use `spine sync --full` to force a complete scan.

## `spine generate` Returns “No files match the requested filters”

- Run `spine sync` first so the index reflects recent edits.
- Check status with `spine status --verbose`.
- Use `--filter current,pending` if you intentionally want other statuses.

## Skeleton Files Missing After AI Update

- Run `spine validate --fix` to mark missing skeletons and restore index health.
- Re-run `spine generate` for the affected files.

## Export Fails with “no current skeletons to export”

- Ensure you have at least one file marked `current`.
- Use `spine sync` + `spine generate` + AI update to bring skeletons current.

## Hash Mismatch Warnings

`spine validate` or `spine sync` may flag mismatched hashes when skeleton content changes without updating `index.json`.

1. Re-run `spine generate` for the affected file.
2. Update the skeleton file.
3. Run `spine sync` or `spine validate --fix` to refresh hashes.

## Still Stuck?

Capture the command output and index snippet, then open an issue or start a discussion in your repository. Include:

- Command run and flags.
- Relevant status output.
- Excerpts from `.spine/index.json`.
