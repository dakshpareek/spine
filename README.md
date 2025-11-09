# Spine – Extract Your Codebase's Architecture

[![GitHub Release](https://img.shields.io/github/v/release/dakshpareek/spine)](https://github.com/dakshpareek/spine/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go 1.23+](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org)

Maintain an AI-ready snapshot of your codebase. `spine` extracts the structural skeleton of your project into a lightweight `.spine/` workspace so coding assistants can understand architecture before they touch your source files.

## Installation

### Option 1: Pre-built Binary (Recommended)

Download from [GitHub Releases](https://github.com/dakshpareek/spine/releases):

```bash
# macOS (Apple Silicon)
tar xzf spine_v1.0.0_darwin_arm64.tar.gz
sudo mv spine /usr/local/bin/

# macOS (Intel)
tar xzf spine_v1.0.0_darwin_amd64.tar.gz
sudo mv spine /usr/local/bin/

# Linux
tar xzf spine_v1.0.0_linux_amd64.tar.gz
sudo mv spine /usr/local/bin/
```

### Option 2: From Source (requires Go 1.23+)

```bash
go install github.com/dakshpareek/spine@latest
```

Ensure `$HOME/go/bin` is in your PATH:

```bash
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Note:** When installing from source, `spine --version` will show "dev". The correct version is embedded in pre-built binaries from GitHub Releases. The actual code version is correct despite the version string.

## Updating

### Pre-built Binary

Download the latest version from [GitHub Releases](https://github.com/dakshpareek/spine/releases) and replace your existing binary:

```bash
# Assuming spine is at /usr/local/bin/spine
tar xzf spine_vX.Y.Z_<platform>.tar.gz
sudo mv spine /usr/local/bin/spine
```

### From Source

```bash
go install github.com/dakshpareek/spine@latest
```

Check your version:

```bash
spine --version
```

## Uninstalling

### Pre-built Binary

```bash
sudo rm /usr/local/bin/spine
```

### From Source

```bash
go clean -i github.com/dakshpareek/spine
```

Or manually remove:

```bash
rm $HOME/go/bin/spine
```

## Quick Start

```bash
cd /path/to/project
spine init                # bootstrap .spine and initial index
spine generate --output prompt.md
# or combine the first two steps:
# spine pipeline --output prompt.md
```

1. Paste `prompt.md` into your AI assistant and let it create/update skeletons.
2. Save the generated skeletons under `.spine/skeletons/`.
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

## Development

### Building Locally

```bash
git clone https://github.com/dakshpareek/spine.git
cd spine
go build .
./spine --help
```

### Running Tests

```bash
go test ./... -v -cover
```

### Releasing

Releases are automated via GitHub Actions. To create a release:

```bash
git tag -a v1.0.1 -m "Release v1.0.1"
git push origin v1.0.1
```

This triggers the release workflow which:
1. Runs all tests
2. Builds binaries for all platforms via [GoReleaser](https://goreleaser.com/)
3. Creates a GitHub Release with downloadable artifacts

See [.goreleaser.yml](.goreleaser.yml) for build configuration.

## License

MIT – see [LICENSE](LICENSE).
