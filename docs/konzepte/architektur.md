# Architektur

> Systemuebersicht und technische Details

## Uebersicht

TerminalizCrazy ist ein KI-natives Terminal, gebaut mit Go und dem Charm.sh-Oekosystem (Bubble Tea TUI Framework).

---

## Komponenten-Diagramm

    main.go
       |
       v
    config.Load() -----> tui.Run()
                              |
                              v
                         tui.Model (Bubble Tea)
                              |
            +-----------------+------------------+
            |                 |                  |
            v                 v                  v
       ai.Service      executor.Executor   storage.Storage
       (KI-Anbieter)   (Befehlsausfuehrung)  (SQLite)

---

## Haupt-Packages

### cmd/terminalizcrazy

Einstiegspunkt der Anwendung:
- Konfiguration laden
- TUI starten
- Graceful Shutdown

### internal/config

Konfigurationsmanagement mit Viper:
- TOML-Datei: ~/.terminalizcrazy/config.toml
- Umgebungsvariablen
- Standardwerte

### internal/tui

Bubble Tea TUI-Implementierung:
- Model-Update-View Architektur
- PaneManager fuer Multi-Pane-Layout
- TabBar fuer Tab-Navigation

### internal/ai

KI-Client-Implementierungen:
- Anthropic (Claude)
- OpenAI
- Ollama (lokal)
- Agent-Modus mit Planner

### internal/executor

Befehlsausfuehrung:
- Shell-Integration (bash/cmd)
- Risikobewertung
- Timeout-Handling

### internal/storage

SQLite-Persistenz:
- Sessions und Messages
- Command History
- Agent Plans
- Workflows
- Workspaces

### internal/collab

Echtzeit-Zusammenarbeit:
- WebSocket-Server
- Ende-zu-Ende-Verschluesselung (ECDH + AES-256-GCM)
- Room-Management

### internal/secretguard

Geheimniserkennung:
- Regex-basierte Pattern-Erkennung
- Automatische Maskierung

### internal/project

Projekterkennung:
- Erkennt Go, Node, Python, Rust, Java, etc.
- Kontextinformationen fuer KI

### internal/theme

Theme-System:
- YAML-Theme-Definition
- Hot-Reload
- Default-Themes

### internal/workspace

Workspace-Management:
- Layout-Presets (quad, tall, wide, stack)
- Pane-Zustand
- Auto-Save

### internal/plugins

Hook-basiertes Plugin-System:
- pre_command / post_command
- pre_ai / post_ai
- Prioritaetsbasierte Ausfuehrung

### internal/workflows

Workflow-Templates:
- YAML-Definition
- Variablen
- Bedingte Ausfuehrung

---

## Datenfluss

### Benutzeranfrage

    1. Benutzer tippt Frage
    2. TUI sendet an ai.Service
    3. AI generiert Befehlsvorschlag
    4. SecretGuard maskiert Ausgabe
    5. TUI zeigt Ergebnis

### Befehlsausfuehrung

    1. Benutzer drueckt Ctrl+E
    2. Executor bewertet Risiko
    3. Bei Risiko: Bestaetigung anfordern
    4. Shell fuehrt Befehl aus
    5. Ausgabe durch SecretGuard
    6. Storage speichert History

### Zusammenarbeit

    1. Host startet Server (Ctrl+S)
    2. Server generiert Share-Code
    3. Gast verbindet (Ctrl+J + Code)
    4. WebSocket-Verbindung etabliert
    5. Verschluesselte Nachrichten
    6. Echtzeit-Synchronisation

---

## Bubble Tea Architektur

### Model-Update-View

TerminalizCrazy verwendet das Elm-Pattern:

    type Model struct {
        // Zustand
        messages []ChatMessage
        input    textinput.Model
        viewport viewport.Model
        // ...
    }
    
    func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
    func (m Model) View() string

### Messages

Async-Operationen werden ueber Messages abgewickelt:

    type aiResponseMsg struct {
        response *ai.Response
        err      error
    }
    
    type cmdResultMsg struct {
        result *executor.Result
    }

### Commands

Commands sind Funktionen, die Messages zurueckgeben:

    func (m *Model) sendAIRequest(input string) tea.Cmd {
        return func() tea.Msg {
            resp, err := m.aiService.ProcessInput(ctx, input)
            return aiResponseMsg{response: resp, err: err}
        }
    }

---

## KI-Integration

### Client Interface

    type Client interface {
        Complete(ctx context.Context, req *Request) (*Response, error)
        Provider() Provider
    }

### Request/Response

    type Request struct {
        Messages []Message
        Context  *RequestContext  // Projekt-Info
    }
    
    type Response struct {
        Content string
        Command string  // Extrahierter Befehl
    }

### Planner

Der Agent-Modus nutzt einen Planner:

    type Planner struct {
        client Client
    }
    
    func (p *Planner) CreatePlan(task string) (*Plan, error)
    
    type Plan struct {
        Tasks []Task
    }
    
    type Task struct {
        Command      string
        Verification Verification
    }

---

## Storage Schema

### Sessions

    CREATE TABLE sessions (
        id TEXT PRIMARY KEY,
        name TEXT,
        work_dir TEXT,
        created_at DATETIME,
        updated_at DATETIME
    )

### Messages

    CREATE TABLE messages (
        id INTEGER PRIMARY KEY,
        session_id TEXT,
        role TEXT,
        content TEXT,
        command TEXT,
        success BOOLEAN,
        created_at DATETIME
    )

### Command History

    CREATE TABLE command_history (
        id INTEGER PRIMARY KEY,
        command TEXT,
        output TEXT,
        success BOOLEAN,
        duration_ms INTEGER,
        created_at DATETIME
    )

### Agent Plans

    CREATE TABLE agent_plans (
        id TEXT PRIMARY KEY,
        task TEXT,
        status TEXT,
        created_at DATETIME
    )
    
    CREATE TABLE agent_tasks (
        id INTEGER PRIMARY KEY,
        plan_id TEXT,
        command TEXT,
        verification TEXT,
        status TEXT,
        output TEXT
    )

---

## Erweiterbarkeit

### Neue KI-Anbieter

1. ai.Client Interface implementieren
2. Provider-Konstante hinzufuegen
3. In ai.NewService registrieren

### Neue Plugins

1. plugins.Plugin Interface implementieren
2. Prioritaet festlegen
3. Mit plugins.Register() registrieren

### Neue Themes

1. YAML-Datei in ~/.terminalizcrazy/themes/
2. Farbpalette definieren
3. Theme-Namen in config.toml setzen

---

## Siehe auch

- [Einstellungen](../referenz/einstellungen.md) - Konfiguration
- [Plugins](../anleitungen/plugins.md) - Plugin-Entwicklung
- [Themes](../referenz/themes.md) - Theme-Erstellung
