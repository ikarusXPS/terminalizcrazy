# Architecture Codemap

> Freshness: 2026-04-17 | Auto-generated from source analysis

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      cmd/terminalizcrazy                     │
│                         (entry point)                        │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      internal/config                         │
│              (viper config, env vars, TOML)                  │
│                          ↓ uses crypto                       │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        internal/tui                          │
│                    (Bubble Tea main loop)                    │
│                        1924 lines                            │
└───┬─────────┬─────────┬─────────┬─────────┬─────────┬───────┘
    │         │         │         │         │         │
    ▼         ▼         ▼         ▼         ▼         ▼
┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐
│  ai   │ │storage│ │executor│ │collab │ │project│ │ theme │
└───┬───┘ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘
    │
    ▼
┌───────────────┐
│   executor    │
│ (risk assess) │
└───────────────┘
```

## Dependency Graph

```
tui ──────────────────────────────────────────────┐
 ├─→ ai ─→ executor                               │
 ├─→ clipboard                                    │
 ├─→ collab                                       │
 ├─→ config ─→ crypto                             │
 ├─→ executor                                     │
 ├─→ project                                      │
 ├─→ secretguard                                  │
 ├─→ storage                                      │
 └─→ theme                                        │
                                                  │
workflows ─→ executor                             │
                                                  │
workspace (standalone, no internal deps)          │
plugins (standalone, no internal deps)            │
```

## Package Summary

| Package | Lines | Purpose | Dependencies |
|---------|-------|---------|--------------|
| tui | 1924 | Main UI loop, Bubble Tea | ai, clipboard, collab, config, executor, project, secretguard, storage, theme |
| storage | 1044 | SQLite persistence | - |
| pane_manager | 676 | Multi-pane layout | (part of tui) |
| collab/server | 582 | WebSocket collab server | - |
| ai/agent | 557 | Autonomous task execution | executor |
| plugins | 525 | Hook-based plugin system | - |
| project | 518 | Project type detection | - |
| ai/planner | 515 | Multi-step task planning | - |
| workspace | 470 | Layout management | - |

## External Dependencies

### Core Framework
- `charmbracelet/bubbletea` - TUI framework
- `charmbracelet/lipgloss` - Styling
- `charmbracelet/bubbles` - UI components

### AI Providers
- `google/generative-ai-go` - Gemini
- `liushuangls/go-anthropic` - Claude
- `sashabaranov/go-openai` - OpenAI
- Ollama (HTTP API, no SDK)

### Infrastructure
- `modernc.org/sqlite` - Pure-Go SQLite
- `gorilla/websocket` - WebSocket
- `spf13/viper` - Configuration
- `golang.org/x/crypto` - Encryption

## Data Flow

```
User Input
    │
    ▼
┌─────────────────┐
│   tui.Model     │ ←──── tea.Msg (async)
│   (Update)      │
└────────┬────────┘
         │
    ┌────┴────┬────────────┬─────────────┐
    ▼         ▼            ▼             ▼
┌───────┐ ┌───────┐ ┌──────────┐ ┌───────────┐
│  AI   │ │Executor│ │ Storage  │ │  Collab   │
│Service│ │        │ │ (SQLite) │ │(WebSocket)│
└───┬───┘ └───┬───┘ └────┬─────┘ └─────┬─────┘
    │         │          │             │
    ▼         ▼          ▼             ▼
 Response   Result    Persist      Broadcast
    │         │          │             │
    └─────────┴──────────┴─────────────┘
                    │
                    ▼
              tui.View()
                    │
                    ▼
               Terminal
```

## Key Interfaces

```go
// AI Provider interface (internal/ai/ai.go)
type Client interface {
    Complete(ctx, *Request) (*Response, error)
    Provider() Provider
}

// Plugin interface (internal/plugins/plugin.go)
type Plugin interface {
    Name() string
    Type() PluginType
    Hooks() []HookType
    Execute(ctx, *HookContext) (*HookResult, error)
}

// Storage interface (internal/workspace/manager.go)
type Storage interface {
    SaveWorkspace(*Workspace) error
    LoadWorkspace(id string) (*Workspace, error)
    ListWorkspaces() ([]*Workspace, error)
    DeleteWorkspace(id string) error
}
```

## File Count by Package

| Package | Source Files | Test Files |
|---------|-------------|------------|
| ai | 8 | 5 |
| tui | 10 | 6 |
| collab | 4 | 3 |
| storage | 1 | 1 |
| config | 1 | 1 |
| workspace | 5 | 1 |
| plugins | 2 | 2 |
| workflows | 2 | 2 |
| theme | 3 | 1 |
| project | 1 | 1 |
| executor | 1 | 1 |
| secretguard | 1 | 1 |
| clipboard | 1 | 1 |
| crypto | 1 | 1 |
