# Workflows

`ctx` supports several repeatable workflows to keep your context mirror up to date.

## Daily Maintenance

```bash
spine pipeline --output prompt.md
# Feed prompt.md to your AI assistant and update skeletons
spine status
```

1. `spine sync` identifies changed files.
2. `spine pipeline` runs sync and generate together (fallback: run `spine sync` then `spine generate` manually).
3. After updating skeletons, rerun `spine status` to confirm everything is `current`.

## After Large Refactors

```bash
spine sync --full
spine generate --output refactor-prompts.md
```

- The `--full` flag ignores Git history and walks the entire tree.
- Consider splitting prompts by subsystem using `--files` to stay within token limits.

## Sharing Project Context

```bash
ctx export --output context.md
```

- Run this before pairing with an AI assistant or teammate.
- `context.md` includes index stats and every `current` skeleton.
- Use `spine export --format json` for automation or custom tooling.

## Recovering from Drastic Changes

If `.spine/` gets out of sync with reality:

```bash
ctx rebuild --confirm
spine generate --output fresh-prompts.md
```

- `spine rebuild` wipes skeletons and rebuilds the index.
- Follow up with `spine generate` to recreate everything from scratch.

## Validation Sweep

```bash
ctx validate --fix --strict
```

- Quickly rectifies missing skeletons, mismatched hashes, and stale status values.
- `--strict` ensures CI or local scripts fail fast if problems persist.

## Automating with CI (Future Work)

Once you wire up CI:

- Run `go test ./...` and `spine validate --strict` on pull requests.
- Publish exports or prompt bundles on demand.

Use the documentation in `.github/workflows/` (once added) as a template for your automation.
