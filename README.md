# Spine – Extract Your Codebase's Architecture

Maintain an AI-ready snapshot of your codebase. `spine` extracts the structural skeleton of your project into a lightweight `.code-context/` workspace so coding assistants can understand architecture before they touch your source files.

## Quick Start

```bash
go install github.com/dakshpareek/spine@latest

cd /path/to/project
spine init                # bootstrap .code-context and initial index
spine generate --output prompt.md
# or combine the first two steps:
# spine pipeline --output prompt.md
```

1. Paste `prompt.md` into your AI assistant and let it create/update skeletons.
2. Save the generated skeletons under `.code-context/skeletons/`.
3. Run `spine status` to confirm everything is current.

## Daily Workflow

```bash
spine pipeline --output prompt.md   # sync + generate in one step
spine validate --fix                # mark skeletons current after AI updates
spine export --output context.md    # optional: share full context
```

- `spine pipeline` runs `sync` then prints the AI prompt (or writes to a file).
- `spine validate --fix` recomputes hashes after you save skeletons.
- `spine export` collects all current skeletons for context-heavy coding sessions.

## Documentation

- [Getting Started](docs/getting-started.md)
- [Command Reference](docs/commands.md)
- [Workflows](docs/workflows.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Contributing](CONTRIBUTING.md)

## Testing

```bash
go test ./... -cover
```

Integration coverage lives in `main_test.go`, exercising the init → sync → generate flow end-to-end using sample fixtures.

## Development & Distribution

We ship cross-platform binaries via [GoReleaser](https://goreleaser.com/). See `.goreleaser.yml` for the current configuration.

### Building Locally

```bash
go build .          # Build for current platform
./spine --help      # Test the binary

# Build for all platforms (requires GoReleaser)
goreleaser build --snapshot --clean
```

> Note: GoReleaser must be installed locally (`brew install goreleaser` or `go install github.com/goreleaser/goreleaser@latest`).

## Badges

Add badges once CI and releases are wired up:

- Build Status
- Go Report Card
- License

## License

MIT – see [LICENSE](LICENSE).
