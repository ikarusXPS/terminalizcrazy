# TerminalizCrazy - Feature Specification

> AI-native terminal TUI with multi-provider support, agent mode, and real-time collaboration.

## Feature 1: Multi-Provider AI Integration

### Requirements
1. Support multiple AI providers (Ollama, Gemini, Anthropic, OpenAI)
2. Hot-swap between providers at runtime (Ctrl+M)
3. Streaming responses for supported providers
4. Automatic request type detection (command, explain, chat)

### API Specification
```go
type Client interface {
    Complete(ctx context.Context, req *Request) (*Response, error)
    Provider() Provider
}

type StreamingClient interface {
    Client
    CompleteStream(ctx context.Context, req *Request, handler func(StreamingResponse)) error
}
```

### Data Model
- `Request`: UserMessage, Context, Type
- `Response`: Content, Command, Explanation, Confidence, Provider
- `StreamingResponse`: Delta, Done, Command, FullText, Err

### Business Logic
- Default to Ollama (local, no API key)
- Fall back to cloud providers if Ollama unavailable
- Auto-detect request type from user input patterns

---

## Feature 2: Agent Mode (Task Planning)

### Requirements
1. Multi-step task planning with AI
2. Three modes: off, suggest, auto
3. Verification criteria per task
4. Risk assessment before execution

### API Specification
```go
type Agent struct {
    planner   *Planner
    executor  *executor.Executor
    mode      AgentMode
    maxTasks  int
}

type Plan struct {
    ID, Name, Description string
    Tasks []*Task
    Status PlanStatus
}

type Task struct {
    Command, Description string
    Verification *Verification
}
```

### Business Logic
- `suggest`: Show plan, require user approval
- `auto`: Execute LOW-risk tasks automatically
- Verify each task after execution (exit_code, output_contains)

---

## Feature 3: Real-Time Collaboration

### Requirements
1. Share terminal session with others
2. End-to-end encryption (ECDH + AES-256-GCM)
3. User presence and color coding
4. Broadcast commands and output

### API Specification
- WebSocket server on port 8765
- Share code format: `xxxx-yyyy`
- Message types: chat, command, output, join, leave

### Data Model
```go
type Room struct {
    ID, ShareCode, HostID string
    Users map[string]*User
    Messages []*Message
}

type CryptoSession struct {
    privateKey, publicKey, sharedKey
}
```

### Business Logic
- Host creates room, generates share code
- Guests join with share code
- All messages encrypted with shared key

---

## Feature 4: Multi-Pane TUI

### Requirements
1. Split panes (vertical, horizontal)
2. Layout presets (quad, tall, wide, stack, single)
3. Floating panes with drag support
4. Zoom individual panes

### API Specification
```go
type PaneManager struct {
    panes []*Pane
    layout LayoutType
}

type Pane struct {
    ID, Type, Title string
    Content, Width, Height
}
```

### Business Logic
- Ctrl+\ for vertical split
- Ctrl+- for horizontal split
- Ctrl+Z for zoom toggle
- Alt+Arrow for navigation

---

## Feature 5: Plugin System

### Requirements
1. Hook-based execution (pre/post command, pre/post AI)
2. Priority ordering
3. Built-in safety plugin
4. User plugin directory support

### API Specification
```go
type Plugin interface {
    Name() string
    Hooks() []HookType
    Priority() int
    Execute(ctx, *HookContext) (*HookResult, error)
}
```

### Built-in Plugins
| Plugin | Priority | Purpose |
|--------|----------|---------|
| SafetyPlugin | 1 | Block dangerous commands |
| AliasPlugin | 10 | Command aliases |
| TimestampPlugin | 50 | Add timestamps |

---

## Feature 6: Workspace Persistence

### Requirements
1. Save/restore workspace layouts
2. Auto-save at configurable interval
3. Multiple workspace support
4. SQLite storage

### Data Model
```go
type Workspace struct {
    ID, Name string
    Layout LayoutType
    PaneStates []*PaneState
}
```

### Business Logic
- Auto-save every 60 seconds (configurable)
- Restore last workspace on startup
- Max 10 workspaces (configurable)

---

## Feature 7: Secret Guard

### Requirements
1. Detect secrets in output (API keys, tokens)
2. Mask detected secrets
3. Support multiple secret patterns

### Detection Patterns
- AWS keys: `AKIA[0-9A-Z]{16}`
- GitHub tokens: `gh[ps]_[A-Za-z0-9]{36}`
- JWT: `eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+`
- Generic API keys: `api[_-]?key[=:][A-Za-z0-9]{20,}`

---

## Feature 8: Theme System

### Requirements
1. YAML-based theme files
2. Hot-reload on file change
3. Built-in themes (default, dracula, monokai, nord)
4. User theme directory

### Data Model
```go
type Theme struct {
    Name, Description, Author string
    Colors ColorPalette
}

type ColorPalette struct {
    Primary, Secondary, Background, Foreground string
    Success, Warning, Error, Info, Muted, Border string
}
```

---

## Non-Functional Requirements

### Performance
- Startup time < 500ms
- AI response streaming with < 100ms first token
- Smooth 60fps TUI rendering

### Security
- API keys encrypted at rest
- E2E encryption for collaboration
- No secrets in logs or output

### Compatibility
- Linux, macOS, Windows
- Terminal: iTerm2, Terminal.app, Windows Terminal, Alacritty
- Go 1.25+
