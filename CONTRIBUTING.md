# Contributing to TerminalizCrazy

Thank you for your interest in contributing to TerminalizCrazy! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/ikarusXPS/terminalizcrazy/issues)
2. If not, create a new issue using the bug report template
3. Include:
   - Clear description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Your environment (OS, Go version, terminal)

### Suggesting Features

1. Check existing issues and discussions for similar ideas
2. Create a new issue using the feature request template
3. Describe the use case and expected behavior

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Commit with conventional commits: `git commit -m 'feat: add amazing feature'`
7. Push: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.21+
- Git
- Make (optional but recommended)
- golangci-lint (for linting)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/terminalizcrazy.git
cd terminalizcrazy

# Install dependencies
go mod tidy

# Build
make build

# Run tests
make test

# Run with your API key
export ANTHROPIC_API_KEY="your-key"
./bin/terminalizcrazy
```

### Project Structure

```
terminalizcrazy/
├── cmd/terminalizcrazy/    # Entry point
├── internal/
│   ├── ai/                 # AI providers (Anthropic, OpenAI, Ollama)
│   ├── collab/             # Real-time collaboration
│   ├── config/             # Configuration management
│   ├── executor/           # Command execution & risk assessment
│   ├── plugins/            # Plugin system
│   ├── project/            # Project detection
│   ├── secretguard/        # Secret masking
│   ├── storage/            # SQLite persistence
│   ├── theme/              # Theme system
│   ├── tui/                # Terminal UI (Bubble Tea)
│   ├── workflows/          # Workflow templates
│   └── workspace/          # Workspace management
├── docs/                   # Documentation
└── config.toml.example     # Example configuration
```

## Coding Standards

### Go Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `golangci-lint` before committing

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: fix a bug
docs: documentation changes
test: add or update tests
refactor: code refactoring
chore: maintenance tasks
perf: performance improvements
ci: CI/CD changes
```

### Testing

- Write tests for new functionality
- Maintain or improve code coverage
- Run `make test` before submitting PR

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/ai/...
```

## Pull Request Guidelines

### Before Submitting

- [ ] Tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventions

### PR Description

Include:
- What changes were made
- Why the changes were made
- How to test the changes
- Screenshots for UI changes

### Review Process

1. Automated checks must pass
2. At least one maintainer review required
3. Address feedback promptly
4. Squash commits if requested

## Areas for Contribution

### Good First Issues

Look for issues labeled `good-first-issue` - these are suitable for newcomers.

### Priority Areas

- Performance improvements (GPU acceleration)
- New AI provider integrations
- Plugin development
- Documentation improvements
- Test coverage
- Accessibility improvements

## Getting Help

- Open a [Discussion](https://github.com/ikarusXPS/terminalizcrazy/discussions)
- Check existing documentation in `docs/`
- Ask in PR comments

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to TerminalizCrazy!
