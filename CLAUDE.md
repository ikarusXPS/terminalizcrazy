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
go test -race ./...                        # With race detection
go test -cover ./...                       # With coverage summary
make test-coverage                         # HTML coverage report

# CI Notes: Windows workspace tests are flaky and allowed to fail

# Lint & Format
go fmt ./...
make lint                                  # Requires golangci-lint
```

## Environment Setup

Default: Ollama with Gemma4 (local, no API key required)

Optional cloud providers:
- `GEMINI_API_KEY` - For Google Gemini
- `ANTHROPIC_API_KEY` - For Claude AI
- `OPENAI_API_KEY` - For OpenAI

Config file: `~/.terminalizcrazy/config.toml`

```toml
ai_provider = "ollama"           # ollama (default), gemini, anthropic, openai
ollama_model = "gemma4"          # or gemma4:e4b, codellama, llama3
agent_mode = "suggest"           # off, suggest (recommended), auto
```

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
| `internal/ai/` | AI clients (Gemini, Anthropic, OpenAI, Ollama) + Agent mode + Planner |
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
| `internal/crypto/` | Key management for encrypted API key storage |

### AI Integration Pattern

All AI providers implement `ai.Client` interface:
```go
type Client interface {
    Complete(ctx context.Context, req *Request) (*Response, error)
    Provider() Provider
}

// Optional streaming support (Gemini, Ollama support this)
type StreamingClient interface {
    Client
    CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error
}
```

Providers: `ollama.go` (default), `gemini.go`, `anthropic.go`, `openai.go`

**Request types** (auto-detected from user input):
- `RequestTypeCommand` - Natural language → shell command
- `RequestTypeExplain` - Error/command explanation
- `RequestTypeChat` - General conversation

**Agent Mode** (`ai.Agent` + `ai.Planner`):
- Creates multi-step task plans with verification
- Plans contain Tasks with verification criteria (exit_code, output_contains, run_command)
- Three modes: `off`, `suggest` (shows plan, waits for approval), `auto` (executes LOW-risk automatically)

### TUI Message Flow

Uses Bubble Tea's Elm architecture with async message types:
```go
// Async operation results (defined in tui.go)
aiResponseMsg     // AI completion result
cmdResultMsg      // Command execution result
streamingChunkMsg // Streaming AI response chunk
collabMessageMsg  // Collaboration message from WebSocket
themeChangedMsg   // Hot-reload theme update
```

Key components:
- `tui.Model` - Main state container with ViewMode (Chat, SessionSelect, CollabJoin, ModelSelect)
- `tui.PaneManager` - Multi-pane layout with splits, zoom, floating panes
- `tui.ConfirmState` - Risk confirmation dialog for Medium/High/Critical commands

### Collaboration Architecture

```
Host: Ctrl+S → startSharing() → collab.Server (port 8765) + collab.CollabClient
                                      ↓
                               WebSocket + ECDH key exchange
                                      ↓
Guest: Ctrl+J → joinCollab(shareCode) → collab.CollabClient
```

Message types: `MsgTypeChat`, `MsgTypeCommand`, `MsgTypeOutput`, `MsgTypeJoin`, `MsgTypeLeave`
Encryption: ECDH for key exchange, AES-256-GCM for message encryption (`internal/collab/crypto.go`)

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
| `Ctrl+A` | Toggle agent mode (off → suggest → auto) |
| `Ctrl+M` | Model selector (switch AI models) |
| `Ctrl+S` | Share session (collaboration) |
| `Ctrl+J` | Join session |

### Adding a New AI Provider

1. Create `internal/ai/{provider}.go` implementing `ai.Client` (and optionally `StreamingClient`)
2. Add provider constant to `internal/ai/ai.go`
3. Add case in `ai.NewService()` switch
4. Add config fields in `internal/config/config.go`
5. Add UI handling in `tui.NewModel()` and `loadAvailableModels()`
