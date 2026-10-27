# Contributing to cjk-press

Thank you for your interest in contributing to cjk-press! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

This project adheres to a code of conduct that all contributors must follow. By participating, you are expected to uphold this code. Please report unacceptable behavior to [maintainer email].

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- Environment details (OS, version, etc.)
- Any relevant logs or screenshots

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- Use a clear and descriptive title
- Provide a detailed description of the suggested enhancement
- Explain why this enhancement would be useful
- List any examples or use cases

### Pull Requests

1. **Fork the repository** and create your branch from `main`.
2. **Make your changes** following the coding standards below.
3. **Add tests** for new functionality or bug fixes.
4. **Ensure all tests pass** with `go test ./...`.
5. **Run linters** with `golangci-lint run`.
6. **Update documentation** if needed.
7. **Commit your changes** with clear, descriptive commit messages.
8. **Push to your fork** and submit a pull request.

### Development Setup

1. Clone the repository:
```bash
git clone https://github.com/atetov/cjk-press.git
cd cjk-press
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test -v ./...
```

4. Run linters:
```bash
golangci-lint run
```

5. Build the project:
```bash
go build ./cmd/cjk-press
```

## Coding Standards

### Go Style Guide

Follow the standard Go conventions as described in [Effective Go](https://golang.org/doc/effective_go) and the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Code Formatting

- Use `gofmt` to format your code
- Run `goimports` to manage imports
- Keep line length reasonable (preferably under 120 characters)

### Documentation

- All exported functions, types, and constants must have GoDoc comments
- Comments should be complete sentences
- Use present tense in function descriptions
- Include examples in GoDoc comments where appropriate

### Testing

- Write unit tests for all new functionality
- Aim for high code coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for test changes
- `refactor:` for code refactoring
- `perf:` for performance improvements
- `ci:` for CI/CD changes

Example:
```
feat: add support for streaming compression via stdin/stdout

Implement streaming compression to handle large files without
loading them entirely into memory. This improves memory efficiency
for files larger than available RAM.
```

## Project Structure

```
cjk-press/
├── cmd/
│   └── cjk-press/       # CLI application entry point
├── pkg/
│   ├── compress/        # Compression algorithms and I/O
│   ├── dictionary/      # Dictionary implementation
│   └── segment/         # Text segmentation logic
├── .github/             # GitHub workflows and templates
├── docs/                # Additional documentation
└── tests/               # Integration tests and benchmarks
```

## Adding New Features

When adding new features:

1. Update the ROADMAP.md if the feature is significant
2. Add appropriate tests
3. Update the README.md if the feature is user-facing
4. Consider backward compatibility
5. Update version in go.mod if breaking changes are introduced

## Performance Considerations

- Profile code before and after optimizations
- Use Go's built-in benchmarking tools
- Consider memory usage alongside CPU performance
- Test with realistic CJK text datasets

## Release Process

Releases are managed through semantic versioning and automated via GitHub Actions. The maintainers will:

1. Update the version in go.mod
2. Update CHANGELOG.md
3. Create a git tag
4. GitHub Actions will automatically build and publish releases

## Getting Help

If you need help:

- Check existing documentation and issues
- Ask questions in GitHub Discussions
- Contact maintainers directly for sensitive issues

## License

By contributing to cjk-press, you agree that your contributions will be licensed under the MIT License.
