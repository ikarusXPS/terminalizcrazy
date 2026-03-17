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
go test ./...                          # All tests
go test -v ./internal/secretguard/     # Single package
make test-coverage                     # With coverage report

# Lint & Format
go fmt ./...
make lint                              # Requires golangci-lint

# Cross-platform build
make build-all
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
| `internal/tui/` | Bubble Tea TUI with pane/tab system, zoom, floating panes, input sync |
| `internal/executor/` | Command execution with risk assessment |
| `internal/storage/` | SQLite for sessions, messages, history, plans, workspaces |
| `internal/collab/` | WebSocket collaboration + E2E encryption |
| `internal/workflows/` | Reusable workflow templates |
| `internal/plugins/` | Hook-based plugin system |
| `internal/secretguard/` | API key/token detection and masking |
| `internal/project/` | Project type detection (Go, Node, Python, etc.) |
| `internal/theme/` | YAML theme system with hot-reload (Dracula, Nord, Catppuccin, etc.) |
| `internal/workspace/` | Workspace management with layout presets (quad, tall, wide, stack) |

### AI Integration Pattern

All AI providers implement `ai.Client` interface:
```go
type Client interface {
    Complete(ctx context.Context, req *Request) (*Response, error)
    Provider() Provider
}
```

Agent mode (`ai.Agent`) uses `ai.Planner` to create multi-step task plans that can be approved and executed.

### TUI Architecture

The TUI uses Bubble Tea's Elm architecture (Model-Update-View):
- `tui.Model` - Main state container
- `tui.PaneManager` - Multi-pane layout with splits
- `tui.TabBar` - Tab navigation
- `tui/views/` - Extracted view components (ChatView, PlanView)

### Storage Schema

SQLite database at `~/.terminalizcrazy/terminalizcrazy.db`:
- `sessions` - Terminal sessions
- `messages` - Chat messages per session
- `command_history` - Executed commands
- `agent_plans` / `agent_tasks` - Agent execution plans
- `workflows` - Saved workflow templates
- `workspaces` - Workspace layouts and pane states

### Theme System

YAML-based themes with hot-reload support:
- Built-in themes: Dracula, Nord, Catppuccin Mocha, Gruvbox Dark, Tokyo Night
- Custom themes: `~/.terminalizcrazy/themes/*.yaml`
- Hot-reload: Themes reload automatically when files change

### Workspace System

Multiple workspace layouts with persistence:
- `quad` - 2x2 grid (default)
- `tall` - 1 main (60%) + 2 side stacked
- `wide` - 1 top (60%) + 2 bottom
- `stack` - 4 vertical panes

### Pane Enhancements

- **Zoom**: Toggle pane fullscreen (`Ctrl+Z`)
- **Floating**: Toggle floating mode (`Alt+F`)
- **Broadcast**: Sync input to all panes (`Ctrl+Shift+B`)

### Collaboration

WebSocket-based with optional E2E encryption:
- `collab.Server` - Local collaboration server
- `collab.CollabClient` - WebSocket client
- `collab.CryptoSession` - ECDH + AES-256-GCM encryption

### Plugin System

Hook-based architecture with priority ordering:
- Hooks: `pre_command`, `post_command`, `pre_ai`, `post_ai`, etc.
- Built-in plugins: SafetyPlugin, AliasPlugin, TimestampPlugin
- Plugins implement `PluginHandler` interface
