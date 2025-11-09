# Getting Started

Welcome to `ctx`, the Code Context CLI. Follow this guide to bootstrap a new project and keep your AI companions fully informed.

## Prerequisites

- Go 1.21+ installed.
- A Git repository (recommended for fast change detection, but not required).

## Installation

```bash
go install github.com/yourusername/spine@latest
```

This places the `ctx` binary in your `$GOBIN` (`$GOPATH/bin` by default).

## Initialize a Project

```bash
cd /path/to/project
ctx init
```

The command:

1. Creates `.spine/`, `config.json`, `index.json`, and `skeletons/`.
2. Adds `.spine/` to `.gitignore`.
3. Scans your project for trackable files and marks them as `missing`.

## Sync Daily Changes

```bash
ctx sync
```

`ctx` compares tracked files with your working tree:

- Modified files → `stale`
- New files → `missing`
- Deleted files → removed from index

## Generate Skeleton Prompts

```bash
ctx generate --filter stale,missing --output prompt.md
```

Share `prompt.md` with your AI assistant. It contains:

- Instructions for producing structural skeletons.
- Source code for each selected file.
- Paths where skeletons must be written.

## Update Skeletons

After the AI writes skeletons:

1. Save each skeleton to the path under `.spine/skeletons/…`.
2. Update `.spine/index.json` with the new `skeletonHash`, `status: "current"`, and timestamp (or run `spine sync` with `--full`).

## Verify Status

```bash
ctx status --verbose
```

You’ll see a summary plus detailed lists of `stale`, `missing`, and `pendingGeneration` files.

## Export Context

```bash
ctx export --output context.md
```

Use the export when you start a fresh coding session and need to load complete context into your AI assistant.

## Next Steps

- Automate skeleton updates inside your team’s workflow.
- Track coverage with `go test ./... -cover`.
- Contribute improvements—see [Contributing](../CONTRIBUTING.md).
