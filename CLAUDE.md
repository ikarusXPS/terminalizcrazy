# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

```bash
# Build
go build -o bin/terminalizcrazy ./cmd/terminalizcrazy
make build

# Run
./bin/terminalizcrazy
make run

# Test
go test ./...                              # All tests
go test -v ./internal/ai/...               # Single package with verbose
go test -v -run TestAgentMode ./internal/ai/  # Single test by name
go test -cover ./...                       # With coverage summary
make test-coverage                         # HTML coverage report

# Lint & Format
go fmt ./...
make lint                                  # Requires golangci-lint
```

## Environment Setup

Requires one of:
- `ANTHROPIC_API_KEY` - For Claude AI
- `OPENAI_API_KEY` - For OpenAI
- `OLLAMA_ENABLED=true` + `OLLAMA_MODEL=codellama` - For local Ollama

## Architecture Overview

TerminalizCrazy is an AI-native terminal TUI built with Go and Bubble Tea (Charm.sh ecosystem).

### Core Flow

```
main.go → config.Load() → tui.Run()
                              ↓
                         tui.Model (Bubble Tea)
                              ↓
            ┌─────────────────┼─────────────────┐
            ↓                 ↓                 ↓
      ai.Service       executor.Executor   storage.Storage
      (AI providers)   (command exec)      (SQLite)
```

### Key Packages

| Package | Purpose |
|---------|---------|
| `internal/ai/` | AI clients (Anthropic, OpenAI, Ollama) + Agent mode + Planner |
| `internal/tui/` | Bubble Tea TUI with pane/tab system |
| `internal/executor/` | Command execution with risk assessment |
| `internal/storage/` | SQLite for sessions, messages, history, plans, workspaces |
| `internal/collab/` | WebSocket collaboration + E2E encryption (ECDH + AES-256-GCM) |
| `internal/workflows/` | Reusable workflow templates |
| `internal/plugins/` | Hook-based plugin system (pre_command, post_command, pre_ai, post_ai) |
| `internal/secretguard/` | API key/token detection and masking |
| `internal/project/` | Project type detection (Go, Node, Python, Rust, Java, etc.) |
| `internal/theme/` | YAML theme system with hot-reload |
| `internal/workspace/` | Workspace management with layout presets (quad, tall, wide, stack) |

### AI Integration Pattern

All AI providers implement `ai.Client` interface:
```go
type Client interface {
    Complete(ctx context.Context, req *Request) (*Response, error)
    Provider() Provider
}
```

Agent mode (`ai.Agent`) uses `ai.Planner` to create multi-step task plans that can be approved and executed. Plans contain Tasks with verification (exit_code, output_contains, run_command).

### TUI Architecture

Uses Bubble Tea's Elm architecture (Model-Update-View):
- `tui.Model` - Main state container
- `tui.PaneManager` - Multi-pane layout with splits, zoom, floating panes
- `tui.TabBar` - Tab navigation with keyboard shortcuts

### Storage Schema

SQLite database at `~/.terminalizcrazy/terminalizcrazy.db`:
- `sessions`, `messages` - Chat persistence
- `command_history` - Executed commands
- `agent_plans`, `agent_tasks` - Agent execution plans
- `workflows` - Saved workflow templates
- `workspaces` - Workspace layouts and pane states

### Plugin System

Hook-based with priority ordering. Built-in plugins:
- `SafetyPlugin` (priority 1) - Blocks dangerous commands
- `AliasPlugin` (priority 10) - Command aliases (ll→ls -la, gs→git status)
- `TimestampPlugin` - Adds timestamps to output
- `HistoryLoggerPlugin` - Command history tracking

### Key Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+E` | Execute last suggested command |
| `Ctrl+Y` | Copy command to clipboard |
| `Ctrl+T` | New tab |
| `Ctrl+W` | Close pane |
| `Ctrl+\` | Vertical split |
| `Ctrl+-` | Horizontal split |
| `Ctrl+Z` | Toggle pane zoom |
| `Alt+Arrow` | Navigate panes |
| `Ctrl+S` | Share session (collaboration) |
| `Ctrl+J` | Join session |
