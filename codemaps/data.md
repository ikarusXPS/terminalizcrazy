# Data Models Codemap

> Freshness: 2026-04-17 | Auto-generated from source analysis

## SQLite Schema

Database: `~/.terminalizcrazy/terminalizcrazy.db`

### Sessions Table

```sql
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    work_dir TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Messages Table

```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    role TEXT NOT NULL,        -- user, ai, system, output
    content TEXT NOT NULL,
    command TEXT,              -- extracted command (if any)
    success INTEGER,           -- for output messages
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);
```

### Command History Table

```sql
CREATE TABLE command_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    command TEXT NOT NULL,
    output TEXT,
    success INTEGER NOT NULL,
    duration_ms INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Agent Plans Table

```sql
CREATE TABLE agent_plans (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,      -- pending, approved, running, completed, failed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);
```

### Agent Tasks Table

```sql
CREATE TABLE agent_tasks (
    id TEXT PRIMARY KEY,
    plan_id TEXT NOT NULL,
    sequence INTEGER NOT NULL,
    command TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,      -- pending, running, completed, failed, skipped
    output TEXT,
    error TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES agent_plans(id)
);
```

### Workflows Table

```sql
CREATE TABLE workflows (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    steps TEXT NOT NULL,       -- JSON serialized
    variables TEXT,            -- JSON serialized
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Workspaces Table

```sql
CREATE TABLE workspaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    layout TEXT NOT NULL,      -- quad, tall, wide, stack, single
    pane_states TEXT NOT NULL, -- JSON serialized
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Go Struct Definitions

### AI Domain

```go
// internal/ai/ai.go
type Request struct {
    UserMessage string
    Context     *RequestContext
    Type        RequestType
}

type RequestContext struct {
    CurrentDir       string
    OS               string
    Shell            string
    RecentHistory    []string
    ProjectName      string
    ProjectType      string
    ProjectFramework string
}

type Response struct {
    Content     string
    Command     string
    Explanation string
    Confidence  float64
    Provider    Provider
}

type StreamingResponse struct {
    Delta    string
    Done     bool
    Command  string
    FullText string
    Err      error
}
```

### Agent Domain

```go
// internal/ai/agent.go
type Agent struct {
    planner   *Planner
    executor  *executor.Executor
    mode      AgentMode
    maxTasks  int
    onPlanReady func(*Plan)
    onTaskStart func(*Task)
    onTaskDone  func(*Task, error)
}

type AgentConfig struct {
    Mode          AgentMode
    MaxTasks      int
    AutoApprove   bool
    RiskThreshold executor.RiskLevel
}
```

### Plan Domain

```go
// internal/ai/planner.go
type Plan struct {
    ID          string
    Name        string
    Description string
    Tasks       []*Task
    Status      PlanStatus
    Context     *PlanContext
    CreatedAt   time.Time
}

type Task struct {
    ID           string
    Sequence     int
    Command      string
    Description  string
    Status       TaskStatus
    Output       string
    Error        string
    Verification *Verification
}

type Verification struct {
    Type           VerificationType
    ExitCode       int
    OutputContains string
    RunCommand     string
}
```

### Storage Domain

```go
// internal/storage/storage.go
type Session struct {
    ID        string
    Name      string
    WorkDir   string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Message struct {
    ID        int64
    SessionID string
    Role      string
    Content   string
    Command   string
    Success   bool
    CreatedAt time.Time
}

type CommandHistory struct {
    ID         int64
    Command    string
    Output     string
    Success    bool
    DurationMs int64
    CreatedAt  time.Time
}

type AgentPlan struct {
    ID          string
    SessionID   string
    Name        string
    Description string
    Status      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type AgentTask struct {
    ID          string
    PlanID      string
    Sequence    int
    Command     string
    Description string
    Status      string
    Output      string
    Error       string
    CreatedAt   time.Time
}
```

### Collaboration Domain

```go
// internal/collab/types.go
type User struct {
    ID       string
    Name     string
    Color    string
    IsHost   bool
    JoinedAt time.Time
}

type Message struct {
    Type      MessageType
    UserID    string
    UserName  string
    Content   string
    Command   string
    Timestamp time.Time
    Encrypted bool
}

type Room struct {
    ID        string
    ShareCode string
    HostID    string
    Users     map[string]*User
    Messages  []*Message
    CreatedAt time.Time
}

// internal/collab/crypto.go
type CryptoSession struct {
    privateKey *ecdsa.PrivateKey
    publicKey  *ecdsa.PublicKey
    sharedKey  []byte
}

type EncryptedMessage struct {
    Nonce      []byte
    Ciphertext []byte
}
```

### Workspace Domain

```go
// internal/workspace/workspace.go
type Workspace struct {
    ID         string
    Name       string
    Layout     LayoutType
    PaneStates []*PaneState
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type PaneState struct {
    ID       string
    Type     PaneType
    Title    string
    Content  string
    Position Position
    Focused  bool
}

// internal/workspace/layout.go
type LayoutConfig struct {
    Width      int
    Height     int
    MinWidth   int
    MinHeight  int
    Padding    int
}

type LayoutResult struct {
    Positions []PanePosition
}

type PanePosition struct {
    X      int
    Y      int
    Width  int
    Height int
}
```

### Workflow Domain

```go
// internal/workflows/workflow.go
type Workflow struct {
    ID          string
    Name        string
    Description string
    Steps       []*WorkflowStep
    Variables   []*Variable
    CreatedAt   time.Time
}

type WorkflowStep struct {
    ID          string
    Name        string
    Command     string
    Description string
    OnFail      OnFailAction
    Condition   string
}

type Variable struct {
    Name        string
    Description string
    Default     string
    Required    bool
}

type WorkflowExecution struct {
    ID         string
    WorkflowID string
    Status     ExecutionStatus
    Variables  map[string]string
    Results    []*StepResult
    StartedAt  time.Time
    EndedAt    time.Time
}
```

### Plugin Domain

```go
// internal/plugins/plugin.go
type Plugin interface {
    Name() string
    Type() PluginType
    Hooks() []HookType
    Priority() int
    Initialize(config map[string]interface{}) error
    Execute(ctx context.Context, hookCtx *HookContext) (*HookResult, error)
}

type HookContext struct {
    Type       HookType
    Command    string
    Output     string
    AIRequest  string
    AIResponse string
    Metadata   map[string]interface{}
}

type HookResult struct {
    Modified bool
    Command  string
    Output   string
    Block    bool
    Message  string
}
```

### Theme Domain

```go
// internal/theme/theme.go
type Theme struct {
    Name        string
    Description string
    Author      string
    Colors      ColorPalette
}

type ColorPalette struct {
    Primary    string
    Secondary  string
    Background string
    Foreground string
    Success    string
    Warning    string
    Error      string
    Info       string
    Muted      string
    Border     string
}
```

### Config Domain

```go
// internal/config/config.go
type Config struct {
    AIProvider    string
    AnthropicKey  string
    OpenAIKey     string
    GeminiKey     string
    GeminiModel   string
    OllamaURL     string
    OllamaModel   string
    OllamaEnabled bool
    AgentMode     string
    AgentMaxTasks int
    Theme         string
    Appearance    AppearanceConfig
    Pane          PaneConfig
    Workspace     WorkspaceConfig
    Retention     RetentionConfig
    // ...
}

type RetentionConfig struct {
    MessageRetentionDays        int
    CommandHistoryRetentionDays int
    AgentPlanRetentionDays      int
    AutoCleanupEnabled          bool
}
```

## Enums / Constants

```go
// Providers
ProviderOllama    Provider = "ollama"
ProviderGemini    Provider = "gemini"
ProviderAnthropic Provider = "anthropic"
ProviderOpenAI    Provider = "openai"

// Request Types
RequestTypeCommand RequestType = "command"
RequestTypeExplain RequestType = "explain"
RequestTypeChat    RequestType = "chat"

// Agent Modes
AgentModeOff     AgentMode = "off"
AgentModeSuggest AgentMode = "suggest"
AgentModeAuto    AgentMode = "auto"

// Risk Levels
RiskLow      RiskLevel = 0
RiskMedium   RiskLevel = 1
RiskHigh     RiskLevel = 2
RiskCritical RiskLevel = 3

// Plan Status
PlanStatusPending   = "pending"
PlanStatusApproved  = "approved"
PlanStatusRunning   = "running"
PlanStatusCompleted = "completed"
PlanStatusFailed    = "failed"

// Layout Types
LayoutQuad   LayoutType = "quad"
LayoutTall   LayoutType = "tall"
LayoutWide   LayoutType = "wide"
LayoutStack  LayoutType = "stack"
LayoutSingle LayoutType = "single"

// Message Types (Collab)
MsgTypeChat     MessageType = "chat"
MsgTypeCommand  MessageType = "command"
MsgTypeOutput   MessageType = "output"
MsgTypeJoin     MessageType = "join"
MsgTypeLeave    MessageType = "leave"
MsgTypeUserList MessageType = "user_list"
```
