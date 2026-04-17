# TerminalizCrazy - Feature Specification

> AI-native terminal TUI with multi-provider support, agent mode, and real-time collaboration.

---

## Feature 1: Multi-Provider AI Integration

### User Stories
- As a user, I want to use local AI (Ollama) without API keys so I can work offline
- As a user, I want to switch AI providers at runtime so I can compare responses
- As a user, I want streaming responses so I see output progressively
- As a user, I want the AI to understand my intent (command vs explanation vs chat)

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| AI-1 | Support Ollama as default provider (no API key required) | MUST |
| AI-2 | Support Gemini, Anthropic, OpenAI as cloud alternatives | MUST |
| AI-3 | Hot-swap providers at runtime via Ctrl+M | MUST |
| AI-4 | Streaming responses for Ollama and Gemini | MUST |
| AI-5 | Auto-detect request type from input patterns | SHOULD |
| AI-6 | Display provider name and model in UI | SHOULD |
| AI-7 | Graceful fallback when provider unavailable | SHOULD |
| AI-8 | Request timeout configurable (default 30s) | COULD |

### Acceptance Criteria

- [ ] User can start app with only Ollama installed, no env vars
- [ ] Ctrl+M opens model selector showing available providers
- [ ] Selecting a provider switches immediately, next request uses it
- [ ] Streaming shows character-by-character output in real-time
- [ ] "how do I list files" → detected as command request
- [ ] "what does this error mean: ..." → detected as explain request
- [ ] "tell me about Go channels" → detected as chat request
- [ ] If Ollama is down and no cloud keys, show helpful error message

### Edge Cases & Error Handling

| Scenario | Expected Behavior |
|----------|-------------------|
| Ollama not running | Show "Ollama not available. Start with `ollama serve` or configure cloud provider" |
| Invalid API key | Show "API key invalid for {provider}. Check your configuration" |
| Network timeout | Show "Request timed out after {n}s. Check your connection" |
| Rate limit hit | Show "Rate limit reached. Try again in {n} seconds" |
| Empty response | Show "AI returned empty response. Try rephrasing your request" |
| Streaming interrupted | Display partial response with "[interrupted]" marker |

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
```go
type Request struct {
    UserMessage string
    Context     *RequestContext
    Type        RequestType  // command, explain, chat
}

type RequestContext struct {
    CurrentDir       string   // e.g., /home/user/project
    OS               string   // linux, darwin, windows
    Shell            string   // bash, zsh, powershell
    RecentHistory    []string // last 5 commands
    ProjectName      string   // detected project name
    ProjectType      string   // go, node, python, etc.
    ProjectFramework string   // e.g., "next", "gin", "django"
}

type Response struct {
    Content     string   // full response text
    Command     string   // extracted command (if any)
    Explanation string   // explanation text
    Confidence  float64  // 0.0-1.0 confidence score
    Provider    Provider // which provider responded
}
```

### Testing Requirements
- [ ] Unit tests for each provider client
- [ ] Unit tests for request type detection
- [ ] Integration test: Ollama completion (requires Ollama running)
- [ ] Mock tests for cloud providers
- [ ] Timeout handling test
- [ ] Streaming interruption test

---

## Feature 2: Agent Mode (Task Planning)

### User Stories
- As a user, I want multi-step task execution so I can automate complex workflows
- As a user, I want to review plans before execution so I stay in control
- As a user, I want automatic execution of safe commands so I can work faster
- As a user, I want verification after each step so I know tasks succeeded

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| AG-1 | Three modes: off, suggest, auto | MUST |
| AG-2 | AI generates multi-step plans from natural language | MUST |
| AG-3 | Each task has verification criteria | MUST |
| AG-4 | Risk assessment (LOW/MEDIUM/HIGH/CRITICAL) per task | MUST |
| AG-5 | User approval required for MEDIUM+ risk in suggest mode | MUST |
| AG-6 | Auto mode only executes LOW risk automatically | MUST |
| AG-7 | Plan persistence in SQLite | SHOULD |
| AG-8 | Cancel/abort running plan | SHOULD |
| AG-9 | Retry failed task with modifications | COULD |
| AG-10 | Max 10 tasks per plan (configurable) | SHOULD |

### Acceptance Criteria

- [ ] "set up a new Go module with tests" → generates 4-5 step plan
- [ ] Plan shows: task description, command, risk level, verification
- [ ] In suggest mode, user sees plan and presses Enter to approve
- [ ] In auto mode, LOW risk tasks execute without prompt
- [ ] MEDIUM/HIGH/CRITICAL always require confirmation
- [ ] Failed verification stops plan and shows error
- [ ] Ctrl+C cancels running plan immediately
- [ ] Plan history viewable in session

### Risk Level Definitions

| Level | Examples | Auto-Execute |
|-------|----------|--------------|
| LOW | `ls`, `cat`, `go build`, `npm test` | Yes (auto mode) |
| MEDIUM | `mkdir`, `touch`, `go mod init` | No |
| HIGH | `rm file`, `git commit`, `npm install` | No |
| CRITICAL | `rm -rf`, `sudo`, `chmod 777`, `DROP TABLE` | Never |

### Verification Types

| Type | Description | Example |
|------|-------------|---------|
| exit_code | Check command exit code | `exit_code: 0` |
| output_contains | Check output for string | `output_contains: "ok"` |
| output_not_contains | Ensure output lacks string | `output_not_contains: "error"` |
| file_exists | Check file was created | `file_exists: "go.mod"` |
| run_command | Run verification command | `run_command: "test -f go.mod"` |

### API Specification
```go
type Agent struct {
    planner   *Planner
    executor  *executor.Executor
    mode      AgentMode
    maxTasks  int
    callbacks AgentCallbacks
}

type AgentCallbacks struct {
    OnPlanReady  func(*Plan)
    OnTaskStart  func(*Task)
    OnTaskDone   func(*Task, *Result, error)
    OnPlanDone   func(*Plan, error)
}

type Plan struct {
    ID          string
    Name        string
    Description string
    Tasks       []*Task
    Status      PlanStatus  // pending, approved, running, completed, failed, cancelled
    CreatedAt   time.Time
}

type Task struct {
    ID           string
    Sequence     int
    Command      string
    Description  string
    RiskLevel    RiskLevel
    Status       TaskStatus
    Output       string
    Error        string
    Verification *Verification
}
```

### Edge Cases & Error Handling

| Scenario | Expected Behavior |
|----------|-------------------|
| AI generates invalid command | Show parse error, ask user to rephrase |
| Task fails verification | Stop plan, show which verification failed |
| User cancels mid-plan | Abort cleanly, show completed vs skipped tasks |
| Network loss during plan | Pause plan, offer resume when reconnected |
| Task exceeds timeout (60s) | Kill task, mark as timeout failure |
| Circular dependency detected | Reject plan, show error |

### Testing Requirements
- [ ] Unit tests for plan generation
- [ ] Unit tests for risk assessment
- [ ] Unit tests for each verification type
- [ ] Integration test: full plan execution
- [ ] Cancellation test
- [ ] Timeout test

---

## Feature 3: Real-Time Collaboration

### User Stories
- As a user, I want to share my terminal session with colleagues
- As a user, I want end-to-end encryption so my data is private
- As a user, I want to see who's in the session
- As a user, I want to chat alongside the terminal

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| COL-1 | Host creates shareable session with code | MUST |
| COL-2 | Guest joins with share code | MUST |
| COL-3 | E2E encryption (ECDH + AES-256-GCM) | MUST |
| COL-4 | Broadcast commands and output to all | MUST |
| COL-5 | User presence with color coding | MUST |
| COL-6 | In-session chat | SHOULD |
| COL-7 | Max 5 users per session | SHOULD |
| COL-8 | Host can kick users | COULD |
| COL-9 | Session timeout after 1 hour idle | COULD |

### Acceptance Criteria

- [ ] Ctrl+S starts sharing, displays share code (format: xxxx-yyyy)
- [ ] Ctrl+J prompts for share code, joins session
- [ ] All participants see same terminal output
- [ ] Commands from any user broadcast to all
- [ ] User list shows names with unique colors
- [ ] Messages encrypted before transmission
- [ ] Disconnected user removed from list after 10s
- [ ] Host leaving ends session for all

### Protocol Specification

```
Connection Flow:
1. Host: Ctrl+S → Start WebSocket server on :8765
2. Host: Generate share code, display to user
3. Guest: Ctrl+J → Enter share code
4. Guest: Connect to host WebSocket
5. Both: ECDH key exchange
6. Both: Derive shared AES-256-GCM key
7. All messages encrypted with shared key
```

### Message Types

| Type | Direction | Payload |
|------|-----------|---------|
| join | Guest→Host | `{user_name, public_key}` |
| welcome | Host→Guest | `{room_id, users[], host_public_key}` |
| user_joined | Host→All | `{user_id, user_name, color}` |
| user_left | Host→All | `{user_id}` |
| chat | Any→All | `{user_id, content, timestamp}` |
| command | Any→All | `{user_id, command}` |
| output | Host→All | `{command_id, output, exit_code}` |

### Edge Cases & Error Handling

| Scenario | Expected Behavior |
|----------|-------------------|
| Invalid share code | Show "Invalid share code. Check and try again" |
| Host unreachable | Show "Cannot connect to host. They may have stopped sharing" |
| Connection lost | Attempt reconnect 3x, then show disconnected state |
| Key exchange fails | Abort connection, show "Encryption setup failed" |
| Max users reached | Reject join, show "Session is full (5/5)" |
| Malformed message | Log and ignore, don't crash |

### Testing Requirements
- [ ] Unit tests for crypto (ECDH, AES-GCM)
- [ ] Unit tests for message serialization
- [ ] Integration test: 2-user session
- [ ] Reconnection test
- [ ] Encryption verification test

---

## Feature 4: Multi-Pane TUI

### User Stories
- As a user, I want split panes so I can see multiple outputs
- As a user, I want preset layouts for common workflows
- As a user, I want to zoom a pane temporarily
- As a user, I want floating panes for reference content

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| TUI-1 | Vertical split (Ctrl+\) | MUST |
| TUI-2 | Horizontal split (Ctrl+-) | MUST |
| TUI-3 | Close pane (Ctrl+W) | MUST |
| TUI-4 | Navigate panes (Alt+Arrow) | MUST |
| TUI-5 | Zoom pane (Ctrl+Z) | MUST |
| TUI-6 | Layout presets: quad, tall, wide, stack, single | SHOULD |
| TUI-7 | Floating panes with drag | COULD |
| TUI-8 | Pane resize with mouse | COULD |
| TUI-9 | Minimum pane size enforcement | SHOULD |

### Acceptance Criteria

- [ ] Ctrl+\ splits active pane vertically (left/right)
- [ ] Ctrl+- splits active pane horizontally (top/bottom)
- [ ] Alt+Arrow moves focus to adjacent pane
- [ ] Ctrl+Z zooms active pane to full screen, again to restore
- [ ] Ctrl+W closes active pane (unless it's the last one)
- [ ] Pane titles show type (Chat, Output, History, Help)
- [ ] Minimum pane size: 20 chars wide, 5 lines tall

### Layout Presets

| Layout | Description | Pane Count |
|--------|-------------|------------|
| single | One full-screen pane | 1 |
| tall | Large left + 2 stacked right | 3 |
| wide | Large top + 2 side-by-side bottom | 3 |
| quad | 2x2 grid | 4 |
| stack | 4 vertical strips | 4 |

### Pane Types

| Type | Purpose | Default Content |
|------|---------|-----------------|
| Chat | AI conversation | Input + messages |
| Output | Command output | Last command result |
| History | Command history | Scrollable list |
| Help | Keybindings | Static help text |

### Testing Requirements
- [ ] Unit tests for layout calculations
- [ ] Unit tests for focus navigation
- [ ] Split/close cycle test
- [ ] Minimum size enforcement test
- [ ] Zoom toggle test

---

## Feature 5: Plugin System

### User Stories
- As a user, I want dangerous commands blocked automatically
- As a user, I want command aliases (ll → ls -la)
- As a developer, I want to extend functionality with plugins

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| PLG-1 | Hook types: pre_command, post_command, pre_ai, post_ai | MUST |
| PLG-2 | Priority ordering (lower = runs first) | MUST |
| PLG-3 | Built-in SafetyPlugin blocks dangerous commands | MUST |
| PLG-4 | Built-in AliasPlugin for command shortcuts | SHOULD |
| PLG-5 | User plugin directory (~/.terminalizcrazy/plugins/) | SHOULD |
| PLG-6 | Plugin can block, modify, or passthrough | MUST |
| PLG-7 | Plugin configuration via manifest.yaml | COULD |
| PLG-8 | Plugin enable/disable without restart | COULD |

### Acceptance Criteria

- [ ] `rm -rf /` is blocked with warning message
- [ ] `sudo rm -rf ~` is blocked with warning message
- [ ] `ll` expands to `ls -la`
- [ ] `gs` expands to `git status`
- [ ] Plugins execute in priority order (1, 10, 50, 100)
- [ ] Plugin can modify command before execution
- [ ] Plugin can suppress output display

### Hook Flow

```
User Input
    │
    ▼
┌─────────────────┐
│  pre_command    │ ← SafetyPlugin (block dangerous)
│  hooks          │ ← AliasPlugin (expand aliases)
└────────┬────────┘
         │
         ▼
    [Execute Command]
         │
         ▼
┌─────────────────┐
│  post_command   │ ← TimestampPlugin (add timing)
│  hooks          │ ← HistoryPlugin (log to db)
└────────┬────────┘
         │
         ▼
    [Display Output]
```

### Blocked Command Patterns (SafetyPlugin)

| Pattern | Reason |
|---------|--------|
| `rm -rf /` | Deletes entire filesystem |
| `rm -rf ~` | Deletes home directory |
| `rm -rf *` | Deletes current directory |
| `:(){ :\|:& };:` | Fork bomb |
| `mkfs.` | Formats filesystem |
| `dd if=/dev/zero of=/dev/sda` | Overwrites disk |
| `chmod -R 777 /` | Dangerous permissions |
| `> /etc/passwd` | Destroys user database |

### Testing Requirements
- [ ] Unit tests for SafetyPlugin patterns
- [ ] Unit tests for AliasPlugin expansion
- [ ] Hook execution order test
- [ ] Block action test
- [ ] Modify action test

---

## Feature 6: Workspace Persistence

### User Stories
- As a user, I want my pane layout saved automatically
- As a user, I want to restore my workspace on startup
- As a user, I want multiple saved workspaces

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| WS-1 | Auto-save workspace at interval (default 60s) | MUST |
| WS-2 | Restore last workspace on startup | MUST |
| WS-3 | Save workspace manually | SHOULD |
| WS-4 | Multiple workspace slots (max 10) | SHOULD |
| WS-5 | Name workspaces | SHOULD |
| WS-6 | Delete workspace | SHOULD |
| WS-7 | Export/import workspace | COULD |

### Acceptance Criteria

- [ ] Closing app saves current workspace
- [ ] Opening app restores last workspace layout
- [ ] Pane positions, sizes, and types restored
- [ ] Active pane focus restored
- [ ] Scroll position in panes restored
- [ ] Workspace list shows names and last-used dates

### Data Model
```go
type Workspace struct {
    ID         string
    Name       string
    Layout     LayoutType
    PaneStates []*PaneState
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type PaneState struct {
    ID            string
    Type          PaneType
    Title         string
    Content       string    // last content (truncated)
    ScrollOffset  int
    Position      Position  // x, y, width, height
    Focused       bool
}
```

### Testing Requirements
- [ ] Save/restore cycle test
- [ ] Multiple workspace test
- [ ] Corrupted workspace recovery test

---

## Feature 7: Secret Guard

### User Stories
- As a user, I want secrets automatically masked in output
- As a user, I want to know when a secret was detected
- As a security-conscious user, I want no secrets in logs

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| SG-1 | Detect common secret patterns | MUST |
| SG-2 | Mask detected secrets with [REDACTED] | MUST |
| SG-3 | Show notification when secret detected | SHOULD |
| SG-4 | Configurable enable/disable | SHOULD |
| SG-5 | Custom pattern support | COULD |
| SG-6 | Never log unmasked secrets | MUST |

### Acceptance Criteria

- [ ] AWS key `AKIAIOSFODNN7EXAMPLE` masked as `[AWS_KEY:AKIA...]`
- [ ] GitHub token `ghp_xxxx...` masked as `[GITHUB_TOKEN:ghp_...]`
- [ ] JWT `eyJ...` masked as `[JWT:eyJ...]`
- [ ] Notification appears: "1 secret detected and masked"
- [ ] Original value never written to disk

### Detection Patterns

| Type | Pattern | Mask Format |
|------|---------|-------------|
| AWS Access Key | `AKIA[0-9A-Z]{16}` | `[AWS_KEY:AKIA...]` |
| AWS Secret Key | `[A-Za-z0-9/+=]{40}` (after aws_secret) | `[AWS_SECRET:...]` |
| GitHub Token | `gh[ps]_[A-Za-z0-9]{36}` | `[GITHUB_TOKEN:gh*_...]` |
| GitLab Token | `glpat-[A-Za-z0-9-]{20}` | `[GITLAB_TOKEN:...]` |
| JWT | `eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+` | `[JWT:eyJ...]` |
| Private Key | `-----BEGIN.*PRIVATE KEY-----` | `[PRIVATE_KEY]` |
| Generic API Key | `[aA][pP][iI][_-]?[kK][eE][yY][=:]["']?[A-Za-z0-9]{20,}` | `[API_KEY:...]` |
| Bearer Token | `[Bb]earer [A-Za-z0-9_-]{20,}` | `[BEARER:...]` |
| Slack Token | `xox[baprs]-[A-Za-z0-9-]{10,}` | `[SLACK_TOKEN:...]` |

### Testing Requirements
- [ ] Unit tests for each pattern
- [ ] False positive test (don't mask normal strings)
- [ ] Multiple secrets in one output test
- [ ] Masking doesn't break output formatting

---

## Feature 8: Theme System

### User Stories
- As a user, I want to choose my preferred color theme
- As a user, I want themes to reload without restart
- As a power user, I want to create custom themes

### Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| TH-1 | Built-in themes: default, dracula, monokai, nord | MUST |
| TH-2 | YAML-based theme files | MUST |
| TH-3 | Hot-reload on file change | SHOULD |
| TH-4 | User theme directory (~/.terminalizcrazy/themes/) | SHOULD |
| TH-5 | Theme validation with helpful errors | SHOULD |
| TH-6 | Theme preview before applying | COULD |

### Acceptance Criteria

- [ ] Theme selector in settings
- [ ] Changing theme updates UI immediately
- [ ] Editing theme file updates UI within 1 second
- [ ] Invalid theme shows validation error, keeps current
- [ ] Custom theme in user directory appears in selector

### Theme File Format
```yaml
name: "My Theme"
description: "A custom dark theme"
author: "username"
colors:
  primary: "#7aa2f7"
  secondary: "#bb9af7"
  background: "#1a1b26"
  foreground: "#c0caf5"
  success: "#9ece6a"
  warning: "#e0af68"
  error: "#f7768e"
  info: "#7dcfff"
  muted: "#565f89"
  border: "#3b4261"
```

### Testing Requirements
- [ ] Theme loading test
- [ ] Theme validation test
- [ ] Hot-reload test
- [ ] Invalid theme recovery test

---

## Non-Functional Requirements

### Performance

| Metric | Target | Measurement |
|--------|--------|-------------|
| Startup time | < 500ms | Time from exec to first render |
| First AI token | < 100ms | Time from send to first stream chunk |
| UI frame rate | 60fps | No visible lag during typing |
| Memory usage | < 100MB | Steady state after startup |
| SQLite queries | < 10ms | 95th percentile |

### Security

| Requirement | Implementation |
|-------------|----------------|
| API keys at rest | AES-256 encryption via internal/crypto |
| API keys in memory | Cleared after use where possible |
| Collaboration | ECDH key exchange + AES-256-GCM |
| Secret detection | Mask before display and logging |
| No secrets in logs | SecretGuard applied to all output |

### Reliability

| Scenario | Recovery |
|----------|----------|
| Database corruption | Detect on startup, offer to recreate |
| Config file invalid | Fall back to defaults, show warning |
| AI provider down | Show error, continue with other features |
| WebSocket disconnect | Auto-reconnect with exponential backoff |

### Compatibility

| Platform | Terminal | Status |
|----------|----------|--------|
| Linux | GNOME Terminal, Konsole, Alacritty | Supported |
| macOS | Terminal.app, iTerm2, Alacritty | Supported |
| Windows | Windows Terminal, ConEmu | Supported |
| Windows | CMD, PowerShell (legacy) | Limited |

### Accessibility

| Feature | Status |
|---------|--------|
| Keyboard-only navigation | Supported |
| Screen reader | Not yet supported |
| High contrast themes | Planned |
| Configurable font size | Via terminal settings |

---

## Glossary

| Term | Definition |
|------|------------|
| Agent Mode | Autonomous multi-step task execution |
| Plan | A sequence of tasks generated by AI |
| Task | Single command with verification criteria |
| Risk Level | Security classification of a command |
| Pane | A rectangular area in the TUI |
| Workspace | Saved pane layout configuration |
| Hook | Extension point for plugin execution |
| Share Code | 8-character code to join collaboration session |
