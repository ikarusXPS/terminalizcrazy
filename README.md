# ⚡ TerminalizCrazy

Ein modernes, AI-natives Terminal-Tool mit eingebauter Collaboration, intelligentem Secret-Schutz und Smart Sessions.

## Features

- **AI-Integration** - Natürliche Sprache → Befehle, Fehler-Erklärungen, Smart Autocomplete
- **Command Execution** - Vorgeschlagene Befehle direkt ausführen mit Sicherheitsabfrage
- **History Persistence** - Command-History mit SQLite, überlebt Neustarts
- **SecretGuard** - Automatische Erkennung und Maskierung von API-Keys, Tokens, Passwörtern
- **Risk Assessment** - Automatische Risikobewertung von Befehlen (Low/Medium/High/Critical)
- **Session Management** - Automatische Session-Erstellung, Chat-Persistenz und Session-Wiederherstellung
- **Smart Sessions** - Automatische Projekt-Erkennung (Go, Node, Python, Rust, etc.) mit Framework-Detection
- **Clipboard Integration** - Befehle mit Ctrl+Y in Clipboard kopieren
- **Real-Time Collaboration** - Terminal-Sharing wie Google Docs mit WebSocket

## Voraussetzungen

- **Go 1.21+** (getestet mit Go 1.26.1)
- **Git**
- API-Key von [Anthropic](https://console.anthropic.com/) oder [OpenAI](https://platform.openai.com/)

## Installation

### Aus Source bauen

```bash
# Repository klonen
git clone https://github.com/terminalizcrazy/terminalizcrazy.git
cd terminalizcrazy

# Dependencies installieren
go mod tidy

# Bauen
go build -o bin/terminalizcrazy ./cmd/terminalizcrazy

# Oder mit Make
make build
```

### Mit Go Install

```bash
go install github.com/terminalizcrazy/terminalizcrazy/cmd/terminalizcrazy@latest
```

## Setup

### 1. Umgebungsvariablen setzen

```bash
# .env Datei erstellen (aus Vorlage)
cp .env.example .env

# API-Key eintragen
# Öffne .env und setze ANTHROPIC_API_KEY oder OPENAI_API_KEY
```

**Windows (PowerShell):**
```powershell
$env:ANTHROPIC_API_KEY = "sk-ant-api03-your-key-here"
```

**Linux/macOS:**
```bash
export ANTHROPIC_API_KEY="sk-ant-api03-your-key-here"
```

### 2. Konfiguration (optional)

```bash
# Config-Verzeichnis erstellen
mkdir -p ~/.terminalizcrazy

# Konfiguration kopieren
cp config.toml.example ~/.terminalizcrazy/config.toml

# Nach Wunsch anpassen
```

## Verwendung

### Starten

```bash
# Direkt ausführen
./bin/terminalizcrazy

# Oder mit Go
go run ./cmd/terminalizcrazy

# Oder mit Make
make run
```

### Tastenkürzel

| Taste | Aktion |
|-------|--------|
| `Enter` | Nachricht an AI senden |
| `Ctrl+E` | Letzten Befehl ausführen |
| `Ctrl+Y` | Letzten Befehl in Clipboard kopieren |
| `Ctrl+R` | Letzten Befehl anzeigen |
| `Ctrl+L` | Chat leeren |
| `Ctrl+S` | Session teilen (Share) |
| `Ctrl+J` | Collaboration beitreten (Join) |
| `Ctrl+D` | Collaboration beenden (Disconnect) |
| `↑` / `↓` | History durchsuchen |
| `Y` / `N` | Bestätigung bei gefährlichen Befehlen |
| `Esc` | Abbrechen / Beenden |
| `Ctrl+C` | Sofort beenden |

### Workflow

1. **Session wählen**: Beim Start vorherige Session fortsetzen oder neu starten
2. **Frage stellen**: "how to find large files"
3. **AI antwortet** mit Befehl: `find . -size +100M`
4. **Ausführen**: `Ctrl+E` drücken
5. **Bestätigen**: Bei Medium/High Risk mit `Y` bestätigen
6. **Output sehen**: Ergebnis wird im Chat angezeigt

### Session-Wiederherstellung

Beim Start zeigt TerminalizCrazy eine Liste der letzten 10 Sessions:

```
Select Session:

▶ New Session
  abc123 • Session abc123 • 2 hours ago
  def456 • Session def456 • yesterday
```

- **↑↓**: Session auswählen
- **Enter**: Ausgewählte Session laden
- **N** oder **Esc**: Neue Session starten

Alle Chat-Nachrichten und der letzte Befehl werden wiederhergestellt.

### Smart Sessions

TerminalizCrazy erkennt automatisch den Projekttyp anhand von Konfigurationsdateien:

| Datei | Erkannter Typ | Icon |
|-------|---------------|------|
| `go.mod` | Go | 🐹 |
| `package.json` | Node.js | 📦 |
| `requirements.txt`, `pyproject.toml` | Python | 🐍 |
| `Cargo.toml` | Rust | 🦀 |
| `pom.xml`, `build.gradle` | Java | ☕ |
| `*.csproj` | .NET | 🔷 |
| `Gemfile` | Ruby | 💎 |
| `composer.json` | PHP | 🐘 |
| `Dockerfile` | Docker | 🐳 |
| `*.tf` | Terraform | 🏗️ |

**Features:**
- **Automatische Session-Namen**: `🐹 myproject (Bubble Tea)` statt `Session abc123`
- **Projekt-Kontext für AI**: AI kennt das Framework und gibt passende Befehle
- **Directory-Matching**: Sessions im gleichen Verzeichnis werden mit ★ markiert
- **Framework-Erkennung**: Next.js, Django, Spring Boot, Gin, etc.

### Real-Time Collaboration

Terminal-Sessions mit anderen teilen wie Google Docs:

**Session teilen:**
1. `Ctrl+S` drücken
2. Share-Code wird angezeigt (z.B. `a1b2-c3d4`)
3. Code an Teammitglieder senden

**Session beitreten:**
1. `Ctrl+J` drücken
2. Share-Code eingeben
3. Automatische Synchronisation

**Features:**
- 📡 WebSocket-basierte Echtzeit-Synchronisation
- 👥 User-Präsenz-Anzeige
- 💬 Geteilte Chat-Nachrichten
- ⌨️ Befehlsvorschläge sichtbar für alle
- 🔒 Lokaler Server (standardmäßig Port 8765)

## Development

### Projektstruktur

```
terminalizcrazy/
├── cmd/
│   └── terminalizcrazy/     # Entry point
│       └── main.go
├── internal/
│   ├── ai/                  # AI-Integration (Anthropic/OpenAI)
│   ├── clipboard/           # Clipboard-Integration
│   ├── collab/              # Real-Time Collaboration (WebSocket)
│   ├── config/              # Konfiguration (Viper)
│   ├── executor/            # Command Execution & Risk Assessment
│   ├── project/             # Projekt-Erkennung (Smart Sessions)
│   ├── secretguard/         # Secret-Erkennung und -Maskierung
│   ├── storage/             # SQLite Persistenz (Sessions, History)
│   └── tui/                 # Terminal UI (Bubble Tea)
├── .env.example             # Umgebungsvariablen-Vorlage
├── .gitignore
├── config.toml.example      # Config-Vorlage
├── go.mod / go.sum
├── Makefile
├── PROJECT_PLAN.md          # Projektplan und Architektur
├── SECURITY_CHECKLIST.md    # Security-Guidelines
└── README.md

# Daten-Verzeichnis (automatisch erstellt)
~/.terminalizcrazy/
├── terminalizcrazy.db       # SQLite Datenbank
└── config.toml              # User-Config (optional)
```

### Befehle

```bash
# Bauen
make build

# Tests ausführen
make test

# Tests mit Coverage
make test-coverage

# Code formatieren
make fmt

# Linter ausführen (benötigt golangci-lint)
make lint

# Für alle Plattformen bauen
make build-all

# Aufräumen
make clean
```

### Tests

```bash
# Alle Tests
go test ./...

# Mit Verbose Output
go test -v ./...

# Einzelnes Package
go test -v ./internal/secretguard/
```

## Roadmap

### Phase 1 (MVP) ✅ Complete
- [x] Projekt-Setup
- [x] Basis TUI mit Bubble Tea
- [x] SecretGuard Implementation
- [x] AI-Integration (Claude/OpenAI)
- [x] Command Execution mit Sicherheitsabfrage
- [x] Risk Assessment für Befehle

### Phase 2 ✅ Complete
- [x] Session-Management & Persistenz
- [x] Command-History (SQLite)
- [x] History Navigation mit ↑↓

### Phase 3 ✅ Complete
- [x] Clipboard-Integration
- [x] Session-Wiederherstellung
- [x] Smart Sessions (Projekt-Erkennung)

### Phase 4 ✅ Complete
- [x] Real-Time Collaboration (WebSocket)
- [x] User Presence Indicators
- [x] Message Synchronization

### Phase 5
- [ ] Erweiterte AI-Features
- [ ] Plugin-System (WASM)
- [ ] Web-Version
- [ ] Cloud-Hosted Collaboration Server

## Contributing

1. Fork das Repository
2. Feature-Branch erstellen (`git checkout -b feature/amazing-feature`)
3. Änderungen committen (`git commit -m 'feat: add amazing feature'`)
4. Branch pushen (`git push origin feature/amazing-feature`)
5. Pull Request erstellen

## Lizenz

MIT License - siehe [LICENSE](LICENSE) für Details.

## Links

- [PROJECT_PLAN.md](PROJECT_PLAN.md) - Detaillierter Projektplan
- [SECURITY_CHECKLIST.md](SECURITY_CHECKLIST.md) - Security-Guidelines
- [Charm.sh](https://charm.sh/) - TUI Framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI Library
