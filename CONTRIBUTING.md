# Contributing to Code Context CLI

Thanks for your interest in improving `ctx`! This document explains how to get started.

## Development Setup

1. Fork and clone the repository.
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Ensure Go 1.21+ is available.

## Workflow

1. Create a feature branch off `main`.
2. Run `go test ./... -cover` before opening a pull request.
3. Use `ctx validate --fix --strict` to confirm index integrity if you touch the `.code-context/` workflow.

## Coding Guidelines

- Keep functions small and cohesive.
- Add unit tests for new functionality.
- Add integration coverage if behaviour spans multiple packages.
- Use `gofmt`/`goimports` (run `gofmt -w` before committing).

## Documentation

- Update `README.md` and `docs/` when adding new commands or flags.
- Include examples to help users understand expected usage.

## Commit Messages

- Use clear, present-tense summaries (`Add generate command`).
- Reference related issues when applicable.

## Releases

- Update `.goreleaser.yml` if you change build targets.
- Create annotated tags (`git tag -a vX.Y.Z -m "Release vX.Y.Z"`) before running GoReleaser.

## Support

Open a GitHub issue with reproduction steps, or start a discussion thread. Be sure to include:

- `ctx` version (`ctx --version`)
- Go version (`go version`)
- Operating system
- Command output and relevant index snippets

Thanks again for contributing!
