# Backend Codemap

> Freshness: 2026-04-17 | Auto-generated from source analysis

## Package Structure

```
internal/
├── ai/                 # AI provider integrations
│   ├── ai.go          # Service, Client interface, Request/Response
│   ├── agent.go       # AgentMode, autonomous execution (557 lines)
│   ├── planner.go     # Task planning, verification (515 lines)
│   ├── ollama.go      # Ollama client (446 lines)
│   ├── gemini.go      # Google Gemini client
│   ├── anthropic.go   # Anthropic Claude client
│   └── openai.go      # OpenAI GPT client
│
├── tui/                # Terminal UI (Bubble Tea)
│   ├── tui.go         # Main Model, Update, View (1924 lines)
│   ├── pane_manager.go # Multi-pane layout (676 lines)
│   ├── tab_bar.go     # Tab navigation (480 lines)
│   ├── pane.go        # Pane component (425 lines)
│   ├── floating_pane.go # Floating windows
│   ├── styles.go      # Lipgloss styles
│   ├── sync_input.go  # Synchronized input
│   ├── zoom.go        # Pane zoom
│   └── views/         # Sub-views
│       ├── chat_view.go
│       └── plan_view.go
│
├── storage/           # SQLite persistence
│   └── storage.go     # All CRUD operations (1044 lines)
│
├── collab/            # Real-time collaboration
│   ├── server.go      # WebSocket server (582 lines)
│   ├── client.go      # WebSocket client (449 lines)
│   ├── crypto.go      # ECDH + AES-256-GCM
│   └── types.go       # Message types
│
├── config/            # Configuration
│   └── config.go      # Viper-based config (437 lines)
│
├── executor/          # Command execution
│   └── executor.go    # Execute + risk assessment
│
├── plugins/           # Plugin system
│   ├── plugin.go      # Manager, interfaces (525 lines)
│   └── builtin.go     # Built-in plugins
│
├── workflows/         # Workflow templates
│   ├── workflow.go    # Engine, types (421 lines)
│   └── templates.go   # Built-in workflows
│
├── workspace/         # Workspace management
│   ├── manager.go     # Workspace CRUD (470 lines)
│   ├── layout.go      # Layout calculations
│   ├── persistence.go # SQLite storage
│   ├── workspace.go   # Types
│   └── errors.go      # Error definitions
│
├── theme/             # Theme system
│   ├── theme.go       # Theme struct, validation
│   ├── loader.go      # YAML loader
│   └── builtin.go     # Built-in themes
│
├── project/           # Project detection
│   └── project.go     # Detect 11+ project types (518 lines)
│
├── secretguard/       # Secret masking
│   └── secretguard.go # Detect & mask secrets
│
├── clipboard/         # Clipboard access
│   └── clipboard.go   # Cross-platform clipboard
│
└── crypto/            # Key management
    └── crypto.go      # API key encryption
```

## AI Package Details

### Providers

| Provider | File | Streaming | Notes |
|----------|------|-----------|-------|
| Ollama | ollama.go | Yes | Default, local inference |
| Gemini | gemini.go | Yes | Google AI |
| Anthropic | anthropic.go | No | Claude models |
| OpenAI | openai.go | No | GPT models |

### Request Types

```go
RequestTypeCommand  // Natural language → shell command
RequestTypeExplain  // Error/command explanation
RequestTypeChat     // General conversation
```

### Agent Modes

```go
AgentModeOff     // Single commands only
AgentModeSuggest // Plans with approval (default)
AgentModeAuto    // Auto-execute LOW risk
```

## TUI Package Details

### Message Types (Async)

```go
aiResponseMsg      // AI completion result
cmdResultMsg       // Command execution result
streamingChunkMsg  // Streaming AI chunk
collabMessageMsg   // Collab WebSocket message
themeChangedMsg    // Hot-reload theme
modelsLoadedMsg    // Available models list
sessionsLoadedMsg  // Session list
sessionRestoredMsg // Restored session
historyLoadedMsg   // Command history
```

### View Modes

```go
ViewChat          // Main chat view
ViewSessionSelect // Session picker
ViewCollabJoin    // Join collab room
ViewModelSelect   // Model picker
```

### Pane Types

```go
PaneTypeChat    // Chat/input pane
PaneTypeOutput  // Command output
PaneTypeHistory // Command history
PaneTypeHelp    // Help/keybindings
```

## Collab Package Details

### Message Types

```go
MsgTypeChat     // Chat message
MsgTypeCommand  // Command suggestion
MsgTypeOutput   // Command output
MsgTypeJoin     // User joined
MsgTypeLeave    // User left
MsgTypeUserList // User list update
```

### Encryption Flow

```
1. NewCryptoSession() → Generate ECDH keypair
2. GetKeyExchangePayload() → Send public key
3. SetPeerKey() → Derive shared secret
4. Encrypt()/Decrypt() → AES-256-GCM
```

## Plugin Package Details

### Hook Types

```go
HookPreCommand   // Before command execution
HookPostCommand  // After command execution
HookPreAI        // Before AI request
HookPostAI       // After AI response
HookPreOutput    // Before output display
HookPostOutput   // After output display
HookStartup      // On application start
HookShutdown     // On application exit
```

### Built-in Plugins

| Plugin | Priority | Purpose |
|--------|----------|---------|
| SafetyPlugin | 1 | Block dangerous commands |
| AliasPlugin | 10 | Command aliases |
| TimestampPlugin | 50 | Add timestamps |
| HistoryLoggerPlugin | 100 | Log commands |

## Workspace Package Details

### Layout Types

```go
LayoutQuad   // 2x2 grid
LayoutTall   // Main + 2 side
LayoutWide   // Main + 2 bottom
LayoutStack  // 4 vertical
LayoutSingle // Single pane
```

## Project Package Details

### Detected Types

| Type | Detection File |
|------|---------------|
| Go | go.mod |
| Node | package.json |
| Python | requirements.txt, pyproject.toml |
| Rust | Cargo.toml |
| Java | pom.xml, build.gradle |
| Ruby | Gemfile |
| PHP | composer.json |
| .NET | *.csproj, *.sln |
| Docker | Dockerfile |
| Terraform | *.tf |
| Kubernetes | k8s/, kubernetes/ |
