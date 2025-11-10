# Troubleshooting

Common problems and how to fix them.

## `ctx init` Fails with “Already initialized”

`ctx init` aborts if `.ctx/` exists.

- Run `ctx rebuild --confirm` to reset the workspace.
- Or delete `.ctx/` manually if you are sure it’s stale.

## `ctx sync` Doesn’t Detect Changes

- Ensure you are running from the project root.
- For git repositories, confirm the project is a valid Git worktree (`git status`).
- Use `ctx sync --full` to force a complete scan.

## `ctx generate` Returns “No files match the requested filters”

- Run `ctx sync` first so the index reflects recent edits.
- Check status with `ctx status --verbose`.
- Use `--filter current,pending` if you intentionally want other statuses.

## Skeleton Files Missing After AI Update

- Run `ctx validate --fix` to mark missing skeletons and restore index health.
- Re-run `ctx generate` for the affected files.

## Export Fails with “no current skeletons to export”

- Ensure you have at least one file marked `current`.
- Use `ctx sync` + `ctx generate` + AI update to bring skeletons current.

## Hash Mismatch Warnings

`ctx validate` or `ctx sync` may flag mismatched hashes when skeleton content changes without updating `index.json`.

1. Re-run `ctx generate` for the affected file.
2. Update the skeleton file.
3. Run `ctx sync` or `ctx validate --fix` to refresh hashes.

## Still Stuck?

Capture the command output and index snippet, then open an issue or start a discussion in your repository. Include:

- Command run and flags.
- Relevant status output.
- Excerpts from `.ctx/index.json`.
