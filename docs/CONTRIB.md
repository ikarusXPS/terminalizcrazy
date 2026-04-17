# Contributing Guide

> Generated from source of truth: Makefile, go.mod, config.toml.example

## Development Workflow

### Prerequisites

- Go 1.25.0+ (see go.mod)
- golangci-lint (for linting)
- Ollama (for local AI testing) or cloud API key

### Quick Start

```bash
# Clone and setup
git clone https://github.com/ikarusXPS/terminalizcrazy.git
cd terminalizcrazy
go mod tidy

# Build and run
make build
./bin/terminalizcrazy
```

## Available Scripts

| Command | Description |
|---------|-------------|
| `make build` | Build the application to `bin/terminalizcrazy` |
| `make run` | Run without building binary |
| `make test` | Run all tests with verbose output |
| `make test-coverage` | Generate HTML coverage report |
| `make fmt` | Format code with gofmt |
| `make lint` | Run golangci-lint |
| `make tidy` | Tidy go.mod dependencies |
| `make clean` | Remove build artifacts |
| `make install` | Install to $GOPATH/bin |
| `make build-all` | Cross-compile for all platforms |
| `make help` | Show all available targets |

### Build Flags

The build embeds version info via ldflags:
- `main.version` - Git tag or "dev"
- `main.commit` - Short commit hash

## Environment Setup

### Local Development (Default)

No API key required with Ollama:

```bash
# Pull and start Ollama
ollama pull gemma4
ollama serve

# Run TerminalizCrazy (defaults to Ollama)
./bin/terminalizcrazy
```

### Cloud Providers (Optional)

Set environment variables for cloud AI:

| Variable | Provider | Format |
|----------|----------|--------|
| `GEMINI_API_KEY` | Google Gemini | `AIzaSy...` |
| `ANTHROPIC_API_KEY` | Anthropic Claude | `sk-ant-api03-...` |
| `OPENAI_API_KEY` | OpenAI GPT | `sk-...` |
| `AI_PROVIDER` | Override default | `ollama`, `gemini`, `anthropic`, `openai` |

### Debug Mode

```bash
export DEBUG=true
export LOG_LEVEL=debug
./bin/terminalizcrazy
```

## Testing Procedures

### Run All Tests

```bash
make test
# or
go test -v ./...
```

### Run Specific Package

```bash
go test -v ./internal/ai/...
go test -v ./internal/tui/...
```

### Run Single Test

```bash
go test -v -run TestAgentMode ./internal/ai/
```

### With Race Detection

```bash
go test -race ./...
```

### Coverage Report

```bash
make test-coverage
# Opens coverage.html
```

### CI Notes

- Windows workspace tests are flaky and allowed to fail
- Coverage target: maintain existing percentage

## Code Standards

### Formatting

```bash
make fmt
# or
go fmt ./...
```

### Linting

```bash
make lint
# Requires: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

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

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `charmbracelet/bubbletea` | TUI framework |
| `charmbracelet/lipgloss` | TUI styling |
| `google/generative-ai-go` | Gemini AI client |
| `liushuangls/go-anthropic` | Claude AI client |
| `sashabaranov/go-openai` | OpenAI client |
| `gorilla/websocket` | Collaboration WebSocket |
| `spf13/viper` | Configuration management |
| `modernc.org/sqlite` | Pure-Go SQLite |

## Project Structure

```
internal/
├── ai/           # AI providers + Agent + Planner
├── tui/          # Bubble Tea UI components
├── executor/     # Command execution + risk assessment
├── storage/      # SQLite persistence
├── collab/       # Real-time collaboration + E2E crypto
├── config/       # Viper configuration
├── plugins/      # Hook-based plugin system
├── secretguard/  # Secret detection and masking
├── project/      # Project type detection
├── theme/        # YAML theme system
├── workspace/    # Layout management
├── workflows/    # Workflow templates
└── crypto/       # Key management
```

## Adding New Features

### New AI Provider

1. Create `internal/ai/{provider}.go` implementing `ai.Client`
2. Add provider constant to `internal/ai/ai.go`
3. Add case in `ai.NewService()` switch
4. Add config fields in `internal/config/config.go`
5. Add UI handling in `tui.NewModel()`

### New Plugin

1. Implement `plugins.Plugin` interface
2. Register in `plugins/builtin.go` or user plugin directory
3. Set appropriate priority (lower = runs first)

### New Theme

1. Create `~/.terminalizcrazy/themes/{name}.yaml`
2. Follow schema in `internal/theme/theme.go`
3. Hot-reload enabled by default
